package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Store) ListEvents(ctx context.Context) ([]model.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, time_h, min_block, max_block, days,
		       hour_earliest, hour_latest, hour_start, hour_end,
		       fixed_day, fixed_start, fixed_end
		FROM events ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Event
	for rows.Next() {
		var e model.Event
		var days sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &e.Time, &e.MinBlock, &e.MaxBlock, &days,
			&e.HourEarliest, &e.HourLatest, &e.HourStart, &e.HourEnd,
			&e.FixedDay, &e.FixedStart, &e.FixedEnd); err != nil {
			return nil, err
		}
		if days.Valid && days.String != "" && days.String != "null" {
			_ = json.Unmarshal([]byte(days.String), &e.Days)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) AddEvent(ctx context.Context, e model.Event) (int64, error) {
	days, _ := json.Marshal(e.Days)
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO events (name, time_h, min_block, max_block, days,
		                    hour_earliest, hour_latest, hour_start, hour_end,
		                    fixed_day, fixed_start, fixed_end)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.Name, e.Time, e.MinBlock, e.MaxBlock, days,
		e.HourEarliest, e.HourLatest, e.HourStart, e.HourEnd,
		e.FixedDay, e.FixedStart, e.FixedEnd)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateEvent(ctx context.Context, e model.Event) error {
	days, _ := json.Marshal(e.Days)
	res, err := s.db.ExecContext(ctx, `
		UPDATE events SET name=?, time_h=?, min_block=?, max_block=?, days=?,
		                  hour_earliest=?, hour_latest=?, hour_start=?, hour_end=?,
		                  fixed_day=?, fixed_start=?, fixed_end=?
		WHERE id=?`,
		e.Name, e.Time, e.MinBlock, e.MaxBlock, days,
		e.HourEarliest, e.HourLatest, e.HourStart, e.HourEnd,
		e.FixedDay, e.FixedStart, e.FixedEnd, e.ID)
	if err != nil {
		return err
	}
	return requireAffected(res)
}

func (s *Store) DeleteEvent(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(res)
}
