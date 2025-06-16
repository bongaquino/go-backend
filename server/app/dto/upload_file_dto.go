package dto

type UploadFileDTO struct {
	DirectoryID string `json:"directory_id" form:"directory_id" binding:"omitempty"`
}
