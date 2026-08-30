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
	// Users.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			account TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role INTEGER NOT NULL DEFAULT 0,
			avatar_uri TEXT NOT NULL DEFAULT '',
			profile_background_uri TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL DEFAULT (unixepoch())
		);
		`); err != nil {
		return nil, err
	}
	// Foods and orders.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS foods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			food_id TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			detail TEXT NOT NULL,
			prize REAL NOT NULL,
			rate REAL NOT NULL,
			required_time INTEGER NOT NULL,
			sold_count INTEGER NOT NULL DEFAULT 0,
			image_uris TEXT NOT NULL DEFAULT '[]',
			category INTEGER NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			deleted_at INTEGER NOT NULL DEFAULT -1
		);
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			status INTEGER NOT NULL,
			amount REAL NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			done_at INTEGER NOT NULL DEFAULT -1
		);
		CREATE TABLE IF NOT EXISTS order_foods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id TEXT NOT NULL,
			food_id TEXT NOT NULL,
			food_count INTEGER NOT NULL
		);
		`); err != nil {
		return nil, err
	}
	// Comments of foods and orders.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS food_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			food_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			detail TEXT NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			deleted_at INTEGER NOT NULL DEFAULT -1
		);
		CREATE TABLE IF NOT EXISTS order_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			detail TEXT NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			deleted_at INTEGER NOT NULL DEFAULT -1
		);
	`); err != nil {
		return nil, err
	}

	return db, nil
}
