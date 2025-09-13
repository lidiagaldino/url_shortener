package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLStat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URLID     string             `bson:"url_id" json:"url_id"`
	ClickedAt time.Time          `bson:"clicked_at" json:"clicked_at"`
	IP        string             `bson:"ip" json:"ip"`
	UserAgent string             `bson:"user_agent" json:"user_agent"`
	Referer   string             `bson:"referer,omitempty" json:"referer,omitempty"`
}
