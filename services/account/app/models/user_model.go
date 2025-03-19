package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (User) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "email", Value: 1}},
	}
}
