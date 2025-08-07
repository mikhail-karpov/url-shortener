package domain

import "time"

type ShortURL struct {
	OriginalURL string
	Alias       string
	CreatedAt   time.Time
}
