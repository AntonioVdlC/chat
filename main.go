package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	db := initDB()
	defer db.Close()

	hub := newHub(db)
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
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/robots.txt")
	})

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/",
		gzipHandler(
			cacheHandler(
				http.StripPrefix("/public/", fs))))

	port := getPort()
	log.Printf("Listening on port %s", port)
	
	if env := os.Getenv("ENV"); env == "dev" {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	} else {
		err := http.ListenAndServe(port, http.HandlerFunc(redirect))
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}
