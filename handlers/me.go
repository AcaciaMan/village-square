package handlers

import (
	"database/sql"
	"net/http"

	vsdb "village-square/db"
	"village-square/middleware"
)

// Me returns a handler that responds with the logged-in user's profile.
func Me(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		user, err := vsdb.GetUserByID(db, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeJSON(w, http.StatusOK, user)
	}
}
