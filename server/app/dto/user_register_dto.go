package dto

type RegisterUserRequest struct {
	FirstName       string  `json:"first_name" binding:"required"`
	MiddleName      *string `json:"middle_name"`
	LastName        string  `json:"last_name" binding:"required"`
	Suffix          *string `json:"suffix"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=8"`
	ConfirmPassword string  `json:"confirm_password" binding:"required,eqfield=Password"`
}
