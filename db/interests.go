package db

import (
	"database/sql"
	"errors"
	"strings"
)

// ErrAlreadyInterested is returned when a user tries to express interest twice.
var ErrAlreadyInterested = errors.New("already interested")

// CreateInterest records a user's interest in a post.
func CreateInterest(db *sql.DB, postID, userID int64) error {
	_, err := db.Exec(
		"INSERT INTO interests (post_id, user_id) VALUES (?, ?)",
		postID, userID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrAlreadyInterested
		}
		return err
	}
	return nil
}

// GetInterestCount returns the number of interests for a post.
func GetInterestCount(db *sql.DB, postID int64) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM interests WHERE post_id = ?", postID).Scan(&count)
	return count, err
}

// HasUserInterest returns true if the user has expressed interest in the post.
func HasUserInterest(db *sql.DB, postID, userID int64) (bool, error) {
	var one int
	err := db.QueryRow(
		"SELECT 1 FROM interests WHERE post_id = ? AND user_id = ? LIMIT 1",
		postID, userID,
	).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteInterest removes a user's interest in a post.
func DeleteInterest(db *sql.DB, postID, userID int64) error {
	res, err := db.Exec(
		"DELETE FROM interests WHERE post_id = ? AND user_id = ?",
		postID, userID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
