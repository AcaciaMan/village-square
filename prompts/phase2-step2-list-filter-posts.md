# Phase 2 — Step 2: List Posts & Filter API

## Context

Village Square now has a `posts` table and a `POST /api/posts` endpoint (auth required) that creates posts with type, title, body, and category. Each post has an `author` field populated via JOIN.

**Existing `db/posts.go` exports:**
- `Post` struct (ID, UserID, Author, Type, Title, Body, Category, CreatedAt)
- `CreatePost(db, userID, postType, title, body, category) (*Post, error)`
- `GetPostByID(db, id) (*Post, error)`

## What I need you to do

Add endpoints to list and retrieve posts. **No frontend changes yet.**

### 1. Add to `db/posts.go`

Export `ListPosts(db *sql.DB, postType, category string) ([]Post, error)`:

- SELECT posts joined with users (for author name), ordered by `created_at DESC`.
- If `postType` is non-empty, filter by `type = ?`.
- If `category` is non-empty, filter by `category = ?`.
- Both filters can be active at the same time.
- Return an empty slice (not nil) if no posts match.
- Build the query dynamically based on which filters are set (use `WHERE` clauses and a `[]any` args slice).

### 2. Add to `handlers/posts.go`

Export `ListPosts(db *sql.DB) http.HandlerFunc` for `GET /api/posts`:

- **This endpoint is public** (no auth required) — anyone can browse the feed.
- Read query parameters: `?type=offer` and `?category=fish`.
- Validate filter values if provided:
  - `type` must be one of: `offer`, `request`, `announcement` (or empty). If invalid → `400 {"error":"invalid type filter"}`.
  - `category` must be one of: `fish`, `produce`, `crafts`, `services`, `other` (or empty). If invalid → `400 {"error":"invalid category filter"}`.
- Call `db.ListPosts(db, postType, category)`.
- Return `200` with:
```json
{
  "posts": [ ... ],
  "count": 3
}
```
- If no posts, return `{"posts": [], "count": 0}`.

Export `GetPost(db *sql.DB) http.HandlerFunc` for `GET /api/posts/{id}`:

- **Public** (no auth required).
- Parse `{id}` from the URL path. Use `r.PathValue("id")` (Go 1.22+) or parse from the URL manually.
- Call `db.GetPostByID(db, id)`.
- If not found → `404 {"error":"post not found"}`.
- Return `200` with the post object.

### 3. Update `main.go`

Add routes:
```go
mux.HandleFunc("GET /api/posts", handlers.ListPosts(database))
mux.HandleFunc("GET /api/posts/{id}", handlers.GetPost(database))
```

These are **not** wrapped with `RequireAuth` — the feed is public.

## Acceptance criteria

- `GET /api/posts` → `200` with all posts, newest first.
- `GET /api/posts?type=offer` → only offers.
- `GET /api/posts?category=fish` → only fish category.
- `GET /api/posts?type=offer&category=fish` → both filters combined.
- `GET /api/posts?type=invalid` → `400`.
- `GET /api/posts/{id}` → `200` with single post.
- `GET /api/posts/999` → `404`.
- Both endpoints work **without** a session cookie (public access).
- Create endpoint (`POST /api/posts`) still requires auth.

## Constraints

- No frontend changes.
- Keep the query builder simple — no ORM, just string concatenation with parameterized queries.
- Response wrapper `{"posts": [...], "count": N}` makes it easier for the frontend to handle empty states.
