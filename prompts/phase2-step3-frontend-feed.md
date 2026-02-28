# Phase 2 — Step 3: Frontend — Village Feed Page

## Context

Village Square backend now supports:

- `GET /api/posts` — list all posts (public), with optional `?type=` and `?category=` filters. Returns `{"posts": [...], "count": N}`.
- `GET /api/posts/{id}` — single post detail (public).
- `POST /api/posts` — create post (auth required).
- `GET /api/me` — current user (auth required).

The dashboard (`static/dashboard.html`) has two placeholder cards: "Village Feed" and "Village Day".

## What I need you to do

Replace the "Village Feed" placeholder card on the dashboard with a live post feed. **Do not create a separate page** — the feed lives directly on the dashboard.

### 1. Update `static/dashboard.html`

**Replace the "Village Feed" placeholder card** with a full feed section:

#### Filter bar
- A horizontal row of filter controls above the post list:
  - **Type filter**: three toggle buttons — `All`, `Offers`, `Requests`, `Announcements`. Clicking one highlights it and re-fetches with `?type=`. `All` clears the type filter.
  - **Category dropdown**: a `<select>` with options: `All Categories`, `Fish`, `Produce`, `Crafts`, `Services`, `Other`. Changing it re-fetches with `?category=`.
- Filters combine: selecting "Offers" + "Fish" fetches `?type=offer&category=fish`.

#### Post list
- Fetch `GET /api/posts` (with current filters) and render a list of post cards.
- Each post card shows:
  - **Type badge** — small colored pill: green for offer, blue for request, orange for announcement.
  - **Title** — bold, clickable (for now just highlights/selects, detail view comes in next step).
  - **Category** — small gray tag.
  - **Author name** and **time ago** (e.g., "Jan · 2 hours ago"). Write a simple `timeAgo(date)` JS function.
  - **Body preview** — first 120 characters of the body, with "…" if truncated.
- If no posts match, show an empty state: "No posts yet — be the first to share something!"
- Posts are listed newest first (the API already returns them in this order).

#### New Post button
- A floating or prominent "New Post" button (e.g., bottom-right FAB or a button above the feed).
- Clicking it **for now** just shows a `window.alert('New post form coming soon!')`. The actual form is built in the next step.

### 2. Feed loading

- On page load (after the `/api/me` auth check succeeds), immediately fetch and render the feed.
- While loading, show a simple "Loading…" text in the feed area.
- On fetch error, show "Could not load posts. Please try again." with a retry button.

### 3. Styling

- Post cards: white background, subtle shadow, rounded corners (match existing card style).
- Type badges:
  - Offer: `background: #d4edda; color: #155724;`
  - Request: `background: #cce5ff; color: #004085;`
  - Announcement: `background: #fff3cd; color: #856404;`
- Filter toggle buttons: outlined style by default, filled when active.
- Category dropdown: matches the input style from the landing page.
- The feed section should take the full width of the `.main` container.
- Keep the "Village Day" placeholder card below the feed (or in a sidebar on wide screens — your choice, but keep it simple).

### 4. Layout adjustment

- The `.cards` grid should now have the feed as the primary content area.
- On mobile: feed takes full width, Village Day card stacks below.
- On desktop (>800px): optionally put Village Day as a smaller sidebar card, or just keep it stacked below the feed. Keep it simple.

## Acceptance criteria

- Dashboard loads and shows the post feed fetched from the API.
- Filter buttons work: clicking "Offers" shows only offers, selecting "Fish" category shows only fish, combining works.
- "All" type button clears the type filter; "All Categories" clears the category filter.
- Post cards show type badge (colored), title, category tag, author, time ago, and body preview.
- Empty state shown when no posts match filters.
- Loading state shown while fetching.
- "New Post" button is visible and shows an alert when clicked.
- Mobile layout is clean on 375px viewport.
- The Village Day placeholder card is still visible.

## Constraints

- No JavaScript frameworks or libraries.
- No new HTML files — everything goes in `dashboard.html`.
- No Go backend changes.
- Keep CSS inline in the HTML file.
