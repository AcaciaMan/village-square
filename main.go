package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"village-square/db"
	"village-square/handlers"
	"village-square/middleware"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with demo data and exit")
	flag.Parse()

	// Initialise the SQLite database (creates the file on first run).
	database, err := db.Init("village-square.db")
	if err != nil {
		log.Fatalf("database init: %v", err)
	}
	defer database.Close()

	if *seedFlag {
		if err := db.Seed(database); err != nil {
			fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Register routes.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.Health(database))
	mux.HandleFunc("/api/register", handlers.Register(database))
	mux.HandleFunc("/api/login", handlers.Login(database))
	mux.HandleFunc("/api/logout", handlers.Logout(database))
	mux.HandleFunc("GET /api/me", middleware.RequireAuth(database, handlers.Me(database)))
	mux.HandleFunc("POST /api/posts", middleware.RequireAuth(database, handlers.CreatePost(database)))
	mux.HandleFunc("GET /api/posts", handlers.ListPosts(database))
	mux.HandleFunc("GET /api/posts/{id}", handlers.GetPost(database))
	mux.HandleFunc("DELETE /api/posts/{id}", middleware.RequireAuth(database, handlers.DeletePost(database)))
	mux.HandleFunc("GET /api/posts/{id}/contact", middleware.RequireAuth(database, handlers.GetPostContact(database)))
	mux.HandleFunc("POST /api/events", middleware.RequireAuth(database, handlers.CreateEvent(database)))
	mux.HandleFunc("GET /api/events", handlers.ListEvents(database))
	mux.HandleFunc("GET /api/events/{id}", handlers.GetEvent(database))
	mux.HandleFunc("DELETE /api/events/{id}", middleware.RequireAuth(database, handlers.DeleteEvent(database)))

	// Catch-all for unmatched /api/* routes — return JSON errors.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"not found"}`)
	})

	// Serve everything in ./static at the root path.
	fs := http.FileServer(http.Dir("static"))

	// Custom handler: serve static files, or 404 page for non-existent paths.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Let the root "/" pass through to index.html via FileServer.
		if r.URL.Path == "/" {
			fs.ServeHTTP(w, r)
			return
		}
		// Check if the requested file exists in ./static.
		path := "static" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "static/404.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	// Start background session cleanup goroutine.
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			n, err := db.CleanExpiredSessions(database)
			if err != nil {
				log.Printf("session cleanup error: %v", err)
			} else if n > 0 {
				log.Printf("cleaned %d expired sessions", n)
			}
		}
	}()

	// Wrap all routes with middleware (outermost → innermost).
	addr := ":8080"
	log.Printf("Village Square listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.Logging(middleware.SecurityHeaders(middleware.LimitBody(mux)))))
}
