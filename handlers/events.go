package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"village-square/db"
	"village-square/middleware"
)

// validEventTypes lists allowed event_type values.
var validEventTypes = map[string]bool{
	"garage_sale": true,
	"sport":       true,
	"gathering":   true,
	"other":       true,
}

// CreateEvent handles POST /api/events (auth required).
func CreateEvent(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			EventType   string `json:"event_type"`
			Location    string `json:"location"`
			StartTime   string `json:"start_time"`
			EndTime     string `json:"end_time"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
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

		// Validate description.
		if len(req.Description) > 2000 {
			writeError(w, http.StatusBadRequest, "description must be under 2000 characters")
			return
		}

		// Validate event_type.
		if !validEventTypes[req.EventType] {
			writeError(w, http.StatusBadRequest, "event_type must be garage_sale, sport, gathering, or other")
			return
		}

		// Validate location.
		if len(req.Location) > 200 {
			writeError(w, http.StatusBadRequest, "location must be under 200 characters")
			return
		}

		// Validate start_time.
		if req.StartTime == "" {
			writeError(w, http.StatusBadRequest, "start_time is required")
			return
		}
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			writeError(w, http.StatusBadRequest, "start_time must be a valid RFC3339 datetime")
			return
		}

		// Validate end_time (optional).
		var endTime *time.Time
		if req.EndTime != "" {
			t, err := time.Parse(time.RFC3339, req.EndTime)
			if err != nil {
				writeError(w, http.StatusBadRequest, "end_time must be a valid RFC3339 datetime")
				return
			}
			if !t.After(startTime) {
				writeError(w, http.StatusBadRequest, "end_time must be after start_time")
				return
			}
			endTime = &t
		}

		userID, _ := middleware.GetUserID(r)

		event, err := db.CreateEvent(database, userID, req.Title, req.Description, req.EventType, req.Location, startTime, endTime)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create event")
			return
		}

		writeJSON(w, http.StatusCreated, event)
	}
}

// ListEvents handles GET /api/events (public).
func ListEvents(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventType := r.URL.Query().Get("type")

		if eventType != "" && !validEventTypes[eventType] {
			writeError(w, http.StatusBadRequest, "invalid event type filter")
			return
		}

		events, err := db.ListEvents(database, eventType)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not list events")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"events": events, "count": len(events)})
	}
}

// GetEvent handles GET /api/events/{id} (public).
func GetEvent(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid event id")
			return
		}

		event, err := db.GetEventByID(database, id)
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not retrieve event")
			return
		}

		writeJSON(w, http.StatusOK, event)
	}
}

// DeleteEvent handles DELETE /api/events/{id} (auth required).
func DeleteEvent(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid event id")
			return
		}

		userID, _ := middleware.GetUserID(r)

		err = db.DeleteEvent(database, id, userID)
		if err == db.ErrNotFound {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not delete event")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "event deleted"})
	}
}
