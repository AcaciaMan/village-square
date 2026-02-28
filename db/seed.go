package db

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
}

type seedPost struct {
	UserEmail string
	Type      string
	Title     string
	Body      string
	Category  string
	EventRef  string // title of linked event, empty if none
}

type seedEvent struct {
	UserEmail   string
	EventType   string
	Title       string
	Description string
	Location    string
	StartTime   string // RFC3339
	EndTime     string // RFC3339, empty if none
}

// Seed populates the database with demo users and posts. Idempotent — skips
// users that already exist (matched by email).
func Seed(db *sql.DB) error {
	users := []seedUser{
		{"Jan Visser", "jan@village.nl", "jan123"},
		{"Maria de Boer", "maria@village.nl", "maria123"},
		{"Pieter Bakker", "pieter@village.nl", "pieter123"},
		{"Sophie Jansen", "sophie@village.nl", "sophie123"},
		{"Kees Mulder", "kees@village.nl", "kees123"},
		{"Anna de Vries", "anna@village.nl", "anna123"},
	}

	posts := []seedPost{
		{"jan@village.nl", "offer", "Fresh herring from this morning", "Caught 5kg of herring at the lake. Pick up at harbor before noon. €5/kg.", "fish", ""},
		{"maria@village.nl", "offer", "Homemade apple jam", "Made with apples from our garden. 6 jars available, €3 each.", "produce", ""},
		{"pieter@village.nl", "request", "Need help fixing garden fence", "Few panels blown over in the storm. Can anyone help this Saturday? I'll provide lunch!", "services", ""},
		{"sophie@village.nl", "offer", "Hand-knitted scarves", "Wool scarves in various colors. Perfect for the coming winter. €15 each.", "crafts", ""},
		{"kees@village.nl", "announcement", "Road closure next week", "The Dorpsstraat will be closed Mon-Wed for pipe repairs. Use the Molenweg detour.", "other", ""},
		{"anna@village.nl", "offer", "Free-range eggs", "Our chickens are laying well! Fresh eggs available daily, €2.50 per dozen.", "produce", ""},
		{"jan@village.nl", "request", "Looking for a dog sitter", "Going away for a weekend in March. Need someone to watch our labrador Rex.", "services", ""},
		{"maria@village.nl", "announcement", "Village council meeting", "Next meeting is March 5th at 19:30 in the community hall. All welcome.", "other", ""},
		{"pieter@village.nl", "offer", "Smoked mackerel", "Smoked it myself yesterday. 2kg available. €8/kg, ready to eat.", "fish", ""},
		{"sophie@village.nl", "request", "Looking for wool donations", "Starting a knitting group for teens. Any leftover yarn welcome!", "crafts", ""},
		{"kees@village.nl", "offer", "Tractor available for garden work", "Can help plough or move heavy loads this weekend. Free for neighbours.", "services", ""},
		{"anna@village.nl", "request", "Wanted: rhubarb", "Looking for rhubarb to make a pie for village day. Will trade for eggs!", "produce", ""},
		{"jan@village.nl", "offer", "Old fishing rods at Village Day", "Selling 3 fishing rods and a tackle box at my garage sale. €10-€25 each.", "fish", "Jan's Garage Sale"},
		{"anna@village.nl", "offer", "Fresh eggs at Maria's sale", "I'll have a table at Maria's garden sale with eggs and rhubarb cake!", "produce", "Maria's Garden Sale"},
	}

	events := []seedEvent{
		{"jan@village.nl", "garage_sale", "Jan's Garage Sale", "Old fishing gear, tools, and boat parts. Everything priced to go!", "Housenumber 7, driveway", "2026-06-15T09:00:00Z", "2026-06-15T12:00:00Z"},
		{"maria@village.nl", "garage_sale", "Maria's Garden Sale", "Homemade preserves, old kitchenware, and children's books.", "Housenumber 15, front garden", "2026-06-15T09:30:00Z", "2026-06-15T13:00:00Z"},
		{"pieter@village.nl", "sport", "Village Football Match", "Annual match: East Village vs West Village. All skill levels welcome!", "Sports field behind the church", "2026-06-15T14:00:00Z", "2026-06-15T16:00:00Z"},
		{"kees@village.nl", "sport", "Kids' Sack Race & Games", "Fun games for children under 12. Prizes for everyone!", "Village green", "2026-06-15T13:00:00Z", "2026-06-15T14:30:00Z"},
		{"sophie@village.nl", "gathering", "Evening BBQ & Music", "Bring your own drinks, meat provided. Live acoustic music from 20:00.", "Community hall garden", "2026-06-15T18:00:00Z", "2026-06-15T23:00:00Z"},
	}

	// Insert users (skip if email already exists).
	usersCreated := 0
	emailToID := make(map[string]int64)

	for _, u := range users {
		// Check if user exists.
		var existingID int64
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", u.Email).Scan(&existingID)
		if err == nil {
			emailToID[u.Email] = existingID
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", u.Email, err)
		}

		res, err := db.Exec(
			"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			u.Name, u.Email, string(hash),
		)
		if err != nil {
			return fmt.Errorf("insert user %s: %w", u.Email, err)
		}

		id, _ := res.LastInsertId()
		emailToID[u.Email] = id
		usersCreated++
	}

	// Insert events first (so linked posts can reference them).
	eventsCreated := 0
	for _, ev := range events {
		userID, ok := emailToID[ev.UserEmail]
		if !ok {
			continue
		}

		var exists int
		err := db.QueryRow(
			"SELECT 1 FROM events WHERE user_id = ? AND title = ?",
			userID, ev.Title,
		).Scan(&exists)
		if err == nil {
			continue // already seeded
		}

		startTime, _ := time.Parse(time.RFC3339, ev.StartTime)
		var endTime *time.Time
		if ev.EndTime != "" {
			t, _ := time.Parse(time.RFC3339, ev.EndTime)
			endTime = &t
		}

		_, err = db.Exec(
			"INSERT INTO events (user_id, title, description, event_type, location, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
			userID, ev.Title, ev.Description, ev.EventType, ev.Location, startTime, endTime,
		)
		if err != nil {
			return fmt.Errorf("insert event %q: %w", ev.Title, err)
		}
		eventsCreated++
	}

	// Insert posts (skip if already exists by user_id + title).
	postsCreated := 0
	for _, p := range posts {
		userID, ok := emailToID[p.UserEmail]
		if !ok {
			continue
		}

		// Check if this exact post already exists for this user.
		var exists int
		err := db.QueryRow(
			"SELECT 1 FROM posts WHERE user_id = ? AND title = ?",
			userID, p.Title,
		).Scan(&exists)
		if err == nil {
			continue // already seeded
		}

		// Resolve event_id if this post is linked to an event.
		var eventID *int64
		if p.EventRef != "" {
			var eid int64
			err := db.QueryRow("SELECT id FROM events WHERE title = ?", p.EventRef).Scan(&eid)
			if err == nil {
				eventID = &eid
			}
		}

		_, err = db.Exec(
			"INSERT INTO posts (user_id, type, title, body, category, event_id) VALUES (?, ?, ?, ?, ?, ?)",
			userID, p.Type, p.Title, p.Body, p.Category, eventID,
		)
		if err != nil {
			return fmt.Errorf("insert post %q: %w", p.Title, err)
		}
		postsCreated++
	}

	fmt.Printf("Seeded %d users, %d posts, and %d events.\n", usersCreated, postsCreated, eventsCreated)
	return nil
}
