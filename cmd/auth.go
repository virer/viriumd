package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the base64 part (after "Basic ")
		payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		if err != nil {
			http.Error(w, "Invalid auth header", http.StatusUnauthorized)
			return
		}

		// Split into username and password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != config.API_username || pair[1] != config.API_password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Auth passed, continue to handler
		next.ServeHTTP(w, r)
	})
}
