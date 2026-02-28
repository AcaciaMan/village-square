# Phase 1 — Step 1: Project Structure & SQLite Database Setup

## Context

I'm building a small web app called **Village Square** — a community board for a rural village. The Go project is already initialized with:

- `go.mod` (module `village-square`, Go 1.25.3)
- `main.go` — a simple HTTP server that serves static files from `./static` on `:8080`, with a `securityHeaders` middleware (CSP, X-Content-Type-Options, X-Frame-Options).
- `static/index.html` — a landing page with a name + email form (currently client-side only, logs to console).

## What I need you to do

Restructure the project and add SQLite database support. **Do not touch the frontend yet.**

### 1. Create this folder structure

```
village-square/
├── main.go              # Entry point: init DB, register routes, start server
├── db/
│   └── db.go            # Database init, migration, and helper functions
├── handlers/
│   └── health.go        # A simple GET /api/health endpoint (returns {"status":"ok"})
├── middleware/
│   └── headers.go       # Move the existing securityHeaders middleware here
├── static/
│   └── index.html       # (unchanged)
├── go.mod
└── go.sum
```

### 2. Database (`db/db.go`)

- Use `database/sql` with the `github.com/mattn/go-sqlite3` driver.
- Export an `Init(dbPath string) (*sql.DB, error)` function that:
  - Opens (or creates) the SQLite file at `dbPath`.
  - Enables WAL mode (`PRAGMA journal_mode=WAL`).
  - Enables foreign keys (`PRAGMA foreign_keys=ON`).
  - Calls a `migrate(db)` function to create tables.
- The `migrate` function should create the **users** table:

```sql
CREATE TABLE IF NOT EXISTS users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL,
    email       TEXT    NOT NULL UNIQUE,
    password    TEXT    NOT NULL,
    role        TEXT    NOT NULL DEFAULT 'villager',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Middleware (`middleware/headers.go`)

- Move the existing `securityHeaders` function from `main.go` into this package.
- Export it as `SecurityHeaders(next http.Handler) http.Handler`.

### 4. Health handler (`handlers/health.go`)

- Export `Health(db *sql.DB) http.HandlerFunc` that:
  - Pings the database.
  - Returns `200 {"status":"ok"}` if ping succeeds.
  - Returns `500 {"status":"error","message":"..."}` if ping fails.
  - Sets `Content-Type: application/json`.

### 5. Update `main.go`

- Import and call `db.Init("village-square.db")` at startup; log-fatal on error.
- Defer `db.Close()`.
- Register routes:
  - `GET /api/health` → `handlers.Health(db)`
  - `/` → static file server (unchanged behavior)
- Apply `middleware.SecurityHeaders` to all routes.
- Keep listening on `:8080`.

### 6. Dependencies

- Run `go get github.com/mattn/go-sqlite3` and ensure `go.sum` is updated.
- The project should compile and run with `go build -o village-square.exe .`

## Acceptance criteria

- `go build` succeeds with no errors.
- Running the server creates `village-square.db` on first start.
- `GET /api/health` returns `{"status":"ok"}` with status 200.
- `GET /` still serves the landing page exactly as before.
- The `users` table exists in the SQLite database (verify with a query or tool).

## Constraints

- No external router library — use `http.ServeMux` (or `http.HandleFunc`).
- No ORM — use raw SQL.
- Keep all files concise with clear comments.
- Do not modify `static/index.html` in this step.
