package service

import (
	"koneksi/server/app/repository"
)

type DirectoryService struct {
	directoryRepo *repository.DirectoryRepository
}

// NewDirectoryService initializes a new DirectoryService
func NewDirectoryService(directoryRepo *repository.DirectoryRepository) *DirectoryService {
	return &DirectoryService{
		directoryRepo: directoryRepo,
	}
}
