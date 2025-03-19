package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Role struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}

func (Role) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "name", Value: 1}},
	}
}
