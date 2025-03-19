package models

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
