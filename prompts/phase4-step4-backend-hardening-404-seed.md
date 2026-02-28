# Phase 4 — Step 4: Backend Validation Hardening, 404 Page & Seed Refresh

## Context

Village Square frontend is now polished: shared CSS/JS, toast notifications, skeleton loaders, inline confirms, responsive hamburger menu, full-screen mobile modals, and print stylesheet.

**Current backend routes:**
```
GET  /api/health
POST /api/register
POST /api/login
POST /api/logout
GET  /api/me            (auth)
POST /api/posts          (auth)
GET  /api/posts
GET  /api/posts/{id}
DELETE /api/posts/{id}   (auth)
POST /api/events         (auth)
GET  /api/events
GET  /api/events/{id}
DELETE /api/events/{id}  (auth)
```

## What I need you to do

Backend hardening, a proper 404 page, and a seed data refresh. **No new features.**

### 1. Input sanitization — `handlers/register.go`, `handlers/posts.go`, `handlers/events.go`

Audit all handlers that accept user input and ensure:

- **Trim whitespace** on all text inputs (name, email, title, body, description, location). Some already do this for `title` — do it consistently everywhere.
- **Strip leading/trailing newlines** from body/description fields.
- **Limit request body size** — add a global middleware that wraps `r.Body` with `http.MaxBytesReader(w, r.Body, 1<<20)` (1 MB limit). This prevents abuse. Create `middleware/bodylimit.go`:
  ```go
  func LimitBody(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
          next.ServeHTTP(w, r)
      })
  }
  ```
  Apply it in `main.go` wrapping the mux alongside `SecurityHeaders`.

### 2. Improve error responses — `handlers/response.go`

- All error responses should include a consistent JSON structure. Already using `{"error":"message"}` — verify this is used everywhere (including `MaxBytesReader` errors).
- Add `405 Method Not Allowed` handling: create a small helper or ensure each endpoint returns `405` for wrong methods. Currently some endpoints use `mux.HandleFunc("POST /api/posts", ...)` which auto-rejects wrong methods — verify the mux returns JSON (not HTML) for 405. If it returns plaintext, add a catch-all handler.

### 3. Custom 404 page — `static/404.html` + Go handler

Create a friendly 404 page and wire it up:

**`static/404.html`:**
- Use shared CSS (`shared.css`).
- Simple centered card: "Page not found" title, a friendly message ("This page doesn't exist. Maybe it moved, or perhaps you mistyped the URL."), and two links:
  - "← Back to Feed" → `/dashboard.html`
  - "← Back to Home" → `/index.html`
- Style: same color scheme, simple layout.

**Wire it up in `main.go`:**
- Add a custom `NotFoundHandler` on the mux that serves `static/404.html` with a 404 status code.
- Make sure API routes (`/api/*`) return JSON 404s (not the HTML page): `{"error":"not found"}`.
- Non-API routes return the HTML 404 page.

### 4. Request logging middleware

Add basic request logging so you can see what's happening in the terminal:

Create `middleware/logging.go`:
```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Wrap ResponseWriter to capture status code
        ww := &statusWriter{ResponseWriter: w, status: 200}
        next.ServeHTTP(ww, r)
        log.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.status, time.Since(start).Round(time.Millisecond))
    })
}
```

Apply it as the outermost middleware in `main.go`:
```go
log.Fatal(http.ListenAndServe(addr, middleware.Logging(middleware.SecurityHeaders(middleware.LimitBody(mux)))))
```

This gives terminal output like:
```
GET /dashboard.html 200 2ms
POST /api/login 200 85ms
GET /api/posts 200 3ms
GET /api/posts/999 404 1ms
```

### 5. Session cleanup

Add a background goroutine that cleans up expired sessions every hour:

In `db/sessions.go`, add:
```go
func CleanExpiredSessions(db *sql.DB) (int64, error) {
    res, err := db.Exec("DELETE FROM sessions WHERE expires_at < datetime('now')")
    if err != nil {
        return 0, err
    }
    return res.RowsAffected()
}
```

In `main.go`, start a goroutine before `ListenAndServe`:
```go
go func() {
    for {
        time.Sleep(1 * time.Hour)
        n, err := db.CleanExpiredSessions(database)
        if err != nil {
            log.Printf("session cleanup error: %v", err)
        } else if n > 0 {
            log.Printf("cleaned %d expired sessions", n)
        }
    }
}()
```

### 6. Seed data refresh

Update `db/seed.go` to make the demo more compelling after all the UI polish:

**Add 2 more users:**

| Name | Email | Password |
|---|---|---|
| Dirk van Dam | dirk@village.nl | dirk123 |
| Lotte Smit | lotte@village.nl | lotte123 |

**Add 4 more posts (from the new users + existing):**

| User | Type | Title | Body | Category |
|---|---|---|---|---|
| Dirk van Dam | offer | Firewood for sale | Split oak, seasoned for 2 years. €50 per cubic metre, delivery possible. | services |
| Lotte Smit | request | Looking for babysitter | Need someone for Friday evenings, 18:00–22:00. Our kids are 4 and 7. | services |
| Dirk van Dam | announcement | New cycling path open | The path along the canal is finally finished! Great for morning rides. | other |
| Lotte Smit | offer | Homemade sourdough bread | Baking every Saturday. Reserve by Thursday evening. €4 per loaf. | produce |

**Add 1 more event:**

| User | Type | Title | Description | Location | Start | End |
|---|---|---|---|---|---|---|
| Lotte Smit | other | Village Day Volunteering | Help us set up tents and tables on Friday evening! Drinks provided. | Community hall | 2026-06-14T17:00:00Z | 2026-06-14T20:00:00Z |

**Updated seed summary:** `"Seeded 8 users, 18 posts, and 6 events."`

## Acceptance criteria

- All text inputs are trimmed consistently across all handlers.
- Request body limited to 1 MB — large payloads get rejected with 413 or 400.
- Navigating to a non-existent page (e.g., `/nonexistent`) shows a styled 404 page.
- API requests to non-existent endpoints return JSON 404.
- Terminal shows coloured-style request logs: `METHOD PATH STATUS DURATION`.
- Expired sessions are cleaned up every hour (verify the goroutine starts in logs).
- `.\village-square.exe --seed` creates 8 users, 18 posts, 6 events.
- Running `--seed` twice creates no duplicates.
- All existing functionality still works.

## Constraints

- Keep middleware composable and chainable.
- Don't over-engineer logging — no log levels, no log files, just stdout.
- The 404 page should be a real HTML file served by Go (not a redirect).
- Keep the session cleanup simple — one goroutine, one timer, no channels.
