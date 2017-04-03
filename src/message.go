package main

// Message emitted by a client and broadcasted to the channel
type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Content string `json:"content"`
}
