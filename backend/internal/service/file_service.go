package service

import (
    "context"
    "backend/internal/models"
    "backend/internal/repository"
)

type FileService struct {
    repo *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository) *FileService {
    return &FileService{repo: repo}
}

func (s *FileService) CreateFileMetadata(ctx context.Context, file *models.FileMetadata) error {
    return s.repo.Create(ctx, file)
}