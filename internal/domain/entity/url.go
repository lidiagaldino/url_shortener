package entity

import "time"

type URL struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	OriginalURL string    `bson:"original_url" json:"original_url"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
