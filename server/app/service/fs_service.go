package service

import (
	"context"
	"errors"
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

func (fs *FSService) ReadRootDirectory(ctx context.Context, userID string) (*model.Directory,
	[]*model.Directory, []*model.File, error) {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByUserIDName(ctx, userID, "root")
	if err != nil {
		return nil, nil, nil, err
	}

	// Fetch the subdirectories within the root directory
	subDirectories, err := fs.directoryRepo.ListByDirectoryIDUserID(ctx, directory.ID.Hex(), userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// @TODO: Fetch files within the root directory

	// Return the directory details
	return directory, subDirectories, nil, nil
}

func (fs *FSService) ReadDirectory(ctx context.Context, ID string, userID string) (*model.Directory,
	[]*model.Directory, []*model.File, error) {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check if the directory exists
	if directory == nil {
		return nil, nil, nil, errors.New("directory not found")
	}

	// Fetch the subdirectories within the specified directory
	subDirectories, err := fs.directoryRepo.ListByDirectoryIDUserID(ctx, directory.ID.Hex(), userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// @TODO: Fetch files within the specified directory

	// Return the directory details
	return directory, subDirectories, nil, nil
}

func (fs *FSService) CreateDirectory(ctx context.Context, directory *model.Directory) error {
	// Create the directory in the repository
	err := fs.directoryRepo.Create(ctx, directory)
	if err != nil {
		return err
	}

	return nil
}
