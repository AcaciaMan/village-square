# Phase 3 — Step 1: Events Data Model & Create Event API

## Context

Village Square is a Go web app for a rural village community. Phases 1–2 are complete:

**Project structure:**
```
village-square/
├── main.go                  # Entry point: init DB, --seed flag, register routes, listen :8080
├── db/
│   ├── db.go                # Init(), migrate() — SQLite with WAL + foreign keys
│   ├── users.go             # User struct, GetUserByID, GetUserByEmail
│   ├── sessions.go          # CreateSession, GetSession, DeleteSession
│   ├── posts.go             # Post struct, CreatePost, GetPostByID, ListPosts, DeletePost, ErrNotFound
│   └── seed.go              # Seed() — 6 demo users + 12 posts
├── handlers/
│   ├── health.go            # GET /api/health
│   ├── register.go          # POST /api/register
│   ├── login.go             # POST /api/login
│   ├── logout.go            # POST /api/logout
│   ├── me.go                # GET /api/me (protected)
│   ├── posts.go             # POST /api/posts (protected), GET /api/posts, GET /api/posts/{id}, DELETE /api/posts/{id} (protected)
│   └── response.go          # writeJSON(w, status, data), writeError(w, status, msg)
├── middleware/
│   ├── auth.go              # RequireAuth(db, next), GetUserID(r), UserIDKey
│   └── headers.go           # SecurityHeaders(next)
├── static/
│   ├── index.html           # Landing page — register/login form
│   └── dashboard.html       # Dashboard — feed with filters, new post modal, post detail expand, delete
├── go.mod / go.sum
└── village-square.db
```

**Existing tables:** `users`, `sessions`, `posts` (id, user_id, type, title, body, category, created_at).

**Key patterns:**
- Routes: Go 1.22+ method patterns on `http.ServeMux` (e.g., `mux.HandleFunc("POST /api/posts", ...)`).
- Protected routes: `middleware.RequireAuth(db, handler)`.
- `middleware.GetUserID(r)` → `(int64, bool)`.
- `handlers.writeJSON(w, status, data)` / `handlers.writeError(w, status, msg)`.
- `db.ErrNotFound` sentinel error for "not found / not yours".

## What I need you to do

Add the `events` table and a `POST /api/events` endpoint to create Village Day events. **No frontend changes yet.**

### 1. Add the events table — update `db/db.go` migrate()

Add after the posts table creation:

```sql
CREATE TABLE IF NOT EXISTS events (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       TEXT     NOT NULL,
    description TEXT     NOT NULL DEFAULT '',
    event_type  TEXT     NOT NULL CHECK(event_type IN ('garage_sale', 'sport', 'gathering', 'other')),
    location    TEXT     NOT NULL DEFAULT '',
    start_time  DATETIME NOT NULL,
    end_time    DATETIME,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Create `db/events.go`

Define an `Event` struct:
```go
type Event struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Author      string    `json:"author"`       // populated from JOIN
    Title       string    `json:"title"`
    Description string    `json:"description"`
    EventType   string    `json:"event_type"`    // garage_sale | sport | gathering | other
    Location    string    `json:"location"`
    StartTime   time.Time `json:"start_time"`
    EndTime     *time.Time `json:"end_time"`     // nullable
    CreatedAt   time.Time `json:"created_at"`
}
```

Export these functions:

- `CreateEvent(db *sql.DB, userID int64, title, description, eventType, location string, startTime time.Time, endTime *time.Time) (*Event, error)`
  - Insert the event, then query it back with author name via JOIN.
  - Return the populated `Event` struct.

- `GetEventByID(db *sql.DB, id int64) (*Event, error)`
  - SELECT with JOIN to users for author name.
  - Return `sql.ErrNoRows` if not found.

### 3. Create `handlers/events.go`

Export `CreateEvent(db *sql.DB) http.HandlerFunc` for `POST /api/events` (auth required):

**Request** (JSON body):
```json
{
  "title": "Garage sale at Housenumber 12",
  "description": "Old furniture, kids toys, and vintage books. Everything must go!",
  "event_type": "garage_sale",
  "location": "Housenumber 12, garden",
  "start_time": "2026-06-15T09:00:00Z",
  "end_time": "2026-06-15T13:00:00Z"
}
```

**Validation:**
1. `title` — required, non-empty (trimmed), max 200 chars.
2. `description` — optional, max 2000 chars.
3. `event_type` — required, must be one of: `garage_sale`, `sport`, `gathering`, `other`.
4. `location` — optional, max 200 chars.
5. `start_time` — required, must be a valid datetime string (parse with `time.Parse(time.RFC3339, ...)`). Return `400` if unparseable.
6. `end_time` — optional. If provided, must be valid datetime and must be after `start_time`. Error: `"end_time must be after start_time"`.

Return `400` with `{"error":"<message>"}` on first validation failure.

**Response** (201):
```json
{
  "id": 1,
  "user_id": 1,
  "author": "Jan Visser",
  "title": "Garage sale at Housenumber 12",
  "description": "Old furniture, kids toys, and vintage books...",
  "event_type": "garage_sale",
  "location": "Housenumber 12, garden",
  "start_time": "2026-06-15T09:00:00Z",
  "end_time": "2026-06-15T13:00:00Z",
  "created_at": "2026-02-28T12:00:00Z"
}
```

### 4. Update `main.go`

Add the route (protected):
```go
mux.HandleFunc("POST /api/events", middleware.RequireAuth(database, handlers.CreateEvent(database)))
```

## Acceptance criteria

- `go build` succeeds.
- Server starts and auto-creates the `events` table.
- `POST /api/events` with valid session + valid JSON → `201` + event object with `author` name.
- Missing/invalid fields → `400` + clear error.
- `end_time` before `start_time` → `400`.
- No auth cookie → `401`.
- All Phase 1 & 2 endpoints still work.

## Constraints

- No frontend changes in this step.
- Use the existing `writeJSON` / `writeError` / `RequireAuth` patterns.
- Keep `db/events.go` and `handlers/events.go` each under ~120 lines.
