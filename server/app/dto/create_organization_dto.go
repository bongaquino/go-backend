package dto

type CreateOrgDTO struct {
	Name     string  `bson:"name" binding:"required"`
	Domain   string  `bson:"domain" binding:"required"`
	Contact  string  `bson:"contact" binding:"required"`
	PolicyID string  `bson:"policy_id" binding:"required"`
	ParentID *string `bson:"parent_id"`
}
