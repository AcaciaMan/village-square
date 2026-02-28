# Phase 3 — Step 4: Link Posts to Events & Seed Village Day Data

## Context

Village Square now has a fully functional Village Day page:

- `events` table with event_type, location, start_time, end_time.
- CRUD API for events (create, list, get, delete).
- Frontend: Village Day page with timeline view, filters, event creation modal, delete button.
- Dashboard shows a live preview of upcoming events in the Village Day card.
- Nav bar on both dashboard and village-day pages.

**What's missing:** Posts can't be linked to events (e.g., "Selling old bikes at Village Day"), and there's no seed data for events.

## What I need you to do

Add an optional `event_id` foreign key to posts, and seed the database with Village Day demo events.

### 1. Add `event_id` column to posts — update `db/db.go` migrate()

Add an `ALTER TABLE` migration (safe with `IF NOT EXISTS`-style check).

Since SQLite doesn't support `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`, use a pattern like:
```go
// Add event_id column to posts if it doesn't exist yet.
_, err := db.Exec("ALTER TABLE posts ADD COLUMN event_id INTEGER REFERENCES events(id) ON DELETE SET NULL")
// Ignore "duplicate column name" error — means migration already ran.
```

The column is nullable (a post doesn't have to be linked to an event).

### 2. Update `db/posts.go`

- Add `EventID *int64 `json:"event_id"`` to the `Post` struct (pointer for nullable).
- Add `EventTitle *string `json:"event_title"`` — populated via LEFT JOIN to events, for display convenience.
- Update **all** SQL queries (`CreatePost`, `GetPostByID`, `ListPosts`) to:
  - SELECT `p.event_id` and `e.title AS event_title`
  - Use `LEFT JOIN events e ON e.id = p.event_id`
  - Scan into the new nullable fields.

### 3. Update `handlers/posts.go` — CreatePost

Accept an optional `event_id` in the JSON request body:
```json
{
  "type": "offer",
  "title": "Selling old bikes at Village Day",
  "body": "3 bikes, all working condition. €20-€50.",
  "category": "other",
  "event_id": 1
}
```

- If `event_id` is provided (non-zero), validate that the event exists: `db.GetEventByID(db, eventID)`. If not found → `400 {"error":"event not found"}`.
- If `event_id` is 0 or omitted, leave it as NULL.

### 4. Update `db.CreatePost` signature

Change to accept an optional event ID:
```go
func CreatePost(db *sql.DB, userID int64, postType, title, body, category string, eventID *int64) (*Post, error)
```

Update the INSERT to include `event_id`.

### 5. Update the frontend — dashboard.html

In the **New Post modal**, add an optional "Link to Event" section:

- Below the category dropdown, add a `<select>` labeled "Link to Village Day event (optional)":
  - Default option: "None"
  - Populated dynamically by fetching `GET /api/events` when the modal opens.
  - Each option shows the event title.
- If an event is selected, include `event_id` in the POST body.

In the **post cards** on the feed:
- If a post has `event_title` (non-null), show a small linked badge: "📅 [Event Title]" below the category tag.
- Make it a link to `/village-day.html` (no deep-linking to a specific event needed for now).

### 6. Update `db/seed.go` — add Village Day events

Add seed events to the `Seed` function. Use the same idempotent pattern (check before insert).

**Seed events (5):**

| Created by | Type | Title | Description | Location | Start | End |
|---|---|---|---|---|---|---|
| Jan Visser | garage_sale | Jan's Garage Sale | Old fishing gear, tools, and boat parts. Everything priced to go! | Housenumber 7, driveway | 2026-06-15T09:00:00Z | 2026-06-15T12:00:00Z |
| Maria de Boer | garage_sale | Maria's Garden Sale | Homemade preserves, old kitchenware, and children's books. | Housenumber 15, front garden | 2026-06-15T09:30:00Z | 2026-06-15T13:00:00Z |
| Pieter Bakker | sport | Village Football Match | Annual match: East Village vs West Village. All skill levels welcome! | Sports field behind the church | 2026-06-15T14:00:00Z | 2026-06-15T16:00:00Z |
| Kees Mulder | sport | Kids' Sack Race & Games | Fun games for children under 12. Prizes for everyone! | Village green | 2026-06-15T13:00:00Z | 2026-06-15T14:30:00Z |
| Sophie Jansen | gathering | Evening BBQ & Music | Bring your own drinks, meat provided. Live acoustic music from 20:00. | Community hall garden | 2026-06-15T18:00:00Z | 2026-06-15T23:00:00Z |

**Add 2 posts linked to events:**

| User | Post Type | Title | Body | Category | Linked Event |
|---|---|---|---|---|---|
| Jan Visser | offer | Old fishing rods at Village Day | Selling 3 fishing rods and a tackle box at my garage sale. €10-€25 each. | fish | Jan's Garage Sale |
| Anna de Vries | offer | Fresh eggs at Maria's sale | I'll have a table at Maria's garden sale with eggs and rhubarb cake! | produce | Maria's Garden Sale |

**Seed behavior:**
- Check events by title + user_id before inserting (idempotent).
- For linked posts, first look up the event ID by title, then insert the post with that event_id.
- Update the final print: `"Seeded X users, X posts, and X events."`

### 7. Run `--seed` update

Make sure `.\village-square.exe --seed` now also seeds events (and linked posts). Running it twice should produce no duplicates.

## Acceptance criteria

- `posts` table now has a nullable `event_id` column.
- `POST /api/posts` with `event_id` → creates a post linked to that event.
- `POST /api/posts` with invalid `event_id` → `400`.
- `GET /api/posts` and `GET /api/posts/{id}` return `event_id` and `event_title` (null when not linked).
- New Post modal on dashboard has an "Link to event" dropdown.
- Post cards show event link badge when linked.
- `.\village-square.exe --seed` creates 5 events + 2 linked posts + existing 6 users + 12 posts.
- Running `--seed` twice creates no duplicates.
- All existing functionality still works.

## Constraints

- Keep the ALTER TABLE migration safe for existing databases.
- Use LEFT JOIN (not INNER JOIN) so posts without events still appear.
- Keep the event dropdown in the modal simple — fetch events once when modal opens.
