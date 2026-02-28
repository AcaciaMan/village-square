package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"village-square/db"
	"village-square/middleware"
)

// CreatePost handles POST /api/posts (auth required).
func CreatePost(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req struct {
			Type     string `json:"type"`
			Title    string `json:"title"`
			Body     string `json:"body"`
			Category string `json:"category"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		// Validate type.
		if req.Type != "offer" && req.Type != "request" && req.Type != "announcement" {
			writeError(w, http.StatusBadRequest, "type must be offer, request, or announcement")
			return
		}

		// Validate title.
		req.Title = strings.TrimSpace(req.Title)
		if req.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		if len(req.Title) > 200 {
			writeError(w, http.StatusBadRequest, "title must be under 200 characters")
			return
		}

		// Validate body.
		if len(req.Body) > 2000 {
			writeError(w, http.StatusBadRequest, "body must be under 2000 characters")
			return
		}

		// Validate / default category.
		if req.Category == "" {
			req.Category = "other"
		}
		validCats := map[string]bool{"fish": true, "produce": true, "crafts": true, "services": true, "other": true}
		if !validCats[req.Category] {
			writeError(w, http.StatusBadRequest, "category must be fish, produce, crafts, services, or other")
			return
		}

		userID, _ := middleware.GetUserID(r)

		post, err := db.CreatePost(database, userID, req.Type, req.Title, req.Body, req.Category)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create post")
			return
		}

		writeJSON(w, http.StatusCreated, post)
	}
}

// validTypes and validCategories for filter validation.
var validTypes = map[string]bool{"offer": true, "request": true, "announcement": true}
var validCategories = map[string]bool{"fish": true, "produce": true, "crafts": true, "services": true, "other": true}

// ListPosts handles GET /api/posts (public).
func ListPosts(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postType := r.URL.Query().Get("type")
		category := r.URL.Query().Get("category")

		if postType != "" && !validTypes[postType] {
			writeError(w, http.StatusBadRequest, "invalid type filter")
			return
		}
		if category != "" && !validCategories[category] {
			writeError(w, http.StatusBadRequest, "invalid category filter")
			return
		}

		posts, err := db.ListPosts(database, postType, category)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not list posts")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"posts": posts, "count": len(posts)})
	}
}

// GetPost handles GET /api/posts/{id} (public).
func GetPost(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid post id")
			return
		}

		post, err := db.GetPostByID(database, id)
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "post not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not retrieve post")
			return
		}

		writeJSON(w, http.StatusOK, post)
	}
}

// DeletePost handles DELETE /api/posts/{id} (auth required).
func DeletePost(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid post id")
			return
		}

		userID, _ := middleware.GetUserID(r)

		err = db.DeletePost(database, id, userID)
		if err == db.ErrNotFound {
			writeError(w, http.StatusNotFound, "post not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not delete post")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
	}
}
