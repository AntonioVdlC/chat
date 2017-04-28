package main

import (
	"log"
)

// Hub represents a chat room and holds a list of all connected
// clients. Besides, it has methods to broadcast messages, and
// register and unregister clients 
type Hub struct {
	clients map[*Client]bool
	broadcast chan Message
	register chan *Client
	unregister chan *Client
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

// run launches the Hub instance!
func (h *Hub) run() {
	for {
		select {
			case client := <- h.register:
				log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged in.")
				for c := range h.clients {
					c.send <- Message{Type:"notice", Content: client.user.Name + " logged in."}
				}
				h.clients[client] = true
			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged out.")
					for c := range h.clients {
						c.send <- Message{Type:"notice", Content: client.user.Name + " logged out."}
					}
					delete(h.clients, client)
					close(client.send)
				}
			case message := <- h.broadcast:
				log.Println("[" + message.UserID + "] " + message.UserName + " sent '" + message.Content + "'.")
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
