package handlers

import (
	"database/sql"
	"net/http"

	vsdb "village-square/db"
)

// Logout returns a handler that destroys the session and clears the cookie.
func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Read the session cookie; if missing the user is already logged out.
		cookie, err := r.Cookie("session")
		if err != nil {
			writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
			return
		}

		// Delete session from the database (ignore errors — best effort).
		_ = vsdb.DeleteSession(db, cookie.Value)

		// Clear the cookie in the browser.
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})

		writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
	}
}
