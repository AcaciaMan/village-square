package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// CreateSession generates a cryptographically random token, inserts a session
// row that expires in 7 days, and returns the token.
func CreateSession(db *sql.DB, userID int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(b) // 64 hex chars

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)

	_, err := db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("insert session: %w", err)
	}
	return token, nil
}

// GetSession looks up a non-expired session and returns the associated user ID.
// Returns sql.ErrNoRows if the token is missing or expired.
func GetSession(db *sql.DB, token string) (int64, error) {
	var userID int64
	err := db.QueryRow(
		"SELECT user_id FROM sessions WHERE token = ? AND expires_at > datetime('now')",
		token,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// DeleteSession removes the session row. Used for logout.
func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

// CleanExpiredSessions deletes all sessions whose expires_at is in the past.
func CleanExpiredSessions(db *sql.DB) (int64, error) {
	res, err := db.Exec("DELETE FROM sessions WHERE expires_at < datetime('now')")
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
