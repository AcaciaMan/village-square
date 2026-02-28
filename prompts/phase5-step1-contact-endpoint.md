# Phase 5 — Step 1: "I'm Interested" Contact Endpoint (Backend)

## Context

Village Square is a community board app (Go + SQLite + vanilla HTML/CSS/JS). We're adding a new feature: an "I'm interested" / "I can help!" button on posts so villagers can easily contact the post author.

Phase 1 of this feature is the **smallest demo-able slice**: a button that opens a `mailto:` link to the post author. This step covers the **backend only** — the frontend button will be wired up in Step 2.

**Current backend routes:**
```
GET    /api/health
POST   /api/register
POST   /api/login
POST   /api/logout
GET    /api/me              (auth)
POST   /api/posts           (auth)
GET    /api/posts
GET    /api/posts/{id}
DELETE /api/posts/{id}      (auth)
POST   /api/events          (auth)
GET    /api/events
GET    /api/events/{id}
DELETE /api/events/{id}     (auth)
```

**Relevant existing code:**

- `db/users.go` has `GetUserByID(db, id) (*User, error)` — the `User` struct includes `Email` and `Name`.
- `db/posts.go` has `GetPostByID(db, id) (*Post, error)` — the `Post` struct includes `UserID`, `Author`, `Title`, `Type`.
- `middleware/auth.go` has `RequireAuth(db, handler)` middleware and `GetUserID(r)` to extract the authenticated user's ID from context.
- `handlers/response.go` has `writeJSON(w, status, data)` and `writeError(w, status, message)` helpers.

**Design decision:** We do NOT want to expose the author's email in the public post JSON (`GET /api/posts` or `GET /api/posts/{id}`). Instead, the author's email is only returned through a dedicated, **authenticated** contact endpoint. This gives a basic layer of privacy — only logged-in villagers can get the contact link.

## What I need you to do

Create a single new endpoint that returns a `mailto:` URL for a given post's author.

### 1. New handler — `handlers/contact.go`

Create a new file `handlers/contact.go` with a handler `GetPostContact`.

**`GetPostContact(database *sql.DB) http.HandlerFunc`:**

- **Auth required** (will be wrapped with `middleware.RequireAuth` in `main.go`).
- Extract the post `{id}` from `r.PathValue("id")`. Parse it as int64; return 400 if invalid.
- Look up the post via `db.GetPostByID(database, id)`. Return 404 if not found.
- The post type must be `"offer"` or `"request"`. If the post is an `"announcement"`, return 400 with `"contact not available for announcements"`.
- The logged-in user should not be able to contact themselves. Extract the caller's user ID via `middleware.GetUserID(r)`. If `callerID == post.UserID`, return 400 with `"cannot contact yourself"`.
- Look up the post author via `db.GetUserByID(database, post.UserID)`. Return 500 if this fails.
- Build the mailto URL:
  ```
  mailto:<author_email>?subject=Village Square: <post_title>
  ```
  The subject should be URL-encoded (use `net/url`'s `url.QueryEscape` or `url.PathEscape` for the subject value).
- Return 200 with JSON:
  ```json
  {
    "mailto": "mailto:jan@village.nl?subject=Village%20Square%3A%20Fresh%20trout%20available"
  }
  ```

Keep the handler concise — no more than ~40 lines of logic.

### 2. Register the route — `main.go`

Add the new route in `main.go`, grouped with the other post routes:

```go
mux.HandleFunc("GET /api/posts/{id}/contact", middleware.RequireAuth(database, handlers.GetPostContact(database)))
```

Place it after the existing `DELETE /api/posts/{id}` route and before the events routes.

### 3. No other changes

- Do **not** modify any existing handlers, models, or database schema.
- Do **not** modify the frontend.
- Do **not** add any new database tables.

## Acceptance criteria

- `go build` succeeds with no errors.
- `GET /api/posts/1/contact` without a session cookie returns `401`.
- `GET /api/posts/1/contact` with a valid session (where post 1 is by another user) returns:
  ```json
  { "mailto": "mailto:<author_email>?subject=Village%20Square%3A%20<url-encoded-title>" }
  ```
- `GET /api/posts/999/contact` returns 404.
- Contacting your own post returns 400.
- Contacting an announcement post returns 400.
- All existing endpoints continue to work.

## Constraints

- No external libraries — just the standard library (`net/url` for encoding).
- Follow the existing code style: one handler per file where sensible, use `writeJSON` / `writeError` helpers.
- Keep the response minimal — only the `mailto` string, no extra fields.
