# Phase 2 — Step 4: New Post Form & Post Detail View

## Context

The dashboard now shows a live feed of posts fetched from `GET /api/posts`, with type and category filters. There's a "New Post" button that currently shows an alert. Each post card shows a type badge, title, category, author, time ago, and body preview.

Backend endpoints available:
- `POST /api/posts` — create post (auth required). Expects `{type, title, body, category}`.
- `GET /api/posts/{id}` — single post detail (public).
- `GET /api/me` — current user info.

## What I need you to do

Implement the "New Post" form and a post detail modal/view, both on the dashboard page. **No new HTML files, no backend changes.**

### 1. New Post form — modal overlay

Replace the `window.alert` on the "New Post" button with an actual modal:

**Modal overlay:**
- A semi-transparent dark backdrop covering the page.
- A centered white card (similar to the landing page container style, max-width ~480px).
- A close button (×) in the top right corner. Clicking backdrop also closes.
- Pressing `Escape` key closes the modal.

**Form fields:**
- **Type** — three radio-style toggle buttons (not a dropdown): `Offer`, `Request`, `Announcement`. Default: `Offer`.
- **Title** — text input, required, max 200 chars. Placeholder: "What are you sharing or looking for?"
- **Category** — `<select>` dropdown: Fish, Produce, Crafts, Services, Other. Default: Other.
- **Body** — `<textarea>`, optional, max 2000 chars, 4 rows. Placeholder: "Add details…"
- **Submit button** — "Post". Green, full-width, matches existing button style.

**Behavior:**
1. Client-side validation (same rules as the backend): type required, title required/max 200, body max 2000.
2. Show per-field errors below each field (like the landing page pattern).
3. On submit: `fetch('POST /api/posts', ...)` with JSON body.
4. On success:
   - Close the modal.
   - Show a success banner: "Post published!"
   - Prepend the new post to the feed list (without a full re-fetch — just insert the returned post object at the top).
5. On error: show the server error message in a banner or in the modal.
6. Disable the submit button while the request is in flight (prevent double-submit).

### 2. Post detail view — expanding card or modal

When a user clicks a post card title in the feed:

**Option A (recommended — inline expand):**
- The card expands to show the full body text (instead of the 120-char preview).
- Show the full creation date (formatted nicely, e.g., "28 February 2026, 14:30").
- A "Close" or "Collapse" link to go back to the preview.

**Option B (modal — if you prefer):**
- Open a modal similar to the new-post modal, but read-only, showing the full post.

Pick whichever is simpler. The key requirement is that the user can read the full post body.

### 3. Styling

- Modal overlay: `background: rgba(0,0,0,0.4)`, `z-index: 100`, centered with flexbox.
- Modal card: same shadow/border-radius as existing cards, `padding: 2rem`.
- Type toggle buttons in the form: similar to the filter buttons on the feed, with the selected one filled.
- Textarea: match the existing input style, but taller.
- Transitions: fade in the modal overlay (CSS animation, 0.2s).

### 4. Accessibility

- Modal should trap focus when open (Tab cycles through modal fields only).
- Close button and submit button are keyboard-accessible.
- `aria-modal="true"` and `role="dialog"` on the modal container.
- `aria-label` on the close button.

## Acceptance criteria

- Clicking "New Post" opens a modal with the form.
- Filling in the form and submitting creates a post via the API.
- The new post appears at the top of the feed immediately (no page reload).
- Success banner shows after posting.
- Validation errors shown inline for missing/invalid fields.
- Clicking a post title in the feed shows the full post body and details.
- Modal closes on ×, Escape, or backdrop click.
- Double-submit is prevented.
- Works on mobile (375px viewport).

## Constraints

- No JavaScript frameworks.
- No new HTML files — everything in `dashboard.html`.
- No backend changes.
- Keep the modal markup at the bottom of `<body>`, hidden by default, shown/hidden via JS.
