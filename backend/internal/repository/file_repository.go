package repository

import (
    "context"
    "backend/internal/models"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
)

// FileRepository is a struct that holds the MongoDB collection for file metadata.
type FileRepository struct {
    collection *mongo.Collection // MongoDB collection for storing file metadata
}

// NewFileRepository initializes a new instance of FileRepository.
// It sets up a connection to the "file_metadata" collection in the "Storely" database.
func NewFileRepository(client *mongo.Client) *FileRepository {
    collection := client.Database("Storely").Collection("file_metadata")
    return &FileRepository{collection: collection}
}

// Create inserts a new file metadata record into the database.
// It takes a context and a file metadata model as inputs.
func (r *FileRepository) Create(ctx context.Context, file *models.FileMetadata) error {
    // Validate that the file metadata is not nil
    if file == nil {
        return fmt.Errorf("file metadata cannot be nil")
    }

    // Log the operation
    log.Printf("Attempting to insert file metadata: %+v", file)

    // Perform the insertion
    _, err := r.collection.InsertOne(ctx, file)
    if err != nil {
        // Wrap and return the error with context for debugging
        return fmt.Errorf("failed to insert file metadata: %w", err)
    }

    // Log the success of the operation
    log.Printf("Successfully inserted file metadata: %+v", file)
    return nil
}
