package store

import (
	"context"
	"database/sql"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Store) GetConfig(ctx context.Context) (model.Config, error) {
	var c model.Config
	row := s.db.QueryRowContext(ctx, `
		SELECT id, sleep_per_day_work_week, sleep_per_day_weekend,
		       prefer_morning_nap, preferred_rests, allow_short_rests,
		       default_min_block, schedule_extra, prefer_spacing_required_work
		FROM planner_config WHERE id = 1`)
	err := row.Scan(&c.ID, &c.SleepPerDayWorkWeek, &c.SleepPerDayWeekend,
		&c.PreferMorningNap, &c.PreferredRests, &c.AllowShortRests,
		&c.DefaultMinBlock, &c.ScheduleExtra, &c.PreferSpacingRequiredWork)
	if err == sql.ErrNoRows {
		return model.DefaultConfig(), nil
	}
	return c, err
}

func (s *Store) SaveConfig(ctx context.Context, c model.Config) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO planner_config
			(id, sleep_per_day_work_week, sleep_per_day_weekend, prefer_morning_nap,
			 preferred_rests, allow_short_rests, default_min_block, schedule_extra,
			 prefer_spacing_required_work)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			sleep_per_day_work_week = VALUES(sleep_per_day_work_week),
			sleep_per_day_weekend   = VALUES(sleep_per_day_weekend),
			prefer_morning_nap      = VALUES(prefer_morning_nap),
			preferred_rests         = VALUES(preferred_rests),
			allow_short_rests       = VALUES(allow_short_rests),
			default_min_block       = VALUES(default_min_block),
			schedule_extra          = VALUES(schedule_extra),
			prefer_spacing_required_work = VALUES(prefer_spacing_required_work)`,
		c.SleepPerDayWorkWeek, c.SleepPerDayWeekend, c.PreferMorningNap,
		c.PreferredRests, c.AllowShortRests, c.DefaultMinBlock, c.ScheduleExtra,
		c.PreferSpacingRequiredWork)
	return err
}
