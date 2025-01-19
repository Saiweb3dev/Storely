// file_repository.go
package repository

import (
    "context"
    "fmt"
    "log"
    "time"

    "backend/internal/models"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// MinIOFileRepository handles MinIO-specific file operations
type MinIOFileRepository struct {
    collection *mongo.Collection
}

// NewMinIOFileRepository creates a new MinIO file repository
func NewMinIOFileRepository(client *mongo.Client) *MinIOFileRepository {
    return &MinIOFileRepository{
        collection: client.Database("Storely").Collection("minio_files"),
    }
}

// CreateFile_MinIO creates a new MinIO file record
func (r *MinIOFileRepository) CreateFile_MinIO(ctx context.Context, file *models.FileMinIO) error {
    if file == nil {
        return fmt.Errorf("file metadata cannot be nil")
    }

    log.Printf("Attempting to insert MinIO file metadata: %+v", file)

    _, err := r.collection.InsertOne(ctx, file)
    if err != nil {
        return fmt.Errorf("failed to insert MinIO file metadata: %w", err)
    }

    log.Printf("Successfully inserted MinIO file metadata: %+v", file)
    return nil
}

// GetFileByID_MinIO retrieves a MinIO file by ID
func (r *MinIOFileRepository) GetFileByID_MinIO(ctx context.Context, fileID string) (*models.FileMinIO, error) {
    objectID, err := primitive.ObjectIDFromHex(fileID)
    if err != nil {
        return nil, fmt.Errorf("invalid MinIO file ID format: %w", err)
    }

    var file models.FileMinIO
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("MinIO file not found")
        }
        return nil, fmt.Errorf("error retrieving MinIO file: %w", err)
    }

    return &file, nil
}

// MarkFileComplete_MinIO marks a MinIO file as complete
func (r *MinIOFileRepository) MarkFileComplete_MinIO(ctx context.Context, fileID string) error {
    objectID, err := primitive.ObjectIDFromHex(fileID)
    if err != nil {
        return fmt.Errorf("invalid MinIO file ID format: %w", err)
    }

    update := bson.M{
        "$set": bson.M{
            "complete": true,
            "updated_at": primitive.DateTime(time.Now().UnixNano() / 1e6),
        },
    }

    result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
    if err != nil {
        return fmt.Errorf("failed to mark MinIO file as complete: %w", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("MinIO file not found")
    }

    return nil
}

// UpdateMinIOPath updates the MinIO path for a file
func (r *MinIOFileRepository) UpdateMinIOPath(ctx context.Context, fileID string, minioPath string) error {
    objectID, err := primitive.ObjectIDFromHex(fileID)
    if err != nil {
        return fmt.Errorf("invalid MinIO file ID format: %w", err)
    }

    update := bson.M{
        "$set": bson.M{
            "minio_path": minioPath,
            "updated_at": primitive.DateTime(time.Now().UnixNano() / 1e6),
        },
    }

    result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
    if err != nil {
        return fmt.Errorf("failed to update MinIO path: %w", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("MinIO file not found")
    }

    return nil
}

// In repository/minio_file_repository.go
func (r *MinIOFileRepository) DeleteMinIOFile(ctx context.Context, fileID string) error {
    objectID, err := primitive.ObjectIDFromHex(fileID)
    if err != nil {
        return fmt.Errorf("invalid MinIO file ID format: %w", err)
    }
    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return fmt.Errorf("failed to delete MinIO file metadata: %w", err)
    }
    return nil
}