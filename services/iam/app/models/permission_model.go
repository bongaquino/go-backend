package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Permission struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}

func (Permission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "name", Value: 1}},
	}
}
