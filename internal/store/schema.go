package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schemaBytes []byte

func ensureSchema() error {
	dbName := getenv("AI_PLANNER_DB_NAME", "ai_planner")
	host := getenv("AI_PLANNER_DB_HOST", "127.0.0.1")
	port := getenv("AI_PLANNER_DB_PORT", "3306")
	user := getenv("AI_PLANNER_DB_USER", "planner")
	pass := getenv("AI_PLANNER_DB_PASSWORD", "")

	noDBDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4", user, pass, host, port)
	db, err := sql.Open("mysql", noDBDSN)
	if err != nil {
		return fmt.Errorf("connect for db create: %w", err)
	}
	defer db.Close()

	_, err = db.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&multiStatements=true",
		user, pass, host, port, dbName)
	db2, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	defer db2.Close()

	if _, err := db2.ExecContext(context.Background(), string(schemaBytes)); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return db2.PingContext(context.Background())
}
