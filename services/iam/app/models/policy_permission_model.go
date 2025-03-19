package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PolicyPermission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	PolicyID     primitive.ObjectID `bson:"policy_id"`
	PermissionID primitive.ObjectID `bson:"permission_id"`
}

func (PolicyPermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "policy_id", Value: 1}, {Key: "permission_id", Value: 1}},
	}
}
