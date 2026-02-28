package db

import (
	"database/sql"
	"fmt"

	// SQLite driver — imported for side-effect registration.
	_ "github.com/mattn/go-sqlite3"
)

// Init opens (or creates) the SQLite database at dbPath,
// enables WAL mode and foreign keys, and runs migrations.
func Init(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Enable WAL mode for better concurrent-read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}

	// Enable foreign-key constraint enforcement.
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	// Run schema migrations.
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// migrate creates tables that don't already exist.
func migrate(db *sql.DB) error {
	const usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT    NOT NULL,
		email       TEXT    NOT NULL UNIQUE,
		password    TEXT    NOT NULL,
		role        TEXT    NOT NULL DEFAULT 'villager',
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	const sessionsTable = `
	CREATE TABLE IF NOT EXISTS sessions (
		token      TEXT     PRIMARY KEY,
		user_id    INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
	);`

	if _, err := db.Exec(sessionsTable); err != nil {
		return fmt.Errorf("create sessions table: %w", err)
	}

	return nil
}
