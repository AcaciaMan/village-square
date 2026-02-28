# Phase 5 — Step 4: Interest Endpoint & Post JSON Enhancement

## Context

Village Square is a community board app (Go + SQLite + vanilla HTML/CSS/JS). In the previous step we created:

- An `interests` table: `(id, post_id, user_id, created_at)` with `UNIQUE(post_id, user_id)`.
- Data-layer functions in `db/interests.go`:
  - `CreateInterest(db, postID, userID) error`
  - `GetInterestCount(db, postID) (int, error)`
  - `HasUserInterest(db, postID, userID) (bool, error)`
  - `DeleteInterest(db, postID, userID) error`
  - Sentinel: `ErrAlreadyInterested`

**Current post JSON** returned by `GET /api/posts` and `GET /api/posts/{id}`:
```json
{
  "id": 1,
  "user_id": 1,
  "author": "Jan Visser",
  "type": "offer",
  "title": "Fresh herring from this morning",
  "body": "...",
  "category": "fish",
  "event_id": null,
  "event_title": null,
  "created_at": "2026-02-28T..."
}
```

**Current `Post` struct** in `db/posts.go`:
```go
type Post struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"user_id"`
    Author     string    `json:"author"`
    Type       string    `json:"type"`
    Title      string    `json:"title"`
    Body       string    `json:"body"`
    Category   string    `json:"category"`
    EventID    *int64    `json:"event_id"`
    EventTitle *string   `json:"event_title"`
    CreatedAt  time.Time `json:"created_at"`
}
```

**Current routes (relevant):**
```
POST /api/posts           (auth) → handlers.CreatePost
GET  /api/posts                  → handlers.ListPosts
GET  /api/posts/{id}             → handlers.GetPost
DELETE /api/posts/{id}    (auth) → handlers.DeletePost
GET  /api/posts/{id}/contact (auth) → handlers.GetPostContact
```

**Existing handler pattern:** Each handler file uses `writeJSON(w, status, data)` and `writeError(w, status, message)` from `handlers/response.go`. Auth handlers use `middleware.GetUserID(r)` to get the caller's user ID.

## What I need you to do

Two things: (A) add an endpoint for toggling interest, and (B) include interest counts in the existing post JSON.

### 1. New handler — `handlers/interest.go`

Create a new file `handlers/interest.go` with a single handler:

**`ToggleInterest(database *sql.DB) http.HandlerFunc`**

Handles `POST /api/posts/{id}/interest` (auth required):

1. Parse `{id}` from `r.PathValue("id")` as int64. Return 400 if invalid.
2. Verify the post exists via `db.GetPostByID(database, id)`. Return 404 if not found.
3. The post must be `offer` or `request` — return 400 with `"interest not available for announcements"` for announcements.
4. The user cannot express interest in their own post. Compare `callerID == post.UserID` → return 400 with `"cannot express interest in your own post"`.
5. Check if the user has already expressed interest via `db.HasUserInterest(database, id, callerID)`.
   - **If already interested:** call `db.DeleteInterest(database, id, callerID)` to **remove** the interest (toggle off).
   - **If not yet interested:** call `db.CreateInterest(database, id, callerID)` to **add** the interest (toggle on).
6. After toggling, get the updated count via `db.GetInterestCount(database, id)`.
7. Return 200 with JSON:
   ```json
   {
     "interested": true,
     "interest_count": 5
   }
   ```
   Where `interested` is `true` if we just added interest, or `false` if we just removed it.

### 2. Add interest fields to the `Post` struct — `db/posts.go`

Add two new fields to the `Post` struct:

```go
InterestCount  int  `json:"interest_count"`
UserInterested bool `json:"user_interested"`
```

These are **not stored in the posts table** — they are populated after querying, similar to how `Author` and `EventTitle` are populated via JOINs.

### 3. Populate interest data in post queries — `db/posts.go`

Modify `GetPostByID` and `ListPosts` to accept an optional `callerUserID int64` parameter (use 0 to mean "no authenticated user / public request"):

**`GetPostByID(db *sql.DB, id int64, callerUserID int64) (*Post, error)`:**
- After the existing Scan, add two follow-up queries:
  - `GetInterestCount(db, p.ID)` → set `p.InterestCount`
  - If `callerUserID > 0`: `HasUserInterest(db, p.ID, callerUserID)` → set `p.UserInterested`

**`ListPosts(db *sql.DB, postType, category string, callerUserID int64) ([]Post, error)`:**
- After scanning all posts, loop through them and populate `InterestCount` and `UserInterested` for each post using `GetInterestCount` and `HasUserInterest`.
- This is an N+1 query pattern, which is fine for village-scale data (< 100 posts). No need to optimize with subqueries now.

**`CreatePost`** also calls `GetPostByID` internally — update that call too. Use the `userID` that was passed to `CreatePost` as the `callerUserID`.

### 4. Update all callers of `GetPostByID` and `ListPosts`

Since the function signatures changed, update all callers:

**`handlers/posts.go`:**
- `ListPosts` handler: extract the caller user ID. Since `GET /api/posts` is a **public** endpoint (no `RequireAuth`), the user may or may not be logged in. Try to read the session cookie and resolve it — but if there's no cookie or the session is invalid, just pass `callerUserID = 0`. Use this approach:
  ```go
  var callerUserID int64
  if cookie, err := r.Cookie("session"); err == nil {
      if uid, err := db.GetSession(database, cookie.Value); err == nil {
          callerUserID = uid
      }
  }
  ```
  Then call `db.ListPosts(database, postType, category, callerUserID)`.

- `GetPost` handler: same pattern — try to extract user ID from session cookie, pass to `db.GetPostByID`. This endpoint is also public.

- `DeletePost` handler: has auth — use `middleware.GetUserID(r)` before calling any DB function that needs it.

- `CreatePost` handler: already has auth — the `callerUserID` is the creator's `userID`.

**`handlers/contact.go`:**
- `GetPostContact`: calls `db.GetPostByID` — pass the caller's user ID (already available from `middleware.GetUserID`).

### 5. Register the new route — `main.go`

Add the new route, grouped with the other post interest/contact routes:

```go
mux.HandleFunc("POST /api/posts/{id}/interest", middleware.RequireAuth(database, handlers.ToggleInterest(database)))
```

Place it right after the existing `GET /api/posts/{id}/contact` line.

### 6. No other changes

- Do **not** modify the frontend in this step.
- Do **not** modify `db/interests.go` (already done in Step 3).
- Do **not** add seed data for interests.

## Acceptance criteria

- `go build` succeeds with no errors.
- `GET /api/posts` now returns posts with `"interest_count": 0` and `"user_interested": false` for each post (no interests recorded yet).
- `POST /api/posts/1/interest` (with auth, post 1 by another user) returns:
  ```json
  { "interested": true, "interest_count": 1 }
  ```
- Calling `POST /api/posts/1/interest` **again** (same user) toggles it off:
  ```json
  { "interested": false, "interest_count": 0 }
  ```
- `POST /api/posts/1/interest` without auth returns 401.
- Interest on announcements returns 400.
- Interest on your own post returns 400.
- `GET /api/posts` while logged in shows `"user_interested": true` for posts where the current user has expressed interest.
- `GET /api/posts` while not logged in shows `"user_interested": false` for all posts.
- All existing functionality (create, list, delete, contact) still works.

## Constraints

- No external libraries.
- The N+1 query pattern for populating interest counts is acceptable at this scale. Do not add complex subqueries or change the SQL JOIN structure.
- Keep the toggle logic simple — one endpoint, idempotent behavior.
- Follow existing code style and conventions.
