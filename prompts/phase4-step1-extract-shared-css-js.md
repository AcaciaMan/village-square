# Phase 4 — Step 1: Extract Shared CSS & JS into External Files

## Context

Village Square has three HTML pages with heavily duplicated inline CSS and JS:

- `static/index.html` (312 lines) — landing/auth page
- `static/dashboard.html` (1284 lines) — feed, filters, new post modal, post detail
- `static/village-day.html` (1060 lines) — timeline, filters, new event modal

**Duplicated across dashboard.html and village-day.html:**
- ~90 lines of identical CSS: reset, header bar, nav links, banner, logout button, mobile header breakpoint
- ~50 lines of identical JS: `showBanner()`, auth guard (`fetch('/api/me')`), logout handler, banner timer logic
- Modal styles are nearly identical (overlay, card, close button, form groups, submit button)

**Duplicated across all three pages:**
- CSS reset (`*`, `box-sizing`, fonts, background)
- Banner styles (`.success`, `.error`, `fadeIn` animation)

## What I need you to do

Extract shared styles and scripts into external files. **No backend changes. No feature changes.**

### 1. Create `static/shared.css`

Extract all CSS that is shared between **at least two** of the three pages:

**From all three pages:**
- CSS reset (`*`, `*::before`, `*::after` box-sizing, margin, padding)
- `body` font family, background, color
- Banner styles (`#banner`, `#banner.success`, `#banner.error`, `@keyframes fadeIn`)

**From dashboard + village-day (app pages):**
- Header bar (`.header`, `.header-title`, `.header-nav`, `.nav-link`, `.nav-link:hover`, `.nav-link.active`)
- Header right (`.header-right`, `.header-user`, `.logout-btn`, `.logout-btn:hover`)
- Main content container (`.main`, `.welcome`)
- Modal overlay & card (`.modal-overlay`, `.modal-card`, `.modal-close`, `.form-group`, `.form-group label`, `.form-group input`, `.form-group textarea`, `.form-group select`, `.field-error`, `.form-submit-btn`, `.form-submit-btn:disabled`)
- Type/filter toggle buttons (`.type-btn`, `.type-btn.active` — shared pattern between post type and event type filters)
- Mobile breakpoint for header (`@media (max-width: 500px)`)

**Keep page-specific styles inline** — e.g., feed cards, post badges, timeline dividers, event cards stay in their respective pages. Only extract what's truly shared.

### 2. Create `static/shared.js`

Extract shared JS functions into a small utility file. Use an IIFE to avoid global pollution, exposing only what's needed via a `VS` namespace (Village Square):

```js
var VS = (function () {
  // ... private state ...
  return {
    showBanner: function (message, type) { ... },
    authGuard: function (onSuccess) { ... },
    setupLogout: function () { ... },
    timeAgo: function (dateStr) { ... },
    escapeHTML: function (str) { ... }
  };
})();
```

**Functions to extract:**

- `VS.showBanner(message, type)` — shows the `#banner` div with success/error styling, auto-hides after 4 seconds. Assumes `<div id="banner" aria-live="polite"></div>` exists.

- `VS.authGuard(onSuccess)` — calls `GET /api/me`. If 401, redirects to `/index.html`. If 200, calls `onSuccess(user)` with the user object. Used on page load.

- `VS.setupLogout()` — attaches click handler to `#logoutBtn`, calls `POST /api/logout`, redirects to `/index.html`.

- `VS.timeAgo(dateStr)` — returns a human-readable "X minutes ago" / "X hours ago" / "X days ago" string. Currently duplicated in dashboard.html JS.

- `VS.escapeHTML(str)` — escapes `<`, `>`, `&`, `"`, `'` for safe innerHTML insertion. Currently duplicated (or missing — add it if not present).

### 3. Update all three HTML files

**`dashboard.html` and `village-day.html`:**
- Add `<link rel="stylesheet" href="/shared.css">` in `<head>` before the inline `<style>` block.
- Add `<script src="/shared.js"></script>` before the page-specific `<script>` block.
- **Remove** all CSS rules that are now in `shared.css` from the inline `<style>`.
- **Remove** all JS functions that are now in `shared.js` from the inline `<script>`.
- Replace direct calls with `VS.showBanner(...)`, `VS.authGuard(...)`, `VS.setupLogout()`, `VS.timeAgo(...)`.

**`index.html`:**
- Add `<link rel="stylesheet" href="/shared.css">` in `<head>`.
- Remove duplicated CSS reset and banner styles.
- The landing page has its own auth redirect logic (redirect to dashboard if already logged in), so it doesn't use `VS.authGuard` — but it can use `VS.showBanner` and `VS.escapeHTML` if applicable.

### 4. Verify no functional changes

After refactoring:
- All three pages should look and behave **exactly** the same as before.
- CSS specificity should work correctly (inline `<style>` after external `<link>` to allow overrides).
- JS load order: `shared.js` must load before page-specific `<script>`.

## Acceptance criteria

- `static/shared.css` exists with ~100-150 lines of shared styles.
- `static/shared.js` exists with `VS` namespace exposing at least `showBanner`, `authGuard`, `setupLogout`, `timeAgo`, `escapeHTML`.
- `dashboard.html` — reduced by ~100+ lines of CSS and ~40+ lines of JS.
- `village-day.html` — reduced by ~100+ lines of CSS and ~40+ lines of JS.
- `index.html` — reduced by ~20+ lines of duplicated reset/banner CSS.
- All pages render identically to before (visual regression: none).
- All functionality works: auth guard, logout, banners, filters, modals, delete, forms.
- No JavaScript errors in the browser console on any page.

## Constraints

- No backend changes.
- No new features.
- No CSS or JS frameworks/libraries.
- Keep `shared.js` under 80 lines.
- Keep `shared.css` focused — only truly shared rules, not "everything".
- Don't break the mobile layout.
