package store

import (
	"context"
	"database/sql"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Store) listIntervals(ctx context.Context, table string) ([]model.BlockedInterval, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, day, start_h, end_h FROM "+table+" ORDER BY day, start_h")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.BlockedInterval
	for rows.Next() {
		var b model.BlockedInterval
		if err := rows.Scan(&b.ID, &b.Day, &b.Start, &b.End); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (s *Store) ListSleep(ctx context.Context) ([]model.BlockedInterval, error) {
	return s.listIntervals(ctx, "sleep_times")
}

func (s *Store) ListBlacklist(ctx context.Context) ([]model.BlockedInterval, error) {
	return s.listIntervals(ctx, "blacklisted_times")
}

func (s *Store) addInterval(ctx context.Context, table string, b model.BlockedInterval) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO "+table+" (day, start_h, end_h) VALUES (?, ?, ?)",
		b.Day, b.Start, b.End)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) AddSleep(ctx context.Context, b model.BlockedInterval) (int64, error) {
	return s.addInterval(ctx, "sleep_times", b)
}

func (s *Store) AddBlacklist(ctx context.Context, b model.BlockedInterval) (int64, error) {
	return s.addInterval(ctx, "blacklisted_times", b)
}

func (s *Store) deleteInterval(ctx context.Context, table string, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM "+table+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) DeleteSleep(ctx context.Context, id int64) error {
	return s.deleteInterval(ctx, "sleep_times", id)
}

func (s *Store) DeleteBlacklist(ctx context.Context, id int64) error {
	return s.deleteInterval(ctx, "blacklisted_times", id)
}
