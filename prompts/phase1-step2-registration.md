# Phase 1 — Step 2: User Registration Endpoint

## Context

Village Square now has:

- A restructured Go project with `db/`, `handlers/`, `middleware/` packages.
- SQLite database initialized on startup with a `users` table:
  ```sql
  users(id INTEGER PK, name TEXT, email TEXT UNIQUE, password TEXT, role TEXT DEFAULT 'villager', created_at DATETIME)
  ```
- `GET /api/health` working.
- Static file serving at `/`.

## What I need you to do

Implement a `POST /api/register` endpoint that creates a new user. **Do not touch the frontend yet.**

### 1. Create `handlers/register.go`

Export `Register(db *sql.DB) http.HandlerFunc` that handles `POST /api/register`:

**Request** (JSON body):
```json
{
  "name": "Jan",
  "email": "jan@village.nl",
  "password": "secret123"
}
```

**Logic:**
1. Decode the JSON body. Return `400` with `{"error":"invalid request body"}` if it fails.
2. Validate:
   - `name` is non-empty (trimmed). Error: `"name is required"`.
   - `email` is non-empty and looks like a valid email (use a simple regex or `mail.ParseAddress`). Error: `"valid email is required"`.
   - `password` is at least 6 characters. Error: `"password must be at least 6 characters"`.
   - Return `400` with `{"error":"<message>"}` on the first validation failure.
3. Hash the password using `golang.org/x/crypto/bcrypt` with default cost.
4. Insert into `users` (name, email, password_hash).
   - If the email already exists (UNIQUE constraint violation), return `409` with `{"error":"email already registered"}`.
5. Return `201` with:
```json
{
  "id": 1,
  "name": "Jan",
  "email": "jan@village.nl",
  "role": "villager",
  "created_at": "2026-02-28T12:00:00Z"
}
```

**Important details:**
- Never return the password hash in any response.
- Set `Content-Type: application/json` on all responses.
- Only accept `POST` method; return `405` for anything else.

### 2. Create a shared response helper (optional but recommended)

Create `handlers/response.go` with small helpers:
- `writeJSON(w http.ResponseWriter, status int, data any)` — marshals and writes JSON.
- `writeError(w http.ResponseWriter, status int, message string)` — writes `{"error":"..."}`.

These will be reused by every handler going forward.

### 3. Update `main.go`

- Add route: `POST /api/register` → `handlers.Register(db)`
- Run `go get golang.org/x/crypto/bcrypt`.

### 4. Method enforcement

Either handle method checks inside the handler (return 405 for non-POST), or create a tiny helper wrapper — your choice. Keep it simple.

## Acceptance criteria

- `POST /api/register` with valid JSON → `201` + user object (no password).
- Missing/invalid fields → `400` + clear error message.
- Duplicate email → `409`.
- `GET /api/register` → `405`.
- Password is stored as a bcrypt hash in SQLite (verify by querying the DB).
- All existing functionality (health, static files) still works.

## Constraints

- No sessions or login yet — that's the next step.
- No frontend changes.
- Keep the handler in a single file, under 100 lines.
