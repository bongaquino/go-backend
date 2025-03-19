package models

import "go.mongodb.org/mongo-driver/bson"

type UserRole struct {
	ID     string `bson:"_id,omitempty"`
	UserID string `bson:"user_id"`
	RoleID string `bson:"role_id"`
}

func (RolePermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "role_id", Value: 1}, {Key: "permission_id", Value: 1}},
	}
}
