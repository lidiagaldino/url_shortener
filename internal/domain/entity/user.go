package entity

import "time"

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}
