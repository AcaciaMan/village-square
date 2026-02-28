package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// ErrNotFound is returned when a post is not found or doesn't belong to the user.
var ErrNotFound = errors.New("not found")

// Post represents a row in the posts table.
type Post struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Author     string    `json:"author"` // populated from JOIN, not stored in posts table
	Type       string    `json:"type"`   // offer | request | announcement
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Category   string    `json:"category"`
	EventID    *int64    `json:"event_id"`    // nullable FK to events
	EventTitle *string   `json:"event_title"` // populated from LEFT JOIN to events
	CreatedAt  time.Time `json:"created_at"`
}

// CreatePost inserts a new post and returns it with the author name populated.
func CreatePost(db *sql.DB, userID int64, postType, title, body, category string, eventID *int64) (*Post, error) {
	res, err := db.Exec(
		"INSERT INTO posts (user_id, type, title, body, category, event_id) VALUES (?, ?, ?, ?, ?, ?)",
		userID, postType, title, body, category, eventID,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetPostByID(db, id)
}

// GetPostByID returns a single post with the author name via JOIN.
// Returns sql.ErrNoRows if not found.
func GetPostByID(db *sql.DB, id int64) (*Post, error) {
	p := &Post{}
	err := db.QueryRow(`
		SELECT p.id, p.user_id, u.name, p.type, p.title, p.body, p.category, p.event_id, e.title, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN events e ON e.id = p.event_id
		WHERE p.id = ?`, id,
	).Scan(&p.ID, &p.UserID, &p.Author, &p.Type, &p.Title, &p.Body, &p.Category, &p.EventID, &p.EventTitle, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ListPosts returns posts ordered by created_at DESC, optionally filtered by type and/or category.
func ListPosts(db *sql.DB, postType, category string) ([]Post, error) {
	query := `SELECT p.id, p.user_id, u.name, p.type, p.title, p.body, p.category, p.event_id, e.title, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN events e ON e.id = p.event_id`

	var conditions []string
	var args []any

	if postType != "" {
		conditions = append(conditions, "p.type = ?")
		args = append(args, postType)
	}
	if category != "" {
		conditions = append(conditions, "p.category = ?")
		args = append(args, category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY p.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Author, &p.Type, &p.Title, &p.Body, &p.Category, &p.EventID, &p.EventTitle, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// DeletePost deletes a post only if it belongs to the given user.
// Returns ErrNotFound if no matching row was deleted.
func DeletePost(db *sql.DB, postID, userID int64) error {
	res, err := db.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
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
