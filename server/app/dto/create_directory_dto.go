package dto

type CreateDirectoryDTO struct {
	DirectoryID string `json:"directory_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
}
