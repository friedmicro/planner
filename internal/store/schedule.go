package store

import (
	"context"
	"time"
)

func (s *Store) SaveSchedule(ctx context.Context, weekStart time.Time, solution []byte) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO solved_schedules (week_start, schedule) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE schedule = VALUES(schedule), created_at = CURRENT_TIMESTAMP`,
		weekStart.Format("2006-01-02"), solution)
	return err
}

func (s *Store) LatestSchedule(ctx context.Context) ([]byte, error) {
	var b []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT schedule FROM solved_schedules ORDER BY week_start DESC LIMIT 1`).Scan(&b)
	return b, err
}