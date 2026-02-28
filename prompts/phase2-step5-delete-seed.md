# Phase 2 — Step 5: Delete Own Posts & Seed Data

## Context

Village Square now has a working marketplace:

- Users can register, log in, and stay authenticated via session cookies.
- `POST /api/posts` creates a post (auth required).
- `GET /api/posts` lists posts with type/category filters.
- `GET /api/posts/{id}` returns a single post.
- Dashboard shows a live feed with filters, a "New Post" modal form, and expandable/detail view for posts.

What's missing: users can't delete their own posts, and there's no sample data to demo with.

## What I need you to do

Add a delete endpoint, wire it into the frontend, and create a seed command for demo data.

### 1. Add to `db/posts.go`

Export `DeletePost(db *sql.DB, postID, userID int64) error`:

- `DELETE FROM posts WHERE id = ? AND user_id = ?`
- Use `result.RowsAffected()` to check if a row was actually deleted.
- If no rows affected, return a sentinel error (e.g., `ErrNotFound` or `sql.ErrNoRows`) — the post either doesn't exist or doesn't belong to this user.

### 2. Add to `handlers/posts.go`

Export `DeletePost(db *sql.DB) http.HandlerFunc` for `DELETE /api/posts/{id}`:

- **Auth required.**
- Parse `{id}` from the URL path.
- Get user ID from context.
- Call `db.DeletePost(db, postID, userID)`.
- If not found / not owner → `404 {"error":"post not found"}`.
  - Intentionally don't distinguish "not found" from "not yours" (prevent enumeration).
- On success → `200 {"message":"post deleted"}`.

### 3. Update `main.go`

Add the route (protected):
```go
mux.HandleFunc("DELETE /api/posts/{id}", middleware.RequireAuth(database, handlers.DeletePost(database)))
```

### 4. Update `static/dashboard.html` — delete button

- On each post card in the feed, **if the post's `user_id` matches the logged-in user's ID**, show a small delete button (trash icon or "Delete" text, red/subtle).
- Store the current user's ID from the `/api/me` response for comparison.
- Clicking delete:
  1. Show a `confirm('Delete this post?')` dialog.
  2. If confirmed, `fetch('DELETE /api/posts/{id}', ...)`.
  3. On success: remove the card from the DOM, show success banner "Post deleted."
  4. On error: show error banner.

### 5. Create `cmd/seed.go` — seed command

Create a seed mechanism so I can populate the database with demo data. Choose **one** of these approaches (pick whichever is simpler):

**Option A — SQL seed file:**
- Create `db/seed.sql` with INSERT statements.
- Add a function `db.RunSeed(db *sql.DB, filePath string) error` that reads and executes the file.
- In `main.go`, check for a `--seed` flag: `go run . --seed` runs the seed and exits.

**Option B — Go seed function:**
- Create `db/seed.go` with `Seed(db *sql.DB) error`.
- Uses Go code to insert users and posts (bcrypt-hash passwords in Go).
- In `main.go`, check for a `--seed` flag.

**Seed data to include (use option B so passwords are properly hashed):**

**Users** (6):
| Name | Email | Password |
|---|---|---|
| Jan Visser | jan@village.nl | jan123 |
| Maria de Boer | maria@village.nl | maria123 |
| Pieter Bakker | pieter@village.nl | pieter123 |
| Sophie Jansen | sophie@village.nl | sophie123 |
| Kees Mulder | kees@village.nl | kees123 |
| Anna de Vries | anna@village.nl | anna123 |

**Posts** (12, spread across users, types, and categories):

1. Jan — offer / fish: "Fresh herring from this morning" / "Caught 5kg of herring at the lake. Pick up at harbor before noon. €5/kg."
2. Maria — offer / produce: "Homemade apple jam" / "Made with apples from our garden. 6 jars available, €3 each."
3. Pieter — request / services: "Need help fixing garden fence" / "Few panels blown over in the storm. Can anyone help this Saturday? I'll provide lunch!"
4. Sophie — offer / crafts: "Hand-knitted scarves" / "Wool scarves in various colors. Perfect for the coming winter. €15 each."
5. Kees — announcement / other: "Road closure next week" / "The Dorpsstraat will be closed Mon-Wed for pipe repairs. Use the Molenweg detour."
6. Anna — offer / produce: "Free-range eggs" / "Our chickens are laying well! Fresh eggs available daily, €2.50 per dozen."
7. Jan — request / services: "Looking for a dog sitter" / "Going away for a weekend in March. Need someone to watch our labrador Rex."
8. Maria — announcement / other: "Village council meeting" / "Next meeting is March 5th at 19:30 in the community hall. All welcome."
9. Pieter — offer / fish: "Smoked mackerel" / "Smoked it myself yesterday. 2kg available. €8/kg, ready to eat."
10. Sophie — request / crafts: "Looking for wool donations" / "Starting a knitting group for teens. Any leftover yarn welcome!"
11. Kees — offer / services: "Tractor available for garden work" / "Can help plough or move heavy loads this weekend. Free for neighbours."
12. Anna — request / produce: "Wanted: rhubarb" / "Looking for rhubarb to make a pie for village day. Will trade for eggs!"

**Seed behavior:**
- Check if users already exist (by email) before inserting — make it idempotent.
- Print what was created to stdout: "Seeded 6 users and 12 posts."
- Running `.\village-square.exe --seed` seeds and exits. Running without `--seed` starts the server normally.

## Acceptance criteria

- `DELETE /api/posts/{id}` with auth → deletes own post, returns `200`.
- Attempting to delete another user's post → `404`.
- Attempting to delete without auth → `401`.
- Dashboard shows delete button only on the current user's posts.
- Delete confirmation dialog prevents accidental deletion.
- `.\village-square.exe --seed` populates 6 users and 12 posts, then exits.
- After seeding, the feed shows a realistic village marketplace with mixed post types and categories.
- Running `--seed` twice doesn't create duplicates.
- All existing functionality still works.

## Constraints

- No new frontend pages.
- Keep the seed function/data in the `db` package (or a `cmd` file).
- Passwords in seed data are simple (for demo only) — but still bcrypt-hashed.
