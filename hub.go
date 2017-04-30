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
	insert := `
		INSERT INTO messages (
			user_id,
			user_name,
			user_avatar,
			type,
			content,
			date_post
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	for {
		select {
			case client := <- h.register:
				log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged in.")

				// Insert login in DB
				var id string
				err := h.db.QueryRow(insert, client.user.UserID, client.user.Name, client.user.AvatarURL, "notice", "login", time.Now()).Scan(&id)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

				// Send notice to other clients that a new client logged in
				message := Message{ID:id, Type:"notice", Content: h.T("chat_notice_login", client.user)}
				for c := range h.clients {
					c.send <- message
				}

				// Add the new client to the list of clients
				h.clients[client] = true

			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged out.")

					// Insert logout in DB
					var id string
					err := h.db.QueryRow(insert, client.user.UserID, client.user.Name, client.user.AvatarURL, "notice", "logout", time.Now()).Scan(&id)
					if err != nil {
						log.Printf("Error: %v", err)
						return
					}

					// Send notice to other clients that this client logged out
					message := Message{ID:id, Type:"notice", Content: h.T("chat_notice_logout", client.user)}
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
				var id string
				err := h.db.QueryRow(insert, message.UserID, message.UserName, message.UserAvatar, message.Type, message.Content, time.Now()).Scan(&id)
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
