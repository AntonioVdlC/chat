package main

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
	"time"
)

const (
	pongWait = 60 * time.Second
	pingPeriod = 33 * time.Second
	maxMessageSize = 512
)

// WSMessage is an interface that needs to be implemented to be able to
// be sent down the pipe through the Client's send channel
type WSMessage interface {
	CanBeSentDownThePipe()
}

// Message emitted by a client and broadcasted to the channel
type Message struct {
	ID string `json:"id"`
	UserID string `json:"userId"`
	UserName string `json:"userName"`
	UserAvatar string `json:"avatar"`
	Type string `json:"type"`
	Content string `json:"content"`
	Date time.Time `json:"date"`
	Client *Client `json:"-"`
}
// CanBeSentDownThePipe allows Message to be sent through
// the Client's send channel
func (m Message) CanBeSentDownThePipe() {}

// Messages is a JSON sent by the hub when a new client connects
type Messages struct {
	Messages []Message `json:"messages"`
	Type string `json:"type"`
}
// CanBeSentDownThePipe allows Messages to be sent through
// the Client's send channel
func (m Messages) CanBeSentDownThePipe() {}

// User represents a chat user
type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}

// Bootstrap is a JSON sent by the hub when a new client connects
type Bootstrap struct {
	Messages []Message `json:"messages"`
	Users []User `json:"users"`
	Type string `json:"type"`
}
// CanBeSentDownThePipe allows Bootstrap to be sent through
// the Client's send channel
func (b Bootstrap) CanBeSentDownThePipe() {}

// Client is a middleman between the WebSocket connection and the Hub
type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan WSMessage
	user goth.User
}

// read pumps messages from the WebSocket to the Hub
func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			break
		}
		msg.UserID = c.user.UserID
		msg.UserName = c.user.Name

		switch (msg.Type) {
			case "message":
				msg.UserAvatar = c.user.AvatarURL
				msg.Date = time.Now().UTC()

				c.hub.broadcast <- msg

			case "request":
				msg.Client = c
				c.hub.request <- msg

			default:
				log.Printf("Error: unrecognised message %v", msg)
		}
	}
}

// write pumps messages from the Hub to the WebSocket
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	
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
