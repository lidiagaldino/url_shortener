package model

import "time"

type URL struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	OriginalURL string    `bson:"original_url" json:"original_url"`
	OwnerID     string    `bson:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`

	ClickCount int       `bson:"click_count" json:"click_count"`
	LastClick  time.Time `bson:"last_click,omitempty" json:"last_click,omitempty"`
}
