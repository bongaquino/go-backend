package service

import (
	"context"
	"errors"
	"koneksi/server/app/dto"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"

	"go.mongodb.org/mongo-driver/bson"
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

func (fs *FSService) UpdateDirectory(ctx context.Context, ID string, userID string, request *dto.UpdateDirectoryDTO) error {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Check if the directory exists
	if directory == nil {
		return errors.New("directory not found")
	}

	// Update the parent directory if provided
	if *request.DirectoryID != "" {
		parentDirectory, err := fs.directoryRepo.ReadByIDUserID(ctx, *request.DirectoryID, userID)
		if err != nil {
			return err
		}
		if parentDirectory == nil {
			return errors.New("parent directory not found")
		}
		directory.DirectoryID = &parentDirectory.ID
	}

	// Update the directory name
	directory.Name = request.Name

	// Save the updated directory in the repository
	updateData := bson.M{
		"name": directory.Name,
	}
	if directory.DirectoryID != nil {
		updateData["directory_id"] = directory.DirectoryID
	}
	err = fs.directoryRepo.Update(ctx, ID, updateData)
	if err != nil {
		return err
	}

	return nil
}
