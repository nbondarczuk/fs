package rest

import (
	"net/http"
)

// authMiddleware checks for access token
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TBD
        next.ServeHTTP(w, r)
    })
}
