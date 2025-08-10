package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Limit struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	UserID         primitive.ObjectID  `bson:"user_id"`
	OrganizationID *primitive.ObjectID `bson:"organization_id"`
	BytesLimit     int64               `bson:"bytes_limit"`
	BytesUsage     int64               `bson:"bytes_usage"`
	CreatedAt      time.Time           `bson:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at"`
}

func (Limit) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}},
	}
}
