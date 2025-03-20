package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	FirstName  string             `bson:"first_name"`
	MiddleName *string            `bson:"middle_name,omitempty"`
	LastName   string             `bson:"last_name"`
	Suffix     *string            `bson:"suffix,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (Profile) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}},
	}
}
