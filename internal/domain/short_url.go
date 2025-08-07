package domain

import "time"

type ShortURL struct {
	ID        string
	LongURL   string
	CreatedAt time.Time
}
