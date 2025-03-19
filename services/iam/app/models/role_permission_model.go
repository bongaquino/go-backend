package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolePermission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RoleID       primitive.ObjectID `bson:"role_id"`
	PermissionID string             `bson:"permission_id"`
}

func (RolePermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
	}
}
