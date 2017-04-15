package main

// Message emitted by a client and broadcasted to the channel
type Message struct {
	User string `json:"user"`
	Content string `json:"content"`
}
