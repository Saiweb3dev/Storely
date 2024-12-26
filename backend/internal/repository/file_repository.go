package repository

import (
    "context"
    "backend/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
    collection *mongo.Collection
}

func NewFileRepository(client *mongo.Client) *FileRepository {
    collection := client.Database("Storely").Collection("file_metadata")
    return &FileRepository{collection: collection}
}

func (r *FileRepository) Create(ctx context.Context, file *models.FileMetadata) error {
    _, err := r.collection.InsertOne(ctx, file)
    return err
}