package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize the tables.
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			account TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role INTEGER NOT NULL DEFAULT 0,
			avatar_uri TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	); err != nil {
		return nil, err
	}

	return db, nil
}
