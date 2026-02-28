# Phase 4 — Step 3: Responsive Polish & Mobile UX

## Context

Village Square has three pages (`index.html`, `dashboard.html`, `village-day.html`) with shared CSS/JS, toast notifications, skeleton loaders, inline confirm dialogs, and consistent error handling from Steps 1–2.

## What I need you to do

Do a thorough responsive audit and fix every mobile issue. Also add a few micro-interactions that make the app feel polished. Test target: **375px wide** (iPhone SE) and **768px** (tablet). **No backend changes.**

### 1. Mobile navigation — hamburger menu

The header currently uses `flex-wrap` to stack items on small screens. Replace this with a proper hamburger toggle:

**Markup** (update header in `dashboard.html` and `village-day.html` — keep in sync):
```html
<div class="header">
  <span class="header-title">Village Square</span>
  <button class="hamburger" id="hamburgerBtn" aria-label="Toggle menu" aria-expanded="false">
    <span></span><span></span><span></span>
  </button>
  <div class="header-menu" id="headerMenu">
    <nav class="header-nav">
      <a href="/dashboard.html" class="nav-link active">Feed</a>
      <a href="/village-day.html" class="nav-link">Village Day</a>
    </nav>
    <div class="header-right">
      <span class="header-user" id="headerUser"></span>
      <button class="logout-btn" id="logoutBtn">Logout</button>
    </div>
  </div>
</div>
```

**Behavior:**
- On screens **> 600px**: hamburger is hidden, menu is always visible (horizontal flexbox).
- On screens **≤ 600px**: hamburger shows, menu is initially hidden. Clicking hamburger toggles a vertical dropdown menu below the header. Add `aria-expanded` toggle.
- Hamburger animates to an × when open (three spans rotate — pure CSS).
- Menu slides down with a CSS transition (`max-height` or `transform`).
- Clicking a nav link or anywhere outside the menu closes it.

**Add to `shared.css` and `shared.js`** since both app pages use the same header.

### 2. Modal — full-screen on mobile

Currently modals are centered cards. On mobile (≤ 500px):

- Modal card should go **full-screen**: `width: 100%; height: 100%; border-radius: 0; max-width: none;`
- Close button becomes a "← Back" text button in the top-left (instead of × in top-right).
- Form fields take full width with comfortable touch-sized tap targets (min 44px height).
- Submit button sticks to the bottom of the viewport (`position: sticky; bottom: 0;`).

### 3. Touch-friendly tap targets

Audit all interactive elements and ensure:

- All buttons: minimum 44×44px touch area (padding, not just text size).
- Filter toggle buttons: at least 40px tall, with enough gap between them (no accidental taps).
- Category dropdown: at least 44px tall.
- Post/event card click areas: generous padding.
- Delete button: not too small and not too close to other actions.

### 4. Feed filter bar — horizontal scroll on mobile

On narrow screens, the type filter buttons + category dropdown may overflow:

- Wrap the filter bar in a horizontally-scrollable container: `overflow-x: auto; -webkit-overflow-scrolling: touch; white-space: nowrap;`
- Hide the scrollbar visually: `::-webkit-scrollbar { display: none; }` + `scrollbar-width: none;`
- Add subtle fade/gradient on the right edge to hint at scrollability.

### 5. Card layout improvements

**Post cards (dashboard):**
- On mobile: full-width, reduced padding (1rem instead of 1.5rem).
- Type badge and category tag on the same row, wrapping if needed.
- Author + time-ago on a single line, separated by a dot (·).
- Body preview: max 2 lines with CSS line-clamp (`-webkit-line-clamp: 2`).

**Event cards (village-day):**
- On mobile: full-width, time/location/description stack vertically.
- Type badge stays inline with the title.

### 6. Landing page responsive

The `index.html` landing page already uses `max-width: 420px` centered card. Check:
- On very small screens (<360px): card should go full-width with `margin: 0; border-radius: 0;`
- Form toggle text ("Already have an account? Log in") should be tappable and clearly styled.

### 7. Smooth micro-interactions

Add these small CSS transitions for polish:

- **Card hover** (desktop only — use `@media (hover: hover)`): subtle lift effect (`transform: translateY(-2px)`, `box-shadow` increase) on post/event cards.
- **Button press**: slight scale-down (`transform: scale(0.97)`) on `:active` for all buttons.
- **Page transitions**: body fade-in on load (0.3s opacity 0→1).
- **Input focus**: the green focus ring already exists — verify it works on all inputs including date/time pickers.

### 8. Print-friendly (quick win)

Add a `@media print` block to `shared.css`:
- Hide header, nav, buttons, modals.
- Show feed/events content only.
- Remove shadows and backgrounds.
- This lets villagers print the feed or event schedule.

## Acceptance criteria

- At 375px width: hamburger menu works, modals go full-screen, filter bar scrolls, cards are full-width.
- At 768px: layout looks good as a tablet view.
- At 1200px+: desktop layout with card hover effects.
- Hamburger menu toggles with animation, closes on outside click.
- All tap targets meet 44px minimum.
- No horizontal scrollbar on any page at any width (except the intentional filter bar scroll).
- Body text is readable (no text smaller than 14px on mobile).
- `Ctrl+P` on any page produces a clean printable layout.
- All existing functionality still works on desktop and mobile.
- No JavaScript errors on any screen size.

## Constraints

- No backend changes.
- No CSS or JS frameworks.
- Use only CSS media queries for responsive design.
- Hamburger menu is pure CSS + vanilla JS (no library).
- Test with browser DevTools device mode — no need for actual devices.
