package main

import (
	"log"
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
}

// newHub creates a returns a new Hub instance
func newHub() *Hub {
	return &Hub{
		broadcast: make(chan Message),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
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

				// Send notice to other clients that a new client logged in
				for c := range h.clients {
					c.send <- Message{Type:"notice", Content: h.T("chat_notice_login", client.user)}
				}

				// Add the new client to the list of clients
				h.clients[client] = true

			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged out.")

					// Send notice to other clients that this client logged out
					for c := range h.clients {
						c.send <- Message{Type:"notice", Content: h.T("chat_notice_logout", client.user)}
					}

					// Remove the client from the list and close send chan
					delete(h.clients, client)
					close(client.send)
				}

			case message := <- h.broadcast:
				log.Println("[" + message.UserID + "] " + message.UserName + " sent '" + message.Content + "'.")

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
