package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
)

// Hub represents a chat room and holds a list of all connected
// clients. Besides, it has methods to broadcast messages, and
// register and unregister clients 
type Hub struct {
	clients map[*Client]bool
	broadcast chan Message
	register chan *Client
	unregister chan *Client
	T i18n.TranslateFunc
	db *sql.DB
}

// newHub creates a returns a new Hub instance
func newHub(db *sql.DB) *Hub {
	return &Hub{
		broadcast: make(chan Message),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
		db: db,
	}
}

// setT adds the i18n TranslateFunc utility to the Hub
func (h *Hub) setT(T i18n.TranslateFunc) {
	h.T = T
}

// run launches the Hub instance!
func (h *Hub) run() {
	for {
		select {
			case client := <- h.register:
				log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged in.")

				// Create a login message
				message := Message{
					UserID: client.user.UserID,
					UserName: client.user.Name,
					UserAvatar: client.user.AvatarURL,
					Type: "login",
					Content: h.T("chat_notice_login", client.user),
					Date: time.Now().UTC(),
				}

				// Insert login in DB
				id, err := insertMessage(h.db, message)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
				message.ID = id

				// Send notice to other clients that a new client logged in
				for c := range h.clients {
					c.send <- message
				}

				// Add the new client to the list of clients
				h.clients[client] = true

				// Send previous messages to new client
				rows, err := selectPreviousMessage(h.db, client.user.UserID)
				defer rows.Close()
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

				for rows.Next() {
					var message Message
					if err := rows.Scan(&message.ID, &message.UserID, &message.UserName, &message.UserAvatar, &message.Type, &message.Content, &message.Date); err != nil {
						log.Printf("Error: %v", err)
						return
					}
					client.send <- message
				}

			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged out.")

					// Create logout message
					message := Message{
						UserID: client.user.UserID,
						UserName: client.user.Name,
						UserAvatar: client.user.AvatarURL,
						Type:"logout",
						Content: h.T("chat_notice_logout", client.user),
						Date: time.Now().UTC(),
					}

					// Insert logout in DB
					id, err := insertMessage(h.db, message)
					if err != nil {
						log.Printf("Error: %v", err)
						return
					}
					message.ID = id

					// Send notice to other clients that this client logged out
					for c := range h.clients {
						c.send <- message
					}

					// Remove the client from the list and close send chan
					delete(h.clients, client)
					close(client.send)
				}

			case message := <- h.broadcast:
				log.Println("[" + message.UserID + "] " + message.UserName + " sent '" + message.Content + "'.")

				// Insert message in DB
				id, err := insertMessage(h.db, message)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
				message.ID = id

				// Send message to all clients
				for client := range h.clients {
					select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
					}
				}
		}
	}
}
