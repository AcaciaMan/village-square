package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"village-square/db"
	"village-square/middleware"
)

// ToggleInterest handles POST /api/posts/{id}/interest (auth required).
func ToggleInterest(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid post id")
			return
		}

		callerID, _ := middleware.GetUserID(r)

		post, err := db.GetPostByID(database, id, callerID)
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "post not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not retrieve post")
			return
		}

		if post.Type == "announcement" {
			writeError(w, http.StatusBadRequest, "interest not available for announcements")
			return
		}

		if callerID == post.UserID {
			writeError(w, http.StatusBadRequest, "cannot express interest in your own post")
			return
		}

		already, err := db.HasUserInterest(database, id, callerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not check interest")
			return
		}

		var interested bool
		if already {
			if err := db.DeleteInterest(database, id, callerID); err != nil {
				writeError(w, http.StatusInternalServerError, "could not remove interest")
				return
			}
			interested = false
		} else {
			if err := db.CreateInterest(database, id, callerID); err != nil {
				writeError(w, http.StatusInternalServerError, "could not add interest")
				return
			}
			interested = true
		}

		count, err := db.GetInterestCount(database, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not get interest count")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"interested":     interested,
			"interest_count": count,
		})
	}
}
