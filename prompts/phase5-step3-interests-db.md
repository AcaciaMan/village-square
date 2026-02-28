# Phase 5 — Step 3: Interests Table & Data Layer

## Context

Village Square is a community board app (Go + SQLite + vanilla HTML/CSS/JS). Phase 1 (Steps 1–2) added a contact button that opens a `mailto:` to the post author. Now in Phase 2 we're adding **persistent interest tracking** — recording when a user expresses interest in a post, so authors and the community can see how much traction a post has.

This step creates the **database table and data-layer functions only**. No handlers or frontend changes yet.

**Current database schema** (managed in `db/db.go` → `migrate()` function):
- `users` — id, name, email, password, role, created_at
- `sessions` — token, user_id, created_at, expires_at
- `posts` — id, user_id, type, title, body, category, event_id, created_at
- `events` — id, user_id, title, description, event_type, location, start_time, end_time, created_at

**Existing pattern for migrations:** Tables are created with `CREATE TABLE IF NOT EXISTS`. The `event_id` column on posts was added as a separate `ALTER TABLE` migration with a "duplicate column name" error check. Follow this same idempotent pattern.

**Existing data-layer pattern:** Each entity has its own file in `db/` (e.g., `db/posts.go`, `db/users.go`, `db/events.go`). Functions accept `*sql.DB` as the first argument and return typed structs + error.

## What I need you to do

### 1. Add the `interests` table — `db/db.go`

Add a new migration at the end of the `migrate()` function to create the interests table:

```sql
CREATE TABLE IF NOT EXISTS interests (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    post_id    INTEGER  NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(post_id, user_id)
);
```

Key design decisions:
- **`ON DELETE CASCADE`** on both FKs — if a post or user is deleted, their interest rows are cleaned up.
- **`UNIQUE(post_id, user_id)`** — a user can only express interest once per post. Enforced at the DB level.

### 2. Create the interests data layer — `db/interests.go`

Create a new file `db/interests.go` with the following functions:

**`CreateInterest(db *sql.DB, postID, userID int64) error`**
- Inserts a row into `interests`. 
- If the UNIQUE constraint is violated (user already interested), return a specific sentinel error: `var ErrAlreadyInterested = errors.New("already interested")`.
- Check for the SQLite unique constraint error by looking for `"UNIQUE constraint failed"` in the error string (same pattern as other SQLite error checks in the codebase).

**`GetInterestCount(db *sql.DB, postID int64) (int, error)`**
- Returns `SELECT COUNT(*) FROM interests WHERE post_id = ?`.

**`HasUserInterest(db *sql.DB, postID, userID int64) (bool, error)`**
- Returns true if a row exists for this post_id + user_id combination.
- Use `SELECT 1 FROM interests WHERE post_id = ? AND user_id = ? LIMIT 1` and check for `sql.ErrNoRows`.

**`DeleteInterest(db *sql.DB, postID, userID int64) error`**
- Deletes the interest row for the given post_id + user_id.
- Returns `ErrNotFound` (already defined in `db/posts.go`) if no row was affected.

### 3. No other changes

- Do **not** modify the `Post` struct, any existing queries, any handlers, or the frontend.
- Do **not** add seed data for interests yet (that comes in a later step).
- Do **not** create any new endpoints.

## Acceptance criteria

- `go build` succeeds with no errors.
- The server starts and the `interests` table is created automatically.
- Running the server a second time doesn't fail (idempotent migration).
- The functions `CreateInterest`, `GetInterestCount`, `HasUserInterest`, and `DeleteInterest` compile and are exported.
- The sentinel error `ErrAlreadyInterested` is exported from the `db` package.
- All existing functionality still works unchanged.

## Constraints

- No ORM — raw SQL only.
- Follow the existing code conventions: one file per entity in `db/`, `*sql.DB` as first parameter, return `error` as last value.
- Keep `db/interests.go` concise — ~50–70 lines total.
- The `ErrAlreadyInterested` sentinel should be in `db/interests.go` (not in `db/posts.go` with `ErrNotFound`).
