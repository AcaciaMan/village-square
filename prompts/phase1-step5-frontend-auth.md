# Phase 1 — Step 5: Frontend — Auth Flow & Dashboard

## Context

The Village Square backend is now complete for Phase 1:

- `POST /api/register` — creates user (name, email, password).
- `POST /api/login` — returns user JSON + sets `session` cookie.
- `POST /api/logout` — clears session.
- `GET /api/me` — returns current user (requires auth).
- `GET /api/health` — DB health check.

The frontend is a single `static/index.html` with a name + email form that currently only logs to the console.

## What I need you to do

Update the frontend so it actually talks to the backend. Create a second page for the dashboard. Keep everything as plain HTML + inline CSS + vanilla JS — no frameworks, no build step.

### 1. Update `static/index.html` — Landing / Auth Page

**Add a password field** to the form (between email and the button):
- `<input type="password" id="password" name="password" placeholder="Password (min 6 chars)">`
- With a label and a `.field-error` div, matching the existing style.

**Add a toggle** between Register and Login modes:
- Below the button, add a link/text: `Already have an account? Log in` / `Need an account? Register`
- Clicking toggles the form between two modes:
  - **Register mode** (default): shows Name + Email + Password fields, button says "Register". Submits to `POST /api/register`, then auto-logs in via `POST /api/login`.
  - **Login mode**: hides the Name field, button says "Log in". Submits to `POST /api/login` only.

**Form submission** (replace the existing console.log script):
1. Prevent default.
2. Validate fields client-side (same rules as before + password ≥ 6 chars).
3. `fetch()` the appropriate API endpoint with `Content-Type: application/json`.
4. On success:
   - Register mode: after register succeeds, immediately call login, then redirect to `/dashboard.html`.
   - Login mode: redirect to `/dashboard.html`.
5. On error: show the server's error message in a banner/toast above the form.

**Auto-redirect**: On page load, call `GET /api/me`. If it returns `200`, the user is already logged in — redirect to `/dashboard.html` immediately.

### 2. Create `static/dashboard.html`

A simple page with:

**Header bar** (top of page, full width):
- Left: "Village Square" title/logo text.
- Right: user's name + a "Logout" button.

**Main content** (centered):
- A welcome message: "Welcome, [Name]!" with the user's name from `/api/me`.
- A subtitle: "This is your village dashboard. More features coming soon."
- Placeholder sections (empty cards or divs with labels) for future features:
  - "Village Feed" (placeholder)
  - "Village Day" (placeholder)

**Logout button**:
- Calls `POST /api/logout`.
- On success, redirects to `/index.html`.

**Auth guard**: On page load, call `GET /api/me`. If it returns `401`, redirect to `/index.html`.

### 3. Styling guidelines

- Reuse the same color scheme: green `#2d6a4f` / `#1b4332`, background `#f5f5f0`, white cards.
- Dashboard uses a full-page layout (not the centered card), with a header bar and content below.
- Mobile-friendly: header stacks vertically on small screens, cards go full-width.
- Keep CSS inline in each HTML file (we'll extract to a shared file in Phase 4).
- Add a simple fade-in or slide animation for the success/error messages (CSS only, no JS animation library).

### 4. Error/success banner

Add a reusable banner pattern to both pages:
- A `<div id="banner">` at the top of the container, hidden by default.
- JS function `showBanner(message, type)` where type is `'success'` or `'error'`.
- Success: green background (`#d4edda`), error: red background (`#f8d7da`).
- Auto-hides after 4 seconds.

## Acceptance criteria

- **Register flow**: Fill in name + email + password → user created → auto-login → lands on dashboard showing "Welcome, [Name]!".
- **Login flow**: Toggle to login → email + password → lands on dashboard.
- **Logout**: Click logout on dashboard → cookie cleared → back on landing page.
- **Session persistence**: Refresh dashboard → still logged in (cookie + `/api/me` check).
- **Auth guard**: Navigate directly to `/dashboard.html` when not logged in → redirected to `/index.html`.
- **Auto-redirect**: Navigate to `/index.html` when already logged in → redirected to `/dashboard.html`.
- **Validation**: Empty fields or bad email show client-side errors. Server errors (e.g., duplicate email) show in the banner.
- **Mobile**: Both pages look good on a 375px-wide viewport (iPhone SE size).

## Constraints

- No JavaScript frameworks or libraries.
- No CSS frameworks (no Bootstrap, Tailwind, etc.).
- No build tools.
- Keep each HTML file self-contained (inline CSS + JS).
- Do not modify any Go files in this step.
