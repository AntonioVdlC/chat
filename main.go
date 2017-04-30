package main

import (
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()

	loadSession()
	loadLocales()

	http.HandleFunc("/", home)

	http.HandleFunc("/login", login)
	http.HandleFunc("/auth/callback", authCallback)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/auth", auth)

	http.HandleFunc("/chat", chat)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.setT(initT(r.Header.Get("Accept-Language"), "en"))
		serveWs(hub, w, r)
	})

	http.HandleFunc("/service-worker.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/service-worker.js")
	})

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	port := getPort()
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
