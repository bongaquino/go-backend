package dto

type UploadFileDTO struct {
	DirectoryID string `json:"directory_id" form:"directory_id" binding:"omitempty"`
	Passphrase  string `json:"passphrase" form:"passphrase" binding:"omitempty"`
}
