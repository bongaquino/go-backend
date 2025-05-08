package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationUserRole struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	UserID         primitive.ObjectID `bson:"user_id"`
	RoleID         primitive.ObjectID `bson:"access_id"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

func (OrganizationUserRole) GetIndexes() []primitive.D {
	return nil
}
