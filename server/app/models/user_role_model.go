package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	RoleID    primitive.ObjectID `bson:"role_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (UserRole) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
	}
}
