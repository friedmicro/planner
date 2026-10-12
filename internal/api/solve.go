package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/solver"
)

func (s *Server) solve(w http.ResponseWriter, r *http.Request) {
	cfg, err := s.store.GetConfig(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	sleep, err := s.store.ListSleep(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	blackout, err := s.store.ListBlacklist(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	events, err := s.store.ListEvents(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	extras, err := s.store.ListExtras(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}

	sol := solver.New(cfg, sleep, blackout, events, extras).Solve()
	body, err := json.Marshal(sol)
	if err != nil {
		s.internalErr(w, err)
		return
	}
	if err := s.store.SaveSchedule(r.Context(), sol.WeekStart, body); err != nil {
		s.log.Printf("warning: could not persist schedule: %v", err)
	}
	writeJSON(w, http.StatusOK, sol)
}

func (s *Server) getSchedule(w http.ResponseWriter, r *http.Request) {
	body, err := s.store.LatestSchedule(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "no schedule solved yet; POST /api/solve first")
			return
		}
		s.internalErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
