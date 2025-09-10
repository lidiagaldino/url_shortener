package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	HashedPassword string             `bson:"password" json:"password"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}
