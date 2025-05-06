package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Name                 string             `bson:"name"`
	Domain               string             `bson:"domain"`
	Contact              string             `bson:"contact"`
	SubscriptionPlanID   primitive.ObjectID `bson:"subscription_plan_id"`
	SubscriptionStatusID primitive.ObjectID `bson:"subscription_status_id"`
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
}

func (Organization) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "name", Value: 1}},
	}
}
