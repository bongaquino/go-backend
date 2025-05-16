package service

import (
	"koneksi/server/app/repository"
)

type FileService struct {
	fileRepo *repository.FileRepository
}

// NewFileService initializes a new FileService
func NewFileService(fileRepo *repository.FileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}
