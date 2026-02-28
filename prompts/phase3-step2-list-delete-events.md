# Phase 3 — Step 2: List Events, Delete Event & Filter API

## Context

Village Square now has an `events` table and a `POST /api/events` endpoint (auth required) that creates events with title, description, event_type, location, start_time, and end_time.

**Existing `db/events.go` exports:**
- `Event` struct (ID, UserID, Author, Title, Description, EventType, Location, StartTime, EndTime, CreatedAt)
- `CreateEvent(db, userID, title, description, eventType, location, startTime, endTime) (*Event, error)`
- `GetEventByID(db, id) (*Event, error)`

## What I need you to do

Add endpoints to list, retrieve, and delete events. **No frontend changes yet.**

### 1. Add to `db/events.go`

**`ListEvents(db *sql.DB, eventType string) ([]Event, error)`:**
- SELECT events joined with users (for author name), ordered by `start_time ASC` (upcoming first).
- If `eventType` is non-empty, filter by `event_type = ?`.
- Return an empty slice (not nil) if no events match.

**`DeleteEvent(db *sql.DB, eventID, userID int64) error`:**
- `DELETE FROM events WHERE id = ? AND user_id = ?`
- Check `RowsAffected()`. If 0 → return `ErrNotFound` (reuse the existing sentinel from `posts.go`).

### 2. Add to `handlers/events.go`

**`ListEvents(db *sql.DB) http.HandlerFunc`** for `GET /api/events`:
- **Public** (no auth required).
- Read query parameter: `?type=garage_sale`.
- Validate: if provided, `type` must be one of: `garage_sale`, `sport`, `gathering`, `other`. If invalid → `400 {"error":"invalid event type filter"}`.
- Call `db.ListEvents(db, eventType)`.
- Return `200` with:
```json
{
  "events": [ ... ],
  "count": 3
}
```

**`GetEvent(db *sql.DB) http.HandlerFunc`** for `GET /api/events/{id}`:
- **Public**.
- Parse `{id}` from `r.PathValue("id")`.
- Call `db.GetEventByID(db, id)`.
- If not found → `404 {"error":"event not found"}`.
- Return `200` with the event object.

**`DeleteEvent(db *sql.DB) http.HandlerFunc`** for `DELETE /api/events/{id}`:
- **Auth required**.
- Parse `{id}`, get user ID from context.
- Call `db.DeleteEvent(db, eventID, userID)`.
- If not found / not owner → `404 {"error":"event not found"}`.
- On success → `200 {"message":"event deleted"}`.

### 3. Update `main.go`

Add routes:
```go
mux.HandleFunc("GET /api/events", handlers.ListEvents(database))
mux.HandleFunc("GET /api/events/{id}", handlers.GetEvent(database))
mux.HandleFunc("DELETE /api/events/{id}", middleware.RequireAuth(database, handlers.DeleteEvent(database)))
```

## Acceptance criteria

- `GET /api/events` → `200` with all events, sorted by start_time ASC.
- `GET /api/events?type=garage_sale` → only garage sales.
- `GET /api/events?type=invalid` → `400`.
- `GET /api/events/{id}` → `200` with single event.
- `GET /api/events/999` → `404`.
- `DELETE /api/events/{id}` with auth (own event) → `200`.
- `DELETE /api/events/{id}` with auth (other's event) → `404`.
- `DELETE /api/events/{id}` without auth → `401`.
- Both GET endpoints work without a session cookie (public access).
- All existing endpoints still work.

## Constraints

- No frontend changes.
- Reuse `db.ErrNotFound` from `posts.go` — don't define a second sentinel.
- Follow the same handler patterns as `handlers/posts.go`.
