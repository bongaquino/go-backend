package dto

type UpdateDirectoryDTO struct {
	DirectoryID *string `json:"directory_id"`
	Name        string  `json:"name"`
}
