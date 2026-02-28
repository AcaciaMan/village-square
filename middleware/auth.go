package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	vsdb "village-square/db"
)

// contextKey is used to store values in request context.
type contextKey string

// UserIDKey is the context key for the authenticated user ID.
const UserIDKey contextKey = "userID"

// RequireAuth wraps a handler and ensures the request has a valid session cookie.
func RequireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		userID, err := vsdb.GetSession(db, cookie.Value)
		if err != nil {
			// Clear the stale/invalid cookie.
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   -1,
			})
			writeAuthError(w, http.StatusUnauthorized, "session expired")
			return
		}

		// Store user ID in context and call the next handler.
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// GetUserID extracts the authenticated user ID from the request context.
// Returns 0, false if not present.
func GetUserID(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(UserIDKey).(int64)
	return id, ok
}

// writeAuthError writes a JSON error from the middleware package.
func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
