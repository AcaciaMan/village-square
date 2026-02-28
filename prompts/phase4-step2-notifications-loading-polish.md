# Phase 4 — Step 2: Consistent Notifications, Loading States & Error Handling

## Context

Village Square now has shared CSS/JS extracted into `static/shared.css` and `static/shared.js` (with `VS.showBanner`, `VS.authGuard`, `VS.setupLogout`, `VS.timeAgo`, `VS.escapeHTML`). Three pages: `index.html`, `dashboard.html`, `village-day.html`.

## What I need you to do

Audit and polish every user-facing interaction to have consistent notifications, loading states, and error handling. **No backend changes. No new features.**

### 1. Upgrade banner to toast notification

Replace the simple banner with a more polished toast system in `shared.js` and `shared.css`:

**Visual upgrade:**
- Toast slides in from the top-right corner (not a full-width banner).
- Max-width: 360px, rounded corners, subtle shadow.
- Left accent border: 4px solid — green for success, red for error, blue for info.
- Close button (×) on the right side.
- Auto-dismiss after 4 seconds, with a small progress bar at the bottom that shrinks.
- Multiple toasts can stack (up to 3 visible, newest on top).
- Fade-out animation when dismissed.

**New API in `VS`:**
```js
VS.toast(message, type)  // type: 'success' | 'error' | 'info'
```
- Replaces `VS.showBanner`.
- Update all call sites across all three pages.

**Add a toast container** to shared markup expectations:
```html
<div id="toastContainer" aria-live="polite"></div>
```
Add this div to all three pages. The CSS in `shared.css` positions it fixed top-right.

### 2. Loading states — consistent pattern

Audit every `fetch()` call and ensure consistent loading UX:

**Feed (dashboard.html):**
- While loading: show a pulsing skeleton placeholder (3 gray rectangular blocks mimicking post cards) instead of plain "Loading…" text.
- On error: show "Could not load posts" with a "Retry" button. Clicking retries the fetch.

**Events timeline (village-day.html):**
- Same skeleton pattern (3 blocks shaped like event cards).
- Same error/retry pattern.

**Village Day preview (dashboard.html sidebar card):**
- Show "Loading…" text (small area, skeleton is overkill here).
- On error: "Could not load events."

**Modal forms (new post, new event):**
- Submit button shows a spinner or "Posting…" / "Creating…" text while the request is in flight.
- Button stays disabled until response.
- On error: button re-enables, toast shows the error.

**Delete actions:**
- Button changes text to "Deleting…" while in flight.
- On success: card fades out (0.3s CSS transition) before removal from DOM.
- On error: button re-enables, toast shows error.

### 3. Skeleton loader CSS

Add to `shared.css`:

```css
.skeleton {
  background: linear-gradient(90deg, #eee 25%, #e0e0e0 50%, #eee 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

Each page generates the appropriate skeleton HTML (3 divs sized to match their card dimensions).

### 4. Form validation error polish

Audit all form submissions (register, login, new post, new event):

- Field-level errors (`.field-error` divs) should appear with a subtle slide-down animation.
- When a field is corrected (user types in it), clear its error immediately via an `input` event listener.
- Required fields that fail validation should get a red border (`border-color: #d32f2f`), cleared on input.
- Order of validation: check all fields, show **all** errors at once (not just the first one). Update both frontend forms if they currently stop at the first error.

### 5. Confirm dialogs — replace `window.confirm`

Replace `confirm('Delete this post?')` / `confirm('Delete this event?')` with a custom inline confirmation:

- Instead of a browser dialog, transform the delete button inline:
  - First click: button changes to a two-button group: "Cancel" and "Confirm Delete" (red).
  - Confirm button triggers the actual delete.
  - Cancel button (or clicking elsewhere / 3-second timeout) reverts to the original delete button.
- This avoids jarring browser dialogs and feels more modern.

### 6. Empty states — consistent illustrations

Audit all empty states and make them consistent:

- **Feed with no posts:** "No posts yet — be the first to share something!" with a simple CSS-only illustration (or emoji: 📋).
- **Feed with no matching filter results:** "No posts match your filters." with a "Clear filters" button.
- **Village Day with no events:** "No events planned yet — be the first to add one!" with 📅 emoji.
- **Village Day preview with no events (dashboard):** "No upcoming events" with an "Add one →" link to `/village-day.html`.

All empty states should use a centered `.empty-state` class from `shared.css` with consistent font size, color, and spacing.

## Acceptance criteria

- Toast notifications appear top-right, stack, auto-dismiss with progress bar, have close button.
- All `showBanner` calls replaced with `VS.toast`.
- Skeleton loaders shown during feed and events loading.
- All modals show "Posting…" / "Creating…" spinner state on submit.
- Delete buttons use inline confirm pattern (no `window.confirm`).
- Delete cards fade out before removal.
- Form validation shows all errors at once, clears on input, highlights fields with red border.
- Empty states are consistent with emoji and action buttons/links.
- No JavaScript errors in console on any page.
- Mobile layout still works (toasts full-width on small screens).

## Constraints

- No backend changes.
- No new features — this is pure polish.
- No CSS/JS frameworks.
- Keep skeleton CSS simple (the shimmer animation + a few sized divs).
- Toasts and inline confirms should be composable utilities in `shared.js`.
