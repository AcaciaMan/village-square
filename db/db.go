package db

import (
	"database/sql"
	"fmt"
	"strings"

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

	// Valid categories: fish, produce, crafts, services, other
	const postsTable = `
	CREATE TABLE IF NOT EXISTS posts (
		id          INTEGER  PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type        TEXT     NOT NULL CHECK(type IN ('offer', 'request', 'announcement')),
		title       TEXT     NOT NULL,
		body        TEXT     NOT NULL DEFAULT '',
		category    TEXT     NOT NULL DEFAULT 'other',
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(postsTable); err != nil {
		return fmt.Errorf("create posts table: %w", err)
	}

	const eventsTable = `
	CREATE TABLE IF NOT EXISTS events (
		id          INTEGER  PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title       TEXT     NOT NULL,
		description TEXT     NOT NULL DEFAULT '',
		event_type  TEXT     NOT NULL CHECK(event_type IN ('garage_sale', 'sport', 'gathering', 'other')),
		location    TEXT     NOT NULL DEFAULT '',
		start_time  DATETIME NOT NULL,
		end_time    DATETIME,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(eventsTable); err != nil {
		return fmt.Errorf("create events table: %w", err)
	}

	// Add event_id column to posts if it doesn't exist yet.
	_, alterErr := db.Exec("ALTER TABLE posts ADD COLUMN event_id INTEGER REFERENCES events(id) ON DELETE SET NULL")
	if alterErr != nil {
		// Ignore "duplicate column name" error — means migration already ran.
		if !strings.Contains(alterErr.Error(), "duplicate column name") {
			return fmt.Errorf("add event_id column: %w", alterErr)
		}
	}

	return nil
}
