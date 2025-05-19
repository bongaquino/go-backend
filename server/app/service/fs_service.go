package service

import (
	"context"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
)

type FSService struct {
	directoryRepo *repository.DirectoryRepository
	fileRepo      *repository.FileRepository
}

// NewFSService initializes a new FSService
func NewFSService(directoryRepo *repository.DirectoryRepository, fileRepo *repository.FileRepository) *FSService {
	return &FSService{
		directoryRepo: directoryRepo,
		fileRepo:      fileRepo,
	}
}

func (fs *FSService) ReadRootDirectory(ctx context.Context, userID string) (*model.Directory, error) {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByUserIDName(ctx, userID, "root")
	if err != nil {
		return nil, err
	}

	// Fetch the files and directories within the root directory

	// Return the directory details
	return directory, nil
}
