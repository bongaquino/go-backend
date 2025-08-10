package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolePermission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RoleID       primitive.ObjectID `bson:"role_id"`
	PermissionID primitive.ObjectID `bson:"permission_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func (RolePermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "role_id", Value: 1}, {Key: "permission_id", Value: 1}},
	}
}
