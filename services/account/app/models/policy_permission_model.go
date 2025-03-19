package models

import "go.mongodb.org/mongo-driver/bson"

type PolicyPermission struct {
	ID           string `bson:"_id,omitempty"`
	PolicyID     string `bson:"policy_id"`
	PermissionID string `bson:"permission_id"`
}

func (PolicyPermission) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "policy_id", Value: 1}, {Key: "permission_id", Value: 1}},
	}
}
