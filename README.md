# Village Square

A community web app for a rural village — connecting villagers, local producers (fishermen, farmers, crafters), and the yearly **Village Day** celebration.

Built with **Go** (stdlib only) + **SQLite** + plain **HTML/CSS/JS** — no frameworks, no build step.

Short introduction video:
https://github.com/user-attachments/assets/cd99a68b-b613-4101-8467-571a12c5924e

## Quick Start

```bash
# Prerequisites: Go 1.22+, GCC (for SQLite cgo)

# Clone and build
git clone https://github.com/<you>/village-square.git
cd village-square
go build -o village-square.exe .

# Seed with demo data (8 users, 18 posts, 6 events)
./village-square.exe --seed

# Start the server
./village-square.exe
# → http://localhost:8080
```

**Demo login:** `jan@village.nl` / `jan123` (or any seeded user — see `db/seed.go`)

## Features

### Phase 1 — Foundation & Auth
- User registration and login with bcrypt-hashed passwords
- Server-side sessions via secure HttpOnly cookies
- Auth guard on protected pages (auto-redirect if not logged in)
- `GET /api/me` returns the current user

### Phase 2 — Posts & Marketplace
- Create, view, and delete posts (offers, requests, announcements)
- Five categories: fish, produce, crafts, services, other
- Feed with type and category filters
- Post detail view, "New Post" modal with validation
- Users can only delete their own posts

### Phase 3 — Village Day
- Dedicated events system (garage sales, sports, gatherings)
- Timeline view grouped by morning / afternoon / evening
- Event creation modal with date, time, and location
- Posts can be linked to events (e.g., "Selling bikes at Village Day")
- Navigation bar across feed and Village Day pages
- Dashboard shows live preview of upcoming events

### Phase 4 — Polish & Hardening
- Shared CSS & JS extracted (`shared.css`, `shared.js`) — DRY across pages
- Toast notification system (replaces basic banners)
- Skeleton loaders during data fetching
- Inline confirm for deletions (no `window.confirm`)
- Responsive hamburger menu, full-screen mobile modals
- 44px touch targets, horizontal-scroll filter bar, print stylesheet
- Request logging middleware, 1 MB body size limit
- Custom 404 page (HTML for browsers, JSON for API)
- Automatic expired-session cleanup (hourly)

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/health` | No | Database health check |
| `POST` | `/api/register` | No | Create a new user |
| `POST` | `/api/login` | No | Log in, set session cookie |
| `POST` | `/api/logout` | No | Clear session |
| `GET` | `/api/me` | Yes | Current user profile |
| `GET` | `/api/posts` | No | List posts (`?type=`, `?category=`) |
| `GET` | `/api/posts/{id}` | No | Single post detail |
| `POST` | `/api/posts` | Yes | Create a post |
| `DELETE` | `/api/posts/{id}` | Yes | Delete own post |
| `GET` | `/api/events` | No | List events (`?type=`) |
| `GET` | `/api/events/{id}` | No | Single event detail |
| `POST` | `/api/events` | Yes | Create an event |
| `DELETE` | `/api/events/{id}` | Yes | Delete own event |

## Project Structure

```
village-square/
├── main.go                  # Entry point, routes, middleware chain
├── db/
│   ├── db.go                # SQLite init, migrations
│   ├── users.go             # User queries
│   ├── sessions.go          # Session CRUD + cleanup
│   ├── posts.go             # Post CRUD + filters
│   ├── events.go            # Event CRUD + filters
│   └── seed.go              # Demo data (--seed flag)
├── handlers/
│   ├── register.go          # POST /api/register
│   ├── login.go             # POST /api/login
│   ├── logout.go            # POST /api/logout
│   ├── me.go                # GET /api/me
│   ├── posts.go             # Post endpoints
│   ├── events.go            # Event endpoints
│   ├── health.go            # GET /api/health
│   └── response.go          # writeJSON / writeError helpers
├── middleware/
│   ├── auth.go              # RequireAuth, GetUserID
│   ├── headers.go           # SecurityHeaders (CSP, etc.)
│   ├── bodylimit.go         # 1 MB request body limit
│   └── logging.go           # Request logging
├── static/
│   ├── index.html           # Landing / register / login
│   ├── dashboard.html       # Feed, filters, new post modal
│   ├── village-day.html     # Event timeline, new event modal
│   ├── 404.html             # Custom not-found page
│   ├── shared.css           # Shared styles
│   └── shared.js            # Shared utilities (VS namespace)
├── go.mod
└── go.sum
```

## Tech Stack

- **Backend:** Go stdlib (`net/http`, `database/sql`, `crypto/rand`, `bcrypt`)
- **Database:** SQLite via `github.com/mattn/go-sqlite3` (WAL mode, foreign keys)
- **Frontend:** Vanilla HTML, CSS, JavaScript — zero dependencies
- **Auth:** bcrypt passwords + server-side session tokens in HttpOnly cookies

## License

See [LICENSE](LICENSE).
