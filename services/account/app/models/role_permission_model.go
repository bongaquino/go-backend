package models

import "go.mongodb.org/mongo-driver/bson"

type RolePermission struct {
	ID           string `bson:"_id,omitempty"`
	RoleID       string `bson:"role_id"`
	PermissionID string `bson:"permission_id"`
}

func (RolePermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
	}
}
