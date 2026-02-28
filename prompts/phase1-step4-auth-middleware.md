# Phase 1 — Step 4: Auth Middleware & GET /api/me

## Context

Village Square now has:

- `POST /api/register` — creates user with bcrypt password.
- `POST /api/login` — verifies credentials, creates session, sets `session` cookie.
- `POST /api/logout` — deletes session, clears cookie.
- `sessions` table with token, user_id, expires_at.
- Helper functions: `CreateSession`, `GetSession`, `DeleteSession`.

## What I need you to do

Build an auth middleware that extracts the current user from the session cookie, and a `GET /api/me` endpoint that returns the logged-in user's profile. **No frontend changes yet.**

### 1. Create `middleware/auth.go`

Define a context key type and export the middleware:

```go
// ContextKey is used to store the user ID in request context.
type contextKey string
const UserIDKey contextKey = "userID"
```

Export `RequireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc`:

1. Read the `session` cookie from the request.
2. If missing → `401 {"error":"authentication required"}`.
3. Call `db.GetSession(token)` to look up the user ID.
4. If not found or expired → `401 {"error":"session expired"}`. Also clear the stale cookie.
5. Store the user ID in the request context using `context.WithValue`.
6. Call `next` with the updated request.

Also export a helper:
```go
func GetUserID(r *http.Request) (int64, bool)
```
that extracts the user ID from the context (returns 0, false if not present).

### 2. Create `handlers/me.go`

Export `Me(db *sql.DB) http.HandlerFunc` for `GET /api/me`:

1. Extract user ID from context via `middleware.GetUserID(r)`.
2. Query the `users` table for that ID.
3. Return `200` with:
```json
{
  "id": 1,
  "name": "Jan",
  "email": "jan@village.nl",
  "role": "villager",
  "created_at": "2026-02-28T12:00:00Z"
}
```
4. If user not found (deleted account) → `404 {"error":"user not found"}`.

### 3. Create `db/users.go`

Move or add user-related DB queries here:

- `GetUserByID(db *sql.DB, id int64) (*User, error)` — returns a User struct or error.
- `GetUserByEmail(db *sql.DB, email string) (*User, error)` — used by login handler.
- Define a `User` struct:
```go
type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`       // never serialized
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
```

If the login handler currently does its own SQL query inline, refactor it to use `GetUserByEmail`.

### 4. Update `main.go`

Register the new protected route:
```go
http.HandleFunc("GET /api/me", middleware.RequireAuth(db, handlers.Me(db)))
```

Use Go 1.22+ method-based routing patterns if available (`"GET /api/me"`, `"POST /api/login"`, etc.) — or check the method inside each handler. Be consistent with whatever approach the codebase already uses.

## Acceptance criteria

- `GET /api/me` without a cookie → `401`.
- `GET /api/me` with a valid session cookie → `200` + user profile JSON.
- `GET /api/me` with an expired/invalid token → `401` + cookie cleared.
- The `password` field is **never** included in any JSON response (verify via `json:"-"` tag).
- Login handler now uses `db.GetUserByEmail` instead of inline SQL.
- All existing endpoints still work.

## Constraints

- Do not use any third-party auth library.
- Do not modify the frontend.
- Keep middleware composable — `RequireAuth` wraps a `http.HandlerFunc` and returns a `http.HandlerFunc`.
