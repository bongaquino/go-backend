package dto

type UpdateOrgDTO struct {
	Name     string  `json:"name" binding:"required"`
	Domain   string  `json:"domain" binding:"required"`
	Contact  string  `json:"contact" binding:"required"`
	PolicyID string  `json:"policy_id" binding:"required"`
	ParentID *string `json:"parent_id"`
}
