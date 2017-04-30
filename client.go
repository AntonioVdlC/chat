package main

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
	"time"
)

// Message emitted by a client and broadcasted to the channel
type Message struct {
	ID string `json:"id"`
	UserID string `json:"userId"`
	UserName string `json:"userName"`
	UserAvatar string `json:"avatar"`
	Type string `json:"type"`
	Content string `json:"content"`
	Date time.Time `json:"date"`
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
		msg.Date = time.Now()

		c.hub.broadcast <- msg
	}
}

// write pumps messages from the Hub to the WebSocket
func (c *Client) write() {
	ticker := time.NewTicker(33 * time.Second)
	
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
			case msg := <- c.send:
				err := c.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

			case <- ticker.C:
				if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
		}
	}
}
