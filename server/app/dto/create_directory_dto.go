package dto

type CreateDirectoryDTO struct {
	DirectoryID string `json:"directory_id" binding:"omitempty"`
	Name        string `json:"name" binding:"required"`
}
