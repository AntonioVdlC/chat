package main

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
)

// Message emitted by a client and broadcasted to the channel
type Message struct {
	UserID string `json:"id"`
	UserName string `json:"user"`
	UserAvatar string `json:"avatar"`
	Content string `json:"content"`
}

// Client is a middleman between the WebSocket connection and the Hub
type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan Message
	user goth.User
}

// read pumps messages from the WebSocket to the Hub
func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			break
		}

		msg.UserID = c.user.UserID
		msg.UserName = c.user.Name
		msg.UserAvatar = c.user.AvatarURL

		c.hub.broadcast <- msg
	}
}

// write pumps messages from the Hub to the WebSocket
func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		msg := <- c.send
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
	}
}
