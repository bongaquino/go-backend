package models

type PolicyPermission struct {
	ID           string `bson:"_id,omitempty"`
	PolicyID     string `bson:"policy_id"`
	PermissionID string `bson:"permission_id"`
}
