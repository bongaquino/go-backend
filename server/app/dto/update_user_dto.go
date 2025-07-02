package dto

type UpdateUserDTO struct {
	FirstName  string  `json:"first_name" binding:"required"`
	MiddleName *string `json:"middle_name"`
	LastName   string  `json:"last_name" binding:"required"`
	Suffix     *string `json:"suffix"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password"`
	Role       string  `json:"role" binding:"required"`
	IsVerified *bool   `json:"is_verified" binding:"required"`
	IsLocked   *bool   `json:"is_locked" binding:"required"`
	IsDeleted  *bool   `json:"is_deleted" binding:"required"`
}
