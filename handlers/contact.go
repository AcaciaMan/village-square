package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"village-square/db"
	"village-square/middleware"
)

// GetPostContact returns a mailto URL for the author of the given post.
func GetPostContact(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid post id")
			return
		}

		post, err := db.GetPostByID(database, id)
		if err != nil {
			writeError(w, http.StatusNotFound, "post not found")
			return
		}

		if post.Type == "announcement" {
			writeError(w, http.StatusBadRequest, "contact not available for announcements")
			return
		}

		callerID, _ := middleware.GetUserID(r)
		if callerID == post.UserID {
			writeError(w, http.StatusBadRequest, "cannot contact yourself")
			return
		}

		author, err := db.GetUserByID(database, post.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to look up author")
			return
		}

		mailto := fmt.Sprintf("mailto:%s?subject=%s",
			author.Email,
			url.QueryEscape("Village Square: "+post.Title),
		)

		writeJSON(w, http.StatusOK, map[string]string{"mailto": mailto})
	}
}
