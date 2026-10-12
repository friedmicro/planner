package model

type Config struct {
	ID                        int     `json:"id"`
	SleepPerDayWorkWeek       float64 `json:"sleep_per_day_work_week"`
	SleepPerDayWeekend        float64 `json:"sleep_per_day_weekend"`
	PreferMorningNap          bool    `json:"prefer_morning_nap"`
	PreferredRests            int     `json:"preferred_rests"`
	AllowShortRests           bool    `json:"allow_short_rests"`
	DefaultMinBlock           float64 `json:"default_min_block"`
	ScheduleExtra             bool    `json:"schedule_extra"`
	PreferSpacingRequiredWork bool    `json:"prefer_spacing_required_work"`
}

func DefaultConfig() Config {
	return Config{
		ID:                        1,
		SleepPerDayWorkWeek:       7,
		SleepPerDayWeekend:        8,
		PreferMorningNap:          false,
		PreferredRests:            2,
		AllowShortRests:           false,
		DefaultMinBlock:           1,
		ScheduleExtra:             false,
		PreferSpacingRequiredWork: false,
	}
}

type BlockedInterval struct {
	ID    int     `json:"id"`
	Day   int     `json:"day"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type Event struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Time         float64 `json:"time"`
	MinBlock     float64 `json:"min_block,omitempty"`
	MaxBlock     float64 `json:"max_block,omitempty"`
	Days         []int   `json:"days,omitempty"`
	HourEarliest float64 `json:"hour_earliest,omitempty"`
	HourLatest   float64 `json:"hour_latest,omitempty"`
	HourStart    float64 `json:"hour_start,omitempty"`
	HourEnd      float64 `json:"hour_end,omitempty"`
	FixedDay     int     `json:"fixed_day,omitempty"`
	FixedStart   float64 `json:"fixed_start,omitempty"`
	FixedEnd     float64 `json:"fixed_end,omitempty"`
}

type ExtraEvent struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Time      float64 `json:"time"`
	MaxTime   float64 `json:"max_time,omitempty"`
	StartHour float64 `json:"start_hour,omitempty"`
	EndHour   float64 `json:"end_hour,omitempty"`
}
