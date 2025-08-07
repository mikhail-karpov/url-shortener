package web

import "time"

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	OriginalURL string    `json:"original_url"`
	Alias       string    `json:"alias"`
	CreatedAt   time.Time `json:"created_at"`
}
