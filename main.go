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
	"github.com/nicksnyder/go-i18n/i18n"
	"io/ioutil"
)

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SESSION_SECRET")))
	store.MaxLength(math.MaxInt64)

	gothic.Store = store

	host := getHost()
	goth.UseProviders(facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), host + "/auth/callback?provider=facebook"))
}

func main() {
	hub := newHub()
	go hub.run()

	loadLocales()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/auth/callback", authCallback)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/auth", auth)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	port := getPort()
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// loadLocales loads all i18n strings from the `locales` directory
func loadLocales() {
	files, err := ioutil.ReadDir("locales")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	for _, file := range files {
		i18n.MustLoadTranslationFile("locales/" + file.Name())
	}
}

// getPort returns the port by first looking at any environment variable
// nammed PORT and then defaulting to :8000
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8000"
}

// getHost returns the host URL by looking if the app is running in "dev" or
// in production (on Heroku for now)
func getHost() string {
	if env := os.Getenv("ENV"); env == "dev" {
		return "http://localhost:8000"
	}
	return "https://efrei-int-chat.herokuapp.com"
}

// getUser returns the goth.User linked with the current session
// NB: for now we are only using Facebook as an OAuth provider
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
	acceptLang := r.Header.Get("Accept-Language")
	defaultLang := "en"
	T := i18n.MustTfunc(acceptLang, defaultLang)

	user, err := getUser(r, "facebook")

	if err != nil {
		t, _ := template.New("login.html").Funcs(template.FuncMap{"T": T}).ParseFiles("./templates/login.html")
		t.Execute(w, nil)
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
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// logout handles the ... logout!
// It doesn't seem to really work at the moment, but well, maybe one day right?
func logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// auth handles the authentication through Facebook and then redirects to the
// home handler
func auth(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
