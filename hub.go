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
	request chan Message
	register chan *Client
	unregister chan *Client
	T i18n.TranslateFunc
	db *sql.DB
}

// newHub creates a returns a new Hub instance
func newHub(db *sql.DB) *Hub {
	return &Hub{
		broadcast: make(chan Message),
		request: make(chan Message),
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
				// Send notice to other clients that a new client logged in
				message, err := loginMessage(h.db, client, h.T)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
				for c := range h.clients {
					c.send <- message
				}

				// Add the new client to the list of clients
				h.clients[client] = true

				// Send previous messages and logged-in users to new client
				bootstrap, err := bootstrap(h.db, client.user.UserID)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
				client.send <- bootstrap

			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					log.Println("[" + client.user.UserID + "] " + client.user.Name + " logged out.")

					// Create logout message
					// Send notice to other clients that this client logged out
					message, err := logoutMessage(h.db, client, h.T)
					if err != nil {
						log.Printf("Error: %v", err)
						return
					}
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
			
			case request := <- h.request:
				log.Println("[" + request.UserID + "] " + request.UserName + " requested '" + request.Content + "'.")

				messages, err := getOlderMessages(h.db, request.Date)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

				request.Client.send <- messages
		}
	}
}

// bootstrap retrieves the initial data needed to bootstrap the app
// retrieving the previous messages and the connected users.
func bootstrap(db *sql.DB, userID string) (Bootstrap, error) {
	bootstrap := Bootstrap{ []Message{}, []User{}, "bootstrap" }

	// Select previous messages from DB
	rows, err := selectPreviousMessage(db, userID)
	defer rows.Close()
	if err != nil {
		return bootstrap, err
	}
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.UserName, &message.UserAvatar, &message.Type, &message.Content, &message.Date); err != nil {
			return bootstrap, err
		}
		bootstrap.Messages = append(bootstrap.Messages, message)
	}

	// Select connected users from DB
	rows, err = selectConnectedUsers(db, userID)
	if err != nil {
		return bootstrap, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Avatar); err != nil {
			return bootstrap, err
		}
		bootstrap.Users = append(bootstrap.Users, user)
	}
	
	// All good to go!
	return bootstrap, nil
}

// loginMessage creates a message on user login and saves it to the DB
func loginMessage(db *sql.DB, client *Client, T i18n.TranslateFunc) (Message, error) {
	message := Message{
		UserID: client.user.UserID,
		UserName: client.user.Name,
		UserAvatar: client.user.AvatarURL,
		Type: "login",
		Content: T("chat_notice_login", client.user),
		Date: time.Now().UTC(),
	}

	// Insert login in DB
	id, err := insertMessage(db, message)
	if err != nil {
		return message, err
	}
	message.ID = id

	return message, nil
}

// logoutMessage creates a message on user logout and saves it to the DB
func logoutMessage(db *sql.DB, client *Client, T i18n.TranslateFunc) (Message, error) {
	message := Message{
		UserID: client.user.UserID,
		UserName: client.user.Name,
		UserAvatar: client.user.AvatarURL,
		Type: "logout",
		Content: T("chat_notice_logout", client.user),
		Date: time.Now().UTC(),
	}

	// Insert login in DB
	id, err := insertMessage(db, message)
	if err != nil {
		return message, err
	}
	message.ID = id

	return message, nil
}

// getOlderMessages retrieves the previous 10 message from a given date
func getOlderMessages(db *sql.DB, date time.Time) (Messages, error) {
	messages := Messages{ []Message{}, "messages" }

	// Select previous messages from DB
	rows, err := selectOlderMessages(db, date)
	defer rows.Close()
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.UserName, &message.UserAvatar, &message.Type, &message.Content, &message.Date); err != nil {
			return messages, err
		}
		messages.Messages = append(messages.Messages, message)
	}

	return messages, nil
}

