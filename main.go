package main

import (
	"log"
	"net/http"

	"village-square/db"
	"village-square/handlers"
	"village-square/middleware"
)

func main() {
	// Initialise the SQLite database (creates the file on first run).
	database, err := db.Init("village-square.db")
	if err != nil {
		log.Fatalf("database init: %v", err)
	}
	defer database.Close()

	// Register routes.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.Health(database))
	mux.HandleFunc("/api/register", handlers.Register(database))
	mux.HandleFunc("/api/login", handlers.Login(database))
	mux.HandleFunc("/api/logout", handlers.Logout(database))
	mux.HandleFunc("GET /api/me", middleware.RequireAuth(database, handlers.Me(database)))

	// Serve everything in ./static at the root path.
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	// Wrap all routes with security-header middleware.
	addr := ":8080"
	log.Printf("Village Square listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.SecurityHeaders(mux)))
}
