package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New(dsn string) (*Store, error) {
	return newWithAutoInit(dsn, false)
}

func NewWithAutoInit() (*Store, error) {
	dsn := DSNFromEnv()
	return newWithAutoInit(dsn, true)
}

func newWithAutoInit(dsn string, autoInit bool) (*Store, error) {
	if autoInit {
		if err := ensureSchema(); err != nil {
			return nil, fmt.Errorf("ensure tables: %w", err)
		}
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func DSNFromEnv() string {
	host := getenv("AI_PLANNER_DB_HOST", "127.0.0.1")
	port := getenv("AI_PLANNER_DB_PORT", "3306")
	user := getenv("AI_PLANNER_DB_USER", "planner")
	pass := getenv("AI_PLANNER_DB_PASSWORD", "")
	name := getenv("AI_PLANNER_DB_NAME", "ai_planner")
	fullDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&multiStatements=true",
		user, pass, host, port, name)
	return fullDSN
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
