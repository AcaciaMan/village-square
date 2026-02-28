package db

import (
	"database/sql"
	"fmt"

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
		{"jan@village.nl", "offer", "Fresh herring from this morning", "Caught 5kg of herring at the lake. Pick up at harbor before noon. €5/kg.", "fish"},
		{"maria@village.nl", "offer", "Homemade apple jam", "Made with apples from our garden. 6 jars available, €3 each.", "produce"},
		{"pieter@village.nl", "request", "Need help fixing garden fence", "Few panels blown over in the storm. Can anyone help this Saturday? I'll provide lunch!", "services"},
		{"sophie@village.nl", "offer", "Hand-knitted scarves", "Wool scarves in various colors. Perfect for the coming winter. €15 each.", "crafts"},
		{"kees@village.nl", "announcement", "Road closure next week", "The Dorpsstraat will be closed Mon-Wed for pipe repairs. Use the Molenweg detour.", "other"},
		{"anna@village.nl", "offer", "Free-range eggs", "Our chickens are laying well! Fresh eggs available daily, €2.50 per dozen.", "produce"},
		{"jan@village.nl", "request", "Looking for a dog sitter", "Going away for a weekend in March. Need someone to watch our labrador Rex.", "services"},
		{"maria@village.nl", "announcement", "Village council meeting", "Next meeting is March 5th at 19:30 in the community hall. All welcome.", "other"},
		{"pieter@village.nl", "offer", "Smoked mackerel", "Smoked it myself yesterday. 2kg available. €8/kg, ready to eat.", "fish"},
		{"sophie@village.nl", "request", "Looking for wool donations", "Starting a knitting group for teens. Any leftover yarn welcome!", "crafts"},
		{"kees@village.nl", "offer", "Tractor available for garden work", "Can help plough or move heavy loads this weekend. Free for neighbours.", "services"},
		{"anna@village.nl", "request", "Wanted: rhubarb", "Looking for rhubarb to make a pie for village day. Will trade for eggs!", "produce"},
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

	// Insert posts (only if the user was just created — to keep idempotent).
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

		_, err = db.Exec(
			"INSERT INTO posts (user_id, type, title, body, category) VALUES (?, ?, ?, ?, ?)",
			userID, p.Type, p.Title, p.Body, p.Category,
		)
		if err != nil {
			return fmt.Errorf("insert post %q: %w", p.Title, err)
		}
		postsCreated++
	}

	fmt.Printf("Seeded %d users and %d posts.\n", usersCreated, postsCreated)
	return nil
}
