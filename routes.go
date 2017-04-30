package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/markbates/goth/gothic"
)

// serveWs upgrades the HTTP connection to a WebSocket and registers
// a Client (and basically just starts the whole thing!)
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	user, err := getUser(r, "facebook")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	client := &Client{
		hub: hub,
		conn: conn,
		send: make(chan Message),
		user: user,
	}
	client.hub.register <- client

	go client.write()
	client.read()
}

// home handles the / route
// It looks first if there's a session by getting the User then either
// displays the login or the chat screen.
func home(w http.ResponseWriter, r *http.Request) {
	_, err := getUser(r, "facebook")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/chat", http.StatusTemporaryRedirect)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	T := initT(r.Header.Get("Accept-Language"), "en")

	_, err := getUser(r, "facebook")

	if err != nil {
		t, _ := template.New("login.html").Funcs(template.FuncMap{"T": T}).ParseFiles("./templates/login.html")
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/chat", http.StatusTemporaryRedirect)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	T := initT(r.Header.Get("Accept-Language"), "en")

	user, err := getUser(r, "facebook")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		t, _ := template.New("chat.html").Funcs(template.FuncMap{"T": T}).ParseFiles("./templates/chat.html")
		t.Execute(w, user)
	}
}

// authCallback handles the callback (duh!) and always redirects to the
// home handler
func authCallback(w http.ResponseWriter, r *http.Request) {
	_, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// logout handles the ... logout!
// It doesn't seem to really work at the moment, but well, maybe one day right?
func logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// auth handles the authentication through Facebook and then redirects to the
// home handler
func auth(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
