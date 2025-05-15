package dto

type GenerateServiceAccountDTO struct {
	UserID       *string `json:"user_id"`
	Name         string  `json:"name" binding:"required"`
	ClientID     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
}
