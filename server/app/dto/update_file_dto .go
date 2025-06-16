package dto

type UpdateFileDTO struct {
	DirectoryID *string `json:"directory_id"`
	Name        string  `json:"name"`
}
