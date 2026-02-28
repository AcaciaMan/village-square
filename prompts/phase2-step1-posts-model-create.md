# Phase 2 — Step 1: Posts Data Model & Create Post API

## Context

Village Square is a Go web app for a rural village community. Phase 1 is complete:

**Project structure:**
```
village-square/
├── main.go                  # Entry point: init DB, register routes on mux, listen :8080
├── db/
│   ├── db.go                # Init(), migrate() — SQLite with WAL + foreign keys
│   ├── users.go             # User struct, GetUserByID, GetUserByEmail
│   └── sessions.go          # CreateSession, GetSession, DeleteSession
├── handlers/
│   ├── health.go            # GET /api/health
│   ├── register.go          # POST /api/register
│   ├── login.go             # POST /api/login
│   ├── logout.go            # POST /api/logout
│   ├── me.go                # GET /api/me (protected)
│   └── response.go          # writeJSON(w, status, data), writeError(w, status, msg)
├── middleware/
│   ├── auth.go              # RequireAuth(db, next), GetUserID(r), UserIDKey
│   └── headers.go           # SecurityHeaders(next)
├── static/
│   ├── index.html           # Landing page with register/login form
│   └── dashboard.html       # Dashboard with auth guard, welcome msg, placeholder cards
├── go.mod / go.sum
└── village-square.db        # SQLite database
```

**Existing tables:** `users` (id, name, email, password, role, created_at), `sessions` (token, user_id, created_at, expires_at).

**Key patterns:**
- Routes use Go 1.22+ method patterns on `http.ServeMux` (e.g., `mux.HandleFunc("GET /api/me", ...)`).
- Protected routes: `middleware.RequireAuth(db, handlers.Me(db))`.
- `middleware.GetUserID(r)` returns `(int64, bool)` from request context.
- `handlers.writeJSON(w, status, data)` and `handlers.writeError(w, status, msg)` for responses.
- `db.User` struct with `json:"-"` on password.

## What I need you to do

Add the `posts` table and a `POST /api/posts` endpoint to create new posts. **No frontend changes yet.**

### 1. Add the posts table — update `db/db.go`

Add to the `migrate` function:

```sql
CREATE TABLE IF NOT EXISTS posts (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        TEXT     NOT NULL CHECK(type IN ('offer', 'request', 'announcement')),
    title       TEXT     NOT NULL,
    body        TEXT     NOT NULL DEFAULT '',
    category    TEXT     NOT NULL DEFAULT 'other',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Valid categories: `fish`, `produce`, `crafts`, `services`, `other`.
The CHECK constraint for category is optional (we'll validate in Go), but add a comment noting the valid values.

### 2. Create `db/posts.go`

Define a `Post` struct:
```go
type Post struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Author    string    `json:"author"`    // populated from JOIN, not stored in posts table
    Type      string    `json:"type"`      // offer | request | announcement
    Title     string    `json:"title"`
    Body      string    `json:"body"`
    Category  string    `json:"category"`
    CreatedAt time.Time `json:"created_at"`
}
```

Export these functions:

- `CreatePost(db *sql.DB, userID int64, postType, title, body, category string) (*Post, error)`
  - Insert the post, then query it back (with author name via JOIN) to return the full object.
  - Return the populated `Post` struct including the `Author` field.

- `GetPostByID(db *sql.DB, id int64) (*Post, error)`
  - SELECT with JOIN to users to get the author name.
  - Return `sql.ErrNoRows` if not found.

### 3. Create `handlers/posts.go`

Export `CreatePost(db *sql.DB) http.HandlerFunc` for `POST /api/posts` (auth required):

**Request** (JSON body):
```json
{
  "type": "offer",
  "title": "Fresh fish available Saturday",
  "body": "Caught this morning, pick up at the harbor before noon.",
  "category": "fish"
}
```

**Logic:**
1. Only accept POST. Return `405` otherwise.
2. Decode JSON body. Return `400` on failure.
3. Validate:
   - `type` must be one of: `offer`, `request`, `announcement`. Error: `"type must be offer, request, or announcement"`.
   - `title` is non-empty (trimmed), max 200 characters. Error: `"title is required"` / `"title must be under 200 characters"`.
   - `body` is optional but max 2000 characters. Error: `"body must be under 2000 characters"`.
   - `category` must be one of: `fish`, `produce`, `crafts`, `services`, `other`. If empty, default to `"other"`.
   - Return `400` with `{"error":"<message>"}` on first validation failure.
4. Get user ID from context via `middleware.GetUserID(r)`.
5. Call `db.CreatePost(...)`.
6. Return `201` with the full post object.

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "author": "Jan",
  "type": "offer",
  "title": "Fresh fish available Saturday",
  "body": "Caught this morning, pick up at the harbor before noon.",
  "category": "fish",
  "created_at": "2026-02-28T12:00:00Z"
}
```

### 4. Update `main.go`

Add the route (protected):
```go
mux.HandleFunc("POST /api/posts", middleware.RequireAuth(database, handlers.CreatePost(database)))
```

## Acceptance criteria

- `go build` succeeds.
- Server starts and auto-creates the `posts` table.
- `POST /api/posts` with valid session cookie + valid JSON → `201` + post object with `author` name.
- Missing/invalid fields → `400` + clear error.
- No auth cookie → `401`.
- All Phase 1 endpoints still work.

## Constraints

- No frontend changes in this step.
- Keep `db/posts.go` and `handlers/posts.go` each under ~100 lines.
- Use the existing `writeJSON` / `writeError` helpers.
