package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolePermission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RoleID       primitive.ObjectID `bson:"role_id"`
	PermissionID string             `bson:"permission_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func (RolePermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
	}
}
