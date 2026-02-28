package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Serve everything in ./static at the root path.
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	// Wrap all routes with security-header middleware.
	addr := ":8080"
	log.Printf("Village Square listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.SecurityHeaders(mux)))
}
