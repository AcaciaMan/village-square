package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve everything in ./static at the root path.
	// index.html is served automatically for "/".
	fs := http.FileServer(http.Dir("static"))

	// Wrap with a small middleware to set security headers.
	http.Handle("/", securityHeaders(fs))

	addr := ":8080"
	log.Printf("Village Square listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// securityHeaders adds baseline security headers to every response.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
