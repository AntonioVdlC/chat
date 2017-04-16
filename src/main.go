package main

import (
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var port = getPort()

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SESSION_SECRET")))
	store.MaxLength(math.MaxInt64)

	gothic.Store = store

	goth.UseProviders(facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost" + port + "/auth/callback?provider=facebook"))
}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/auth/callback", authCallback)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/auth", auth)

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &Client{
		hub: hub,
		conn: conn,
		send: make(chan Message),
	}
	client.hub.register <- client

	go client.write()
	client.read()
}

func getUser(r *http.Request, p string) (goth.User, error) {
	session, _ := gothic.Store.Get(r, p + gothic.SessionName)
	values := session.Values[p]
	if values == nil {
		return goth.User{}, errors.New("cannot find session values")
	}
	
	provider, _ := goth.GetProvider(p)
	sess, _ := provider.UnmarshalSession(values.(string))
	user, err := provider.FetchUser(sess)

	if err != nil {
		return goth.User{}, err
	}

	return user, nil
}

func home(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r, "facebook")

	if err != nil {
		t, _ := template.ParseFiles("./tpl/login.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("./tpl/chat.html")
		t.Execute(w, user)
	}
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	_, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func auth(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
