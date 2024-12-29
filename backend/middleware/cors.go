// middleware/cors.go
package middleware

import (
    "log"
    "net/http"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log incoming request
        log.Printf("Incoming %s request to %s", r.Method, r.URL.Path)

        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Requested-With")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight
        if r.Method == "OPTIONS" {
            log.Printf("Handling OPTIONS preflight request")
            w.WriteHeader(http.StatusOK)
            return
        }

        // Log request body size for POST requests
        if r.Method == "POST" {
            log.Printf("Content-Length: %d bytes", r.ContentLength)
        }

        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Log response
        log.Printf("Completed request to %s", r.URL.Path)
    })
}