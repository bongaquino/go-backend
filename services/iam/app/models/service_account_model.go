package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ServiceAccount struct {
	ID               string    `bson:"_id,omitempty"`
	Name             string    `bson:"name"`
	Description      string    `bson:"description"`
	ClientID         string    `bson:"client_id"`
	ClientSecretHash string    `bson:"client_secret_hash"`
	PolicyID         string    `bson:"policy_id"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

func (ServiceAccount) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "client_id", Value: 1}},
	}
}
