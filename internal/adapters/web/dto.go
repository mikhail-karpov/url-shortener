package web

import "time"

type ShortenURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortURLResponse struct {
	ID        string    `json:"id"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrResponse struct {
	Error string `json:"error"`
}
