package db

import (
	"database/sql"
	"time"
)

// Event represents a row in the events table.
type Event struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Author      string     `json:"author"` // populated from JOIN
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EventType   string     `json:"event_type"` // garage_sale | sport | gathering | other
	Location    string     `json:"location"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"` // nullable
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateEvent inserts a new event and returns it with the author name populated.
func CreateEvent(db *sql.DB, userID int64, title, description, eventType, location string, startTime time.Time, endTime *time.Time) (*Event, error) {
	res, err := db.Exec(
		`INSERT INTO events (user_id, title, description, event_type, location, start_time, end_time)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, title, description, eventType, location, startTime, endTime,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetEventByID(db, id)
}

// GetEventByID returns a single event with the author name via JOIN.
// Returns sql.ErrNoRows if not found.
func GetEventByID(db *sql.DB, id int64) (*Event, error) {
	e := &Event{}
	err := db.QueryRow(`
		SELECT e.id, e.user_id, u.name, e.title, e.description, e.event_type,
		       e.location, e.start_time, e.end_time, e.created_at
		FROM events e
		JOIN users u ON u.id = e.user_id
		WHERE e.id = ?`, id,
	).Scan(&e.ID, &e.UserID, &e.Author, &e.Title, &e.Description, &e.EventType,
		&e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// ListEvents returns events ordered by start_time ASC, optionally filtered by event_type.
func ListEvents(db *sql.DB, eventType string) ([]Event, error) {
	query := `SELECT e.id, e.user_id, u.name, e.title, e.description, e.event_type,
		       e.location, e.start_time, e.end_time, e.created_at
		FROM events e
		JOIN users u ON u.id = e.user_id`

	var args []any
	if eventType != "" {
		query += " WHERE e.event_type = ?"
		args = append(args, eventType)
	}
	query += " ORDER BY e.start_time ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.UserID, &e.Author, &e.Title, &e.Description, &e.EventType,
			&e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// DeleteEvent deletes an event only if it belongs to the given user.
// Returns ErrNotFound if no matching row was deleted.
func DeleteEvent(db *sql.DB, eventID, userID int64) error {
	res, err := db.Exec("DELETE FROM events WHERE id = ? AND user_id = ?", eventID, userID)
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
