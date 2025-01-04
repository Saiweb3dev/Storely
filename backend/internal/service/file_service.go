package service

import (
    "context"
    "fmt"
    
    "backend/internal/models"
    "backend/internal/repository"

)

// FileService is the service layer that handles business logic for file metadata.
type FileService struct {
    repo *repository.FileRepository // Dependency on the repository for database operations
}

// NewFileService creates a new instance of FileService.
// It accepts a FileRepository as a dependency.
func NewFileService(repo *repository.FileRepository) *FileService {
    return &FileService{repo: repo}
}

// CreateFileMetadata handles the creation of file metadata by calling the repository layer.
// It takes a context and a file metadata model as inputs.
func (s *FileService) CreateFileMetadata(ctx context.Context, file *models.FileMetadata) error {
    // Validate the input to ensure it is not nil
    if file == nil {
        return fmt.Errorf("file metadata cannot be nil")
    }

    // Call the repository's Create method to insert the file metadata
    err := s.repo.Create(ctx, file)
    if err != nil {
        // Wrap and return the error with additional context
        return fmt.Errorf("failed to create file metadata in service: %w", err)
    }

    return nil // Return nil if the operation succeeds
}
