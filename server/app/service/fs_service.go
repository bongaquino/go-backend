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
func NewFSService(directoryRepo *repository.DirectoryRepository,
	fileRepo *repository.FileRepository,
) *FSService {
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

	// Fetch files within the root directory
	files, err := fs.fileRepo.ListByDirectoryIDUserID(ctx, directory.ID.Hex(), userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Return the directory details
	return directory, subDirectories, files, nil
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

	// Fetch files within the specified directory
	files, err := fs.fileRepo.ListByDirectoryIDUserID(ctx, directory.ID.Hex(), userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Return the directory details
	return directory, subDirectories, files, nil
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

	// Check if the directory is the root directory
	if directory.Name == "root" {
		return errors.New("cannot update root directory")
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

func (fs *FSService) DeleteDirectory(ctx context.Context, ID string, userID string) error {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Check if the directory exists
	if directory == nil {
		return errors.New("directory not found")
	}

	// Check if the directory is not the root directory
	if directory.Name == "root" {
		return errors.New("cannot delete root directory")
	}

	// Initialize a queue for BFS traversal
	queue := []string{ID}

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		// Mark the current directory as deleted
		err = fs.directoryRepo.Update(ctx, currentID, bson.M{"is_deleted": true})
		if err != nil {
			return err
		}

		// Mark all files in the current directory as deleted
		files, err := fs.fileRepo.ListByDirectoryIDUserID(ctx, currentID, userID)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = fs.fileRepo.Update(ctx, file.ID.Hex(), bson.M{"is_deleted": true})
			if err != nil {
				return err
			}
		}

		// Fetch all subdirectories of the current directory
		subdirs, err := fs.directoryRepo.ListByDirectoryIDUserID(ctx, currentID, userID)
		if err != nil {
			return err
		}

		// Enqueue all subdirectory IDs
		for _, subdir := range subdirs {
			queue = append(queue, subdir.ID.Hex())
		}
	}

	return nil
}

func (fs *FSService) CheckDirectoryOwnership(ctx context.Context, ID string, userID string) (bool, error) {
	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, ID, userID)

	if err != nil {
		return false, err
	}

	// Check if the directory exists
	if directory == nil {
		return false, errors.New("directory not found")
	}

	// Return true if the user owns the directory
	return true, nil
}

func (fs *FSService) CreateFile(ctx context.Context, file *model.File) error {
	// Create the file in the repository
	err := fs.fileRepo.Create(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FSService) ReadFileByIDUserID(ctx context.Context, ID string, userID string) (*model.File, error) {
	// Fetch the file from the repository
	file, err := fs.fileRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return nil, err
	}

	// Check if the file exists
	if file == nil {
		return nil, errors.New("file not found")
	}

	return file, nil
}

func (fs *FSService) UpdateFile(ctx context.Context, ID string, userID string, request *dto.UpdateFileDTO) error {
	// Fetch the file from the repository
	file, err := fs.fileRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Check if the file exists
	if file == nil {
		return errors.New("file not found")
	}

	// Fetch the directory from the repository
	directory, err := fs.directoryRepo.ReadByIDUserID(ctx, *request.DirectoryID, userID)
	if err != nil {
		return err
	}

	// Check if the directory exists
	if directory == nil {
		return errors.New("directory not found")
	}

	// Save the updated file in the repository
	updateData := bson.M{
		"name": request.Name,
	}
	if request.DirectoryID != nil {
		updateData["directory_id"] = request.DirectoryID
	}
	err = fs.fileRepo.Update(ctx, ID, updateData)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FSService) DeleteFile(ctx context.Context, ID string, userID string) error {
	// Fetch the file from the repository
	file, err := fs.fileRepo.ReadByIDUserID(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Check if the file exists
	if file == nil {
		return errors.New("file not found")
	}
	// Save the updated file in the repository
	updateData := bson.M{
		"is_deleted": true,
	}
	err = fs.fileRepo.Update(ctx, ID, updateData)
	if err != nil {
		return err
	}

	return nil
}
