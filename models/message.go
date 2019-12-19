package models

// TextMessageDTO is the structure that goes inside the request body
type TextMessageDTO struct {
	Text   string `json:"text"`
	Number string `json:"number"`
}

// ImageMessageDTO is the structure that goes inside the request body
type ImageMessageDTO struct {
	URLs   []string `json:"urls"`
	Number string   `json:"number"`
}
