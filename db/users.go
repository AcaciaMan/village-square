package db

import (
	"database/sql"
	"time"
)

// User represents a row in the users table.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // never serialized
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUserByID returns the user with the given ID, or an error if not found.
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		"SELECT id, name, email, password, role, created_at FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetUserByEmail returns the user with the given email, or an error if not found.
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		"SELECT id, name, email, password, role, created_at FROM users WHERE email = ?", email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
