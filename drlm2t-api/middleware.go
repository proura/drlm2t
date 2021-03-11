package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Middleware to log requests
func middlewareLog(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.SetOutput(os.Stdout) // logs go to Stderr by default
		log.Println(r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r) // call ServeHTTP on the original handler
	})
}

// Middleware to check if the recived user token is ok
// User --> API / Web Page user
func middlewareUserToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Request Cookie "session_token"
		c, err := r.Cookie("session_token")
		if err != nil {
			// If no exist token redirec to login
			return
		}

		// First clean old sessions
		new(Session).CleanSessions()

		// Get the session from sessions whit the token value
		session := Session{"", c.Value, 0}
		session, err = session.Get()
		if err != nil {
			// If there is an error fetching from sessions, redirect to login
			return
		}

		session.Timestamp = time.Now().Unix()
		session.Update()

		// Send updated expiration time cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Path:    "/",
			Expires: time.Now().Add(600 * time.Second),
		})

		next.ServeHTTP(w, r)
	})
}
