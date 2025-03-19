package models

type RolePermission struct {
	ID           string `bson:"_id,omitempty"`
	RoleID       string `bson:"role_id"`
	PermissionID string `bson:"permission_id"`
}
