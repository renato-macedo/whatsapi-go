package models

// Message is the structure that goes inside the request body
type Message struct {
	Text   string `json:"text"`
	Number string `json:"number"`
}
