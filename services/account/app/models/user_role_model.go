package models

import "go.mongodb.org/mongo-driver/bson"

type UserRole struct {
	ID     string `bson:"_id,omitempty"`
	UserID string `bson:"user_id"`
	RoleID string `bson:"role_id"`
}

func (UserRole) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
	}
}
