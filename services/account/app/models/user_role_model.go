package models

type UserRole struct {
	ID     string `bson:"_id,omitempty"`
	UserID string `bson:"user_id"`
	RoleID string `bson:"role_id"`
}
