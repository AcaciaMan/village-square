package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	vsdb "village-square/db"

	"golang.org/x/crypto/bcrypt"
)

// loginRequest is the expected JSON body for POST /api/login.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login returns a handler that authenticates a user and sets a session cookie.
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Decode body.
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Validate presence.
		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		// Look up user by email.
		user, err := vsdb.GetUserByEmail(db, req.Email)
		if err != nil {
			// Same message whether user not found or wrong password.
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		// Compare password with stored bcrypt hash.
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		// Create a server-side session.
		token, err := vsdb.CreateSession(db, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create session")
			return
		}

		// Set session cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   604800, // 7 days in seconds
			Secure:   false,  // TODO: set true in production (HTTPS)
		})

		// Return user profile (password excluded via json:"-" tag).
		writeJSON(w, http.StatusOK, user)
	}
}
