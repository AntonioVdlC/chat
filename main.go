package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "dev"
	password = "dev"
	dbname = "chat_dev"
	sslmode = "disable"
)

func initDB() *sql.DB {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	log.Println("Connected to DB.")

	return db
}

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

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	port := getPort()
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
