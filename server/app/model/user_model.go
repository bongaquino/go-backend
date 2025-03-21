package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	IsVerified bool               `bson:"is_verified"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (User) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "email", Value: 1}},
	}
}
