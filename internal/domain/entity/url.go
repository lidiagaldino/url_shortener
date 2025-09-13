package entity

import "time"

type URL struct {
	ID          string
	OriginalURL string
	OwnerID     string
	CreatedAt   time.Time
	ClickCount  int
	LastClick   time.Time
}
