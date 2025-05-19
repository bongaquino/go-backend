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

// New helper to fetch subdirectories and files for a directory
func (fs *FSService) fetchDirectoryContents(ctx context.Context, directory *model.Directory, userID string) ([]*model.Directory, []*model.File, error) {
	subDirectories, err := fs.directoryRepo.ListByDirectoryIDUserID(ctx, directory.ID.Hex(), userID)
	if err != nil {
		return nil, nil, err
	}

	// @TODO: Fetch files within the directory
	// files, err := fs.fileRepo.ListByDirectoryID(ctx, directory.ID.Hex())
	// if err != nil {
	// 	return nil, nil, err
	// }

	return subDirectories, nil, nil
}

func (fs *FSService) ReadRootDirectory(ctx context.Context, userID string) (*model.Directory,
	[]*model.Directory, []*model.File, error) {
	// Fetch the root directory from the repository
	directory, err := fs.directoryRepo.ReadByUserIDName(ctx, userID, "root")
	if err != nil {
		return nil, nil, nil, err
	}

	subDirectories, files, err := fs.fetchDirectoryContents(ctx, directory, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return directory, subDirectories, files, nil
}

func (fs *FSService) ReadDirectory(ctx context.Context, ID string, userID string) (*model.Directory,
	[]*model.Directory, []*model.File, error) {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	subDirectories, files, err := fs.fetchDirectoryContents(ctx, directory, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return directory, subDirectories, files, nil
}
