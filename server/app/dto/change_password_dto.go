package dto

type ChangePasswordDTO struct {
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
}
