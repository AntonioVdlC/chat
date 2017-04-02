package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message emitted by a client and broadcasted to the channel
type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Content string `json:"content"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))

	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	port := getPort()
	log.Printf("Listening on port %s", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8000"
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <- broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
