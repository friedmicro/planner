package store

import (
	"context"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Store) ListExtras(ctx context.Context) ([]model.ExtraEvent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, time_h, max_time, start_hour, end_hour
		FROM extra_events ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ExtraEvent
	for rows.Next() {
		var x model.ExtraEvent
		if err := rows.Scan(&x.ID, &x.Name, &x.Time, &x.MaxTime, &x.StartHour, &x.EndHour); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *Store) AddExtra(ctx context.Context, x model.ExtraEvent) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO extra_events (name, time_h, max_time, start_hour, end_hour)
		VALUES (?, ?, ?, ?, ?)`,
		x.Name, x.Time, x.MaxTime, x.StartHour, x.EndHour)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) DeleteExtra(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM extra_events WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(res)
}
