package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// registerRequest is the expected JSON body for POST /api/register.
type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// registerResponse is returned on successful registration.
type registerResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// Register returns a handler that creates a new user account.
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Method check (belt-and-suspenders; mux pattern already filters).
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Decode body.
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Validate fields.
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		req.Email = strings.TrimSpace(req.Email)
		if _, err := mail.ParseAddress(req.Email); err != nil || req.Email == "" {
			writeError(w, http.StatusBadRequest, "valid email is required")
			return
		}

		if len(req.Password) < 6 {
			writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
			return
		}

		// Hash password.
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}

		// Insert user.
		result, err := db.Exec(
			"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			req.Name, req.Email, string(hash),
		)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				writeError(w, http.StatusConflict, "email already registered")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to create user")
			return
		}

		id, _ := result.LastInsertId()

		// Read back created_at and role from the DB row.
		var role, createdAt string
		db.QueryRow("SELECT role, created_at FROM users WHERE id = ?", id).
			Scan(&role, &createdAt)

		// Normalise created_at to RFC 3339.
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			createdAt = t.UTC().Format(time.RFC3339)
		}

		writeJSON(w, http.StatusCreated, registerResponse{
			ID:        id,
			Name:      req.Name,
			Email:     req.Email,
			Role:      role,
			CreatedAt: createdAt,
		})
	}
}
