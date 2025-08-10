package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceAccount struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         primitive.ObjectID `bson:"user_id"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	Name           string             `bson:"name"`
	ClientID       string             `bson:"client_id"`
	ClientSecret   string             `bson:"client_secret"`
	PolicyID       primitive.ObjectID `bson:"policy_id"`
	LastUsedAt     time.Time          `bson:"last_used_at"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

func (ServiceAccount) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "client_id", Value: 1}},
	}
}
