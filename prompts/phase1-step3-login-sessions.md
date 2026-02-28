# Phase 1 — Step 3: Login & Session Management

## Context

Village Square now has:

- SQLite with a `users` table, bcrypt-hashed passwords.
- `POST /api/register` — creates users, returns `201` with user JSON.
- `GET /api/health` — database health check.
- Response helpers: `writeJSON`, `writeError`.

## What I need you to do

Add a `POST /api/login` endpoint and a server-side session system using secure cookies. **No frontend changes yet.**

### 1. Create a sessions table

Add to the `migrate` function in `db/db.go`:

```sql
CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT     PRIMARY KEY,
    user_id    INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
);
```

### 2. Create `db/sessions.go`

Export these functions:

- `CreateSession(db *sql.DB, userID int64) (token string, err error)`
  - Generate a 32-byte random token using `crypto/rand`, hex-encode it (64 chars).
  - Insert into `sessions` with `expires_at` = now + 7 days.
  - Return the token.

- `GetSession(db *sql.DB, token string) (userID int64, err error)`
  - Query for the token where `expires_at > datetime('now')`.
  - Return the user ID, or an error if not found / expired.

- `DeleteSession(db *sql.DB, token string) error`
  - Delete the row. Used for logout.

### 3. Create `handlers/login.go`

Export `Login(db *sql.DB) http.HandlerFunc` for `POST /api/login`:

**Request** (JSON body):
```json
{
  "email": "jan@village.nl",
  "password": "secret123"
}
```

**Logic:**
1. Decode JSON. Return `400` if invalid.
2. Validate email and password are non-empty. Return `400` if missing.
3. Look up user by email. If not found → `401 {"error":"invalid credentials"}`.
4. Compare password with stored hash via `bcrypt.CompareHashAndPassword`. If mismatch → `401 {"error":"invalid credentials"}`.
   - Use the **same error message** for not-found and wrong-password (prevent user enumeration).
5. Create a session via `db.CreateSession`.
6. Set a cookie:
   - Name: `session`
   - Value: the token
   - Path: `/`
   - HttpOnly: `true`
   - SameSite: `Strict`
   - MaxAge: 7 days in seconds (`604800`)
   - Secure: `false` (for local dev; note in a comment to set `true` in production)
7. Return `200` with:
```json
{
  "id": 1,
  "name": "Jan",
  "email": "jan@village.nl",
  "role": "villager"
}
```

### 4. Create `handlers/logout.go`

Export `Logout(db *sql.DB) http.HandlerFunc` for `POST /api/logout`:

1. Read the `session` cookie. If missing → `200` (already logged out, no error).
2. Delete the session from the database.
3. Clear the cookie (set `MaxAge: -1`).
4. Return `200 {"message":"logged out"}`.

### 5. Update `main.go`

Register new routes:
- `POST /api/login` → `handlers.Login(db)`
- `POST /api/logout` → `handlers.Logout(db)`

## Acceptance criteria

- Register a user, then `POST /api/login` with correct credentials → `200` + user JSON + `Set-Cookie` header with `session` cookie.
- Wrong password → `401 {"error":"invalid credentials"}`.
- Non-existent email → `401 {"error":"invalid credentials"}` (same message).
- `POST /api/logout` with valid cookie → session deleted from DB, cookie cleared.
- Session tokens are cryptographically random (64 hex characters).
- Expired sessions are not accepted by `GetSession`.

## Constraints

- Do not use any third-party session library — keep it hand-rolled with `crypto/rand`.
- Do not modify the frontend.
- Do not build auth middleware yet — that's the next step.
