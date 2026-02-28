# Phase 3 — Step 3: Frontend — Village Day Page & Navigation

## Context

Village Square backend now supports events:

- `POST /api/events` — create event (auth required). Fields: title, description, event_type (garage_sale|sport|gathering|other), location, start_time, end_time.
- `GET /api/events` — list all events sorted by start_time ASC, with optional `?type=` filter. Returns `{"events": [...], "count": N}`.
- `GET /api/events/{id}` — single event detail (public).
- `DELETE /api/events/{id}` — delete own event (auth required).
- `GET /api/me` — current user.

The dashboard (`static/dashboard.html`) currently has a "Village Day" placeholder card in the sidebar that says "Upcoming village events and activities will be listed here."

## What I need you to do

Create a new `static/village-day.html` page and add a navigation bar across all pages. **No backend changes.**

### 1. Add a navigation bar to all pages

Add a simple nav bar **inside** the existing green header on both `dashboard.html` and the new `village-day.html`:

```html
<div class="header">
  <span class="header-title">Village Square</span>
  <nav class="header-nav">
    <a href="/dashboard.html" class="nav-link">Feed</a>
    <a href="/village-day.html" class="nav-link">Village Day</a>
  </nav>
  <div class="header-right">
    <span class="header-user" id="headerUser"></span>
    <button class="logout-btn" id="logoutBtn">Logout</button>
  </div>
</div>
```

- The **active page** link gets a `class="nav-link active"` with a bottom border or underline to indicate the current page.
- On mobile, the nav links sit between the title and user area (the header already wraps with flexbox).
- Style the nav links as white text, no underline, with a subtle hover effect.

### 2. Create `static/village-day.html`

A self-contained page (inline CSS + JS, same pattern as dashboard) with:

#### Header
- Same header bar as dashboard, with nav links. "Village Day" link is active.
- Same auth guard: on load, `fetch('/api/me')` — redirect to `/index.html` if 401.
- Same logout button logic.
- Same banner component for success/error messages.

#### Page title section
- Title: **"Village Day"**
- Subtitle: "Our annual village celebration — garage sales, sports, and evening gathering."

#### Event type filter tabs
- Horizontal tab/button row similar to the feed type filters:
  - `All`, `Garage Sales`, `Sports`, `Gatherings`, `Other`
- Clicking re-fetches events with the `?type=` parameter.

#### Timeline view
- Fetch `GET /api/events` and group events by **time of day** based on `start_time`:
  - **Morning** (before 12:00): label "☀️ Morning"
  - **Afternoon** (12:00–17:00): label "🌤️ Afternoon"
  - **Evening** (after 17:00): label "🌙 Evening"
- Each time-of-day section is a visual group with a label/divider.
- If a section has no events, don't show it (no empty sections).

#### Event cards
Each event renders as a card showing:
- **Event type badge** — colored pill, similar to post type badges:
  - Garage sale: `background: #fff3cd; color: #856404;` (amber)
  - Sport: `background: #cce5ff; color: #004085;` (blue)
  - Gathering: `background: #d4edda; color: #155724;` (green)
  - Other: `background: #e2e3e5; color: #383d41;` (gray)
- **Title** — bold
- **Time** — formatted as "09:00 – 13:00" (or "09:00 onwards" if no end_time)
- **Location** — with a 📍 icon prefix (if non-empty)
- **Description** — full text (events are short, no need to truncate)
- **Author** — small text: "Posted by Jan Visser"
- **Delete button** — only shown if `user_id` matches current user, same confirm pattern as posts

#### "Add Event" button
- A prominent button (top of page or floating) labeled "+ Add Event"
- Opens a modal form (same modal pattern as the new-post modal on dashboard):

**Modal — New Event form:**
- **Title** — text input, required, max 200 chars
- **Event type** — radio-style toggle buttons: Garage Sale, Sport, Gathering, Other. Default: Other.
- **Location** — text input, optional. Placeholder: "e.g., Housenumber 12, garden"
- **Date** — `<input type="date">`, required
- **Start time** — `<input type="time">`, required
- **End time** — `<input type="time">`, optional
- **Description** — textarea, optional, max 2000 chars
- **Submit button** — "Create Event"

**Form behavior:**
1. Client-side validation (all rules matching the backend).
2. Combine date + start_time into an ISO string for `start_time` field (e.g., `2026-06-15T09:00:00Z`). Same for end_time if provided.
3. `fetch('POST /api/events', ...)` with JSON body.
4. On success: close modal, show banner "Event created!", re-fetch and re-render the timeline.
5. On error: show error in banner or modal.
6. Disable button during flight.

#### Empty state
- If no events exist: "No events planned yet — be the first to add one!"

### 3. Update `dashboard.html` — Village Day card

Replace the static "Village Day" placeholder card with a **live preview**:
- Fetch `GET /api/events` (limit display to the first 3 upcoming events).
- Show each as a mini card: event type badge + title + time.
- Add a "View all →" link that navigates to `/village-day.html`.
- If no events: "No upcoming events" with an "Add one →" link.

### 4. Styling guidelines

- Reuse the same color scheme and card styles from the dashboard.
- Timeline section dividers: a horizontal rule with the time-of-day label centered (like `——— ☀️ Morning ———`).
- Event cards: slightly different layout from post cards to accommodate time + location.
- Modal matches the new-post modal style.
- `<input type="date">` and `<input type="time">` should be styled to match existing inputs.
- Mobile: single column, cards full width, modal full-screen on small viewports.

## Acceptance criteria

- Navigation bar appears on both dashboard and village-day pages with correct active state.
- Village Day page loads, fetches events from API, and groups them by morning/afternoon/evening.
- Filter tabs work (Garage Sales shows only garage_sale events, etc.).
- "Add Event" button opens a modal, submitting creates an event and refreshes the timeline.
- Delete button shown only on own events, with confirmation dialog.
- Dashboard's Village Day card shows a live preview of up to 3 upcoming events.
- Auth guard works: navigating to `/village-day.html` without login redirects to `/index.html`.
- Mobile layout works on 375px viewport.
- All existing dashboard functionality (feed, new post, filters) still works.

## Constraints

- No JavaScript frameworks or libraries.
- No backend changes.
- Each HTML file is self-contained (inline CSS + inline JS).
- Keep the nav bar markup/style consistent between the two pages (copy-paste is fine for now — Phase 4 will extract shared CSS).
