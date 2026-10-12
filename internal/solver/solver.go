package solver

import (
	"sort"
	"time"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

type Interval struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

func (i Interval) hours() float64 { return i.End - i.Start }

type Placement struct {
	Name  string  `json:"name"`
	Day   int     `json:"day"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Kind  string  `json:"kind"`
}

type Solution struct {
	WeekStart       time.Time   `json:"week_start"`
	Placements      []Placement `json:"placements"`
	SleepSacrificed bool        `json:"sleep_sacrificed"`
	Feasible        bool        `json:"feasible"`
	Unscheduled     []string    `json:"unscheduled,omitempty"`
}

type dayFree struct {
	day    int
	blocks []Interval
}

type Solver struct {
	cfg      model.Config
	sleep    map[int][]Interval
	blackout map[int][]Interval
	events   []model.Event
	extras   []model.ExtraEvent

	free     map[int][]Interval
	occupied map[int][]Interval
}

func New(cfg model.Config, sleep, blackout []model.BlockedInterval, events []model.Event, extras []model.ExtraEvent) *Solver {
	byDay := func(items []model.BlockedInterval) map[int][]Interval {
		m := map[int][]Interval{}
		for _, it := range items {
			iv := Interval{it.Start, it.End}
			if iv.End <= iv.Start {
				continue
			}
			m[it.Day] = append(m[it.Day], iv)
		}
		for d := range m {
			m[d] = merge(m[d])
		}
		return m
	}
	s := &Solver{
		cfg:      cfg,
		sleep:    byDay(sleep),
		blackout: byDay(blackout),
		events:   events,
		extras:   extras,
		occupied: map[int][]Interval{},
		free:     map[int][]Interval{},
	}
	s.computeFree()
	return s
}

func (s *Solver) computeFree() {
	for d := 1; d <= 7; d++ {
		blocked := append(append([]Interval{}, s.blackout[d]...), s.sleep[d]...)
		s.free[d] = merge(blocked)
		s.free[d] = complement(s.free[d])
	}
}

func (s *Solver) freeBlocks(day int, includeSleep bool) []Interval {
	blocked := append([]Interval{}, s.blackout[day]...)
	blocked = append(blocked, s.occupied[day]...)
	if !includeSleep {
		blocked = append(blocked, s.sleep[day]...)
	}
	blocked = merge(blocked)
	return complement(blocked)
}

func (s *Solver) place(day int, start, end float64) {
	s.occupied[day] = append(s.occupied[day], Interval{start, end})
	s.occupied[day] = merge(s.occupied[day])
}

func (s *Solver) Solve() Solution {
	sol := Solution{WeekStart: weekStart(), Feasible: true, SleepSacrificed: false}

	// 1. Fixed events occupy exact slots (commitments, never moved).
	for _, e := range s.events {
		if e.FixedDay == 0 || e.FixedStart >= e.FixedEnd {
			continue
		}
		s.place(e.FixedDay, e.FixedStart, e.FixedEnd)
		sol.Placements = append(sol.Placements, Placement{
			Name: e.Name, Day: e.FixedDay, Start: e.FixedStart, End: e.FixedEnd, Kind: "event",
		})
	}

	// 2. Required events. Day-constrained first, then free-floating.
	var constrained, floating []model.Event
	for _, e := range s.events {
		if e.FixedDay != 0 {
			continue // already handled
		}
		if len(e.Days) > 0 {
			constrained = append(constrained, e)
		} else {
			floating = append(floating, e)
		}
	}
	for _, e := range constrained {
		s.scheduleRequired(&sol, e, true)
	}
	for _, e := range floating {
		s.scheduleRequired(&sol, e, false)
	}

	// 3. Extra events (best-effort, never sacrifice sleep).
	if s.cfg.ScheduleExtra {
		for _, x := range s.extras {
			s.scheduleExtra(&sol, x)
		}
	}

	// 4. Backfill sleep deficits as naps if enabled.
	s.backfillSleep(&sol)

	return sol
}

// scheduleRequired places one required event, splitting it into blocks.
func (s *Solver) scheduleRequired(sol *Solution, e model.Event, constrained bool) {
	minB, maxB := s.blockBounds(e)
	remaining := e.Time

	dayUsed := map[int]float64{}

	// Try free time first. If the event still has a remainder after the free
	// pass (because free time ran out), fall back into sleep to finish it —
	// commitments win, so drop sleep rather than drop work. Blackout stays hard either way.
	s.placeEvent(sol, e, &remaining, minB, maxB, constrained, false, dayUsed)
	if remaining > 0 {
		if s.placeEvent(sol, e, &remaining, minB, maxB, constrained, true, dayUsed) {
			sol.SleepSacrificed = true
		}
	}

	if remaining > 0 {
		sol.Feasible = false
		sol.Unscheduled = append(sol.Unscheduled, e.Name)
	}
}

func (s *Solver) placeEvent(sol *Solution, e model.Event, remaining *float64, minB, maxB float64, constrained, allowSleep bool, dayUsed map[int]float64) bool {
	placedAny := false
	for *remaining > 1e-9 {
		day, iv := s.bestSlot(e, minB, maxB, allowSleep, constrained, dayUsed)
		if day == 0 {
			break
		}
		block := *remaining
		if block > maxB {
			block = maxB
		}
		if dayCap := maxB - dayUsed[day]; block > dayCap {
			block = dayCap
		}
		if block <= 1e-9 {
			break
		}
		avail := iv.hours()
		if block > avail {
			block = avail
		}
		if block < minB && *remaining-block > 1e-9 {
			break
		}
		s.place(day, iv.Start, iv.Start+block)
		sol.Placements = append(sol.Placements, Placement{
			Name: e.Name, Day: day, Start: iv.Start, End: iv.Start + block, Kind: "event",
		})
		*remaining -= block
		dayUsed[day] += block
		placedAny = true
	}
	return placedAny
}

func (s *Solver) bestSlot(e model.Event, minB, maxB float64, allowSleep, constrained bool, dayUsed map[int]float64) (int, Interval) {
	type cand struct {
		day int
		iv  Interval
	}
	var cands []cand
	days := e.Days
	if !constrained {
		for d := 1; d <= 7; d++ {
			days = append(days, d)
		}
	}
	for _, d := range days {
		if dayUsed[d] >= maxB-1e-9 {
			continue
		}
		for _, iv := range s.freeBlocks(d, allowSleep) {
			lo, hi := 0.0, 24.0
			if constrained {
				lo, hi = e.HourEarliest, e.HourLatest
				if hi <= lo {
					lo, hi = e.HourStart, e.HourEnd
				}
				if hi <= lo {
					lo, hi = 0, 24
				}
			}
			clip := Interval{maxf(iv.Start, lo), minf(iv.End, hi)}
			if clip.hours() < minB-1e-9 {
				continue
			}
			cands = append(cands, cand{d, clip})
		}
	}
	if len(cands) == 0 {
		return 0, Interval{}
	}
	if s.cfg.PreferSpacingRequiredWork {
		load := map[int]float64{}
		for d := 1; d <= 7; d++ {
			for _, iv := range s.freeBlocks(d, allowSleep) {
				load[d] += iv.hours()
			}
		}
		sort.SliceStable(cands, func(i, j int) bool {
			li, lj := load[cands[i].day], load[cands[j].day]
			if li != lj {
				return li > lj
			}
			return cands[i].iv.Start < cands[j].iv.Start
		})
	} else {
		// Earliest slot first.
		sort.SliceStable(cands, func(i, j int) bool {
			if cands[i].day != cands[j].day {
				return cands[i].day < cands[j].day
			}
			return cands[i].iv.Start < cands[j].iv.Start
		})
	}
	return cands[0].day, cands[0].iv
}

func (s *Solver) scheduleExtra(sol *Solution, x model.ExtraEvent) {
	maxB := x.MaxTime
	if maxB <= 0 {
		maxB = x.Time
	}
	remaining := x.Time
	for remaining > 1e-9 {
		day, iv := s.bestExtraSlot(x)
		if day == 0 {
			break
		}
		block := remaining
		if block > maxB {
			block = maxB
		}
		if block > iv.hours() {
			block = iv.hours()
		}
		if block < 1e-9 {
			break
		}
		s.place(day, iv.Start, iv.Start+block)
		sol.Placements = append(sol.Placements, Placement{
			Name: x.Name, Day: day, Start: iv.Start, End: iv.Start + block, Kind: "extra",
		})
		remaining -= block
	}
}

func (s *Solver) bestExtraSlot(x model.ExtraEvent) (int, Interval) {
	var best Interval
	bestDay := 0
	for d := 1; d <= 7; d++ {
		for _, iv := range s.freeBlocks(d, false) {
			lo, hi := 0.0, 24.0
			if x.EndHour > x.StartHour {
				lo, hi = x.StartHour, x.EndHour
			}
			clip := Interval{maxf(iv.Start, lo), minf(iv.End, hi)}
			if clip.hours() < 1e-9 {
				continue
			}
			if bestDay == 0 || clip.Start < best.Start {
				best, bestDay = clip, d
			}
		}
	}
	return bestDay, best
}

func (s *Solver) backfillSleep(sol *Solution) {
	for d := 1; d <= 7; d++ {
		target := s.cfg.SleepPerDayWeekend
		if d <= 5 {
			target = s.cfg.SleepPerDayWorkWeek
		}
		scheduled := 0.0
		for _, iv := range s.sleep[d] {
			scheduled += iv.hours()
		}
		deficit := target - scheduled
		if deficit <= 1e-9 {
			continue
		}
		if !s.cfg.AllowShortRests && deficit < 2 && !s.cfg.PreferMorningNap {
			continue
		}
		gaps := s.freeBlocks(d, true)
		if len(gaps) == 0 {
			continue
		}
		var gap Interval
		best := -1.0
		for _, g := range gaps {
			if g.hours() > best {
				best, gap = g.hours(), g
			}
		}
		nap := deficit
		if nap > gap.hours() {
			nap = gap.hours()
		}
		if nap < 1e-9 {
			continue
		}
		if !s.cfg.AllowShortRests && nap < 2 && !s.cfg.PreferMorningNap {
			continue
		}
		s.place(d, gap.Start, gap.Start+nap)
		sol.Placements = append(sol.Placements, Placement{
			Name: "Nap (sleep deficit)", Day: d, Start: gap.Start, End: gap.Start + nap, Kind: "nap",
		})
	}
}

func (s *Solver) blockBounds(e model.Event) (float64, float64) {
	minB := e.MinBlock
	if minB <= 0 {
		minB = s.cfg.DefaultMinBlock
	}
	maxB := e.MaxBlock
	if maxB <= 0 {
		maxB = e.Time
	}
	if maxB < minB {
		maxB = minB
	}
	return minB, maxB
}

func merge(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	ivs := append([]Interval{}, intervals...)
	sort.Slice(ivs, func(i, j int) bool { return ivs[i].Start < ivs[j].Start })
	out := []Interval{ivs[0]}
	for _, iv := range ivs[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End+1e-9 {
			if iv.End > last.End {
				last.End = iv.End
			}
		} else {
			out = append(out, iv)
		}
	}
	return out
}

func complement(intervals []Interval) []Interval {
	ivs := merge(intervals)
	var out []Interval
	cur := 0.0
	for _, iv := range ivs {
		if iv.Start > cur+1e-9 {
			out = append(out, Interval{cur, iv.Start})
		}
		if iv.End > cur {
			cur = iv.End
		}
	}
	if cur < 24-1e-9 {
		out = append(out, Interval{cur, 24})
	}
	return out
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func weekStart() time.Time {
	now := time.Now()
	wd := int(now.Weekday()) // 0=Sun
	if wd == 0 {
		wd = 7
	}
	monday := now.AddDate(0, 0, -(wd - 1))
	y, m, d := monday.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, monday.Location())
}
