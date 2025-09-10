package entity

import "time"

type URL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
}
