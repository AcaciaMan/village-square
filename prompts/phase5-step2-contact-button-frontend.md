# Phase 5 — Step 2: "I'm Interested" Contact Button (Frontend)

## Context

Village Square is a community board app (Go + SQLite + vanilla HTML/CSS/JS). We just added a backend endpoint in the previous step:

```
GET /api/posts/{id}/contact   (auth required)
```

This endpoint returns a JSON object like:
```json
{ "mailto": "mailto:jan@village.nl?subject=Village%20Square%3A%20Fresh%20trout%20available" }
```

It returns:
- **401** if not authenticated.
- **400** if the post is an announcement, or if the caller is the post author.
- **404** if the post doesn't exist.

**Current post card rendering** is in `static/dashboard.html`, inside the `postCardHTML(p)` function. Each post has `p.id`, `p.user_id`, `p.type` (`offer` | `request` | `announcement`), `p.title`, `p.author`, etc. The global `currentUserID` holds the logged-in user's ID (set during auth guard). Posts already have a Delete button visible only to the post owner.

**Event listeners** are attached in `attachPostCardListeners()` which runs after every render.

**The app uses `shared.js`** which provides `VS.toast(message, type)` for notifications and `VS.inlineConfirm(el, callback)` for inline confirmations.

## What I need you to do

Add a contact button to post cards so logged-in users can reach out to the post author via email. **No backend changes in this step.**

### 1. Add the contact button to `postCardHTML()` — `static/dashboard.html`

Inside the `postCardHTML(p)` function, add a contact button with these rules:

- **Only show on `offer` and `request` posts** — never on `announcement`.
- **Only show if the logged-in user is NOT the post author** — i.e., `currentUserID && p.user_id !== currentUserID`.
- Button label depends on post type:
  - For `offer`: `"I'm interested"`
  - For `request`: `"I can help!"`
- The button should have:
  - A CSS class: `contact-btn`
  - A data attribute: `data-contact-id="<post_id>"`
- Place the button **after the post-body-full div** and **before the post-full-date div**, so it appears in the expanded card detail view.

Example HTML the function should produce (for an offer by another user):
```html
<button class="contact-btn" data-contact-id="3">I'm interested</button>
```

### 2. Wire up the click handler — `static/dashboard.html`

In `attachPostCardListeners()`, add event listeners for all `.contact-btn` elements:

```js
var contactBtns = feedContainer.querySelectorAll('.contact-btn');
for (var i = 0; i < contactBtns.length; i++) {
  contactBtns[i].addEventListener('click', handleContactClick);
}
```

Create the `handleContactClick(e)` function:

1. Get the post ID from `e.currentTarget.getAttribute('data-contact-id')`.
2. Disable the button and change its text to `"Opening…"`.
3. Fetch `GET /api/posts/<id>/contact` with `credentials: 'same-origin'`.
4. If the response is OK, parse the JSON and do `window.location.href = data.mailto` to open the user's email client.
5. Show a toast: `VS.toast('Opening your email client…', 'success')`.
6. After a short delay (500ms), re-enable the button and restore its original text.
7. On error, show a toast with the error message: `VS.toast(err.message || 'Could not get contact info.', 'error')`, re-enable the button, and restore its text.

### 3. Style the contact button — `static/dashboard.html`

Add CSS for `.contact-btn` in the existing `<style>` block, placed near the `.post-delete-btn` styles. The button should:

- Only be visible when the card is expanded: by default `display: none`, and `.post-card.expanded .contact-btn { display: inline-block; }`.
- Have an outline style consistent with the app's green accent:
  - `padding: 0.35rem 0.8rem`
  - `font-size: 0.82rem`
  - `font-weight: 600`
  - `color: #2d6a4f`
  - `background: none` → on hover: `background: #e8f5e9`
  - `border: 2px solid #2d6a4f`
  - `border-radius: 8px`
  - `cursor: pointer`
  - `margin-top: 0.6rem`
  - `transition: all 0.2s`
- Add a `disabled` state: `.contact-btn:disabled { opacity: 0.6; cursor: default; }`

### 4. No other changes

- Do **not** modify any Go backend files.
- Do **not** modify `shared.css` or `shared.js`.
- All changes are in `static/dashboard.html` only.

## Acceptance criteria

- On the dashboard, expanding an **offer** post by another user shows an **"I'm interested"** button.
- Expanding a **request** post by another user shows an **"I can help!"** button.
- **Announcement** posts never show a contact button.
- The logged-in user's **own** posts never show a contact button.
- Clicking the button briefly shows "Opening…", then opens the user's default email client with a pre-filled `To:` and `Subject:`.
- A success toast appears: "Opening your email client…".
- If the backend returns an error (e.g., post not found), an error toast appears.
- When the card is **collapsed**, the button is hidden (same as the full body text).
- The button looks consistent with the existing card design (green accent, outline style).
- All existing functionality (expand/collapse, delete, create post, filters) still works.

## Constraints

- Vanilla JS only — no frameworks, no build tools.
- Keep changes minimal and contained within `dashboard.html`.
- Follow the existing code style: `var` declarations, `for` loops, string concatenation for HTML, event delegation patterns.
- The button must be keyboard-accessible (it's a `<button>`, so this is automatic).
