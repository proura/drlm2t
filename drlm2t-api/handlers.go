package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Serve static content
func staticGet(w http.ResponseWriter, r *http.Request) {

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	log.Println("StaticGet: " + path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(configDRLM2T.WwwPath+"", path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the file, return a 500 internal server error and stop
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.ServeFile(w, r, path)
}

// func login(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.Redirect(w, r, "/", 302)
// 	} else {
// 		http.ServeFile(w, r, configDRLM2T.Path+"/signin.html")
// 	}
// }

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func userSignin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if configDRLM2T.APIPasswd != creds.Password {
		log.Println("Failed login for user: ", creds.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.New().String()
	sessions = append(sessions, Session{creds.Username, sessionToken, time.Now().Unix()})

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 600 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: time.Now().Add(600 * time.Second),
	})

	w.WriteHeader(http.StatusOK)

	// Browsers return "Fetch failed loading" if response is empty
	fmt.Fprintf(w, "ok")
}

// Remove ser session from sessions slice and send delete cookie
func userLogout(w http.ResponseWriter, r *http.Request) {
	// Get Request Cookie "session_token"
	c, err := r.Cookie("session_token")
	if err != nil {
		// If no exist token return
		return
	}

	// Get the session from sessions whit the token value
	session := Session{"", c.Value, 0}
	session, err = session.Get()
	if err != nil {
		// If there is an error fetching from sessions return
		return
	}

	session.Delete()

	// Send updated expiration time cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)

	// Browsers return "Fetch failed loading" if response is empty
	fmt.Fprintf(w, "ok")
}

func checkSession(w http.ResponseWriter, r *http.Request) {
	// Get Request Cookie "session_token"
	c, err := r.Cookie("session_token")
	if err != nil {
		// If no exist token return
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// First clean old sessions
	new(Session).CleanSessions()

	// Get the session from sessions whit the token value
	session := Session{"", c.Value, 0}
	session, err = session.Get()
	if err != nil {
		// If there is an error fetching from sessions return
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// session.Timestamp = time.Now().Unix()
	// session.Update()

	// // Send updated expiration time cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   session.Token,
	// 	Path:    "/",
	// 	Expires: time.Now().Add(600 * time.Second),
	// })

	w.WriteHeader(http.StatusOK)

	// Browsers return "Fetch failed loading" if response is empty
	fmt.Fprintf(w, "ok")
}
