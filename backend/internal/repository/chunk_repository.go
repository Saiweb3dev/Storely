package repository

import (
    "context"
		"fmt"
		"log"

    "backend/internal/models"
		
    "go.mongodb.org/mongo-driver/mongo"
	  "go.mongodb.org/mongo-driver/bson"
	  "go.mongodb.org/mongo-driver/bson/primitive"
)

type ChunkRepository struct {
	db *mongo.Database
}

func NewChunkRepository(client *mongo.Client) *ChunkRepository {
	if client == nil {
			log.Fatal("MongoDB client is nil")
	}
	return &ChunkRepository{
			db: client.Database("Storely"),
	}
}

func (r *ChunkRepository) SaveChunk(ctx context.Context, chunk *models.FileChunk) error {
	collection := r.db.Collection("chunks")
	
	result, err := collection.InsertOne(ctx, chunk)
	if err != nil {
			log.Printf("MongoDB InsertOne error: %v", err)
			return err
	}
	
	log.Printf("Inserted chunk with ID: %v", result.InsertedID)
	return nil
}

func (r *ChunkRepository) GetFileChunks(ctx context.Context, fileID string) ([]*models.FileChunk, error) {
    cursor, err := r.db.Collection("chunks").Find(ctx, bson.M{"file_id": fileID})
    if err != nil {
        return nil, err
    }
    
    var chunks []*models.FileChunk
    if err = cursor.All(ctx, &chunks); err != nil {
        return nil, err
    }
    return chunks, nil
}

func (r *ChunkRepository) CountChunks(ctx context.Context, fileID string) (int, error) {
	count, err := r.db.Collection("chunks").CountDocuments(ctx, bson.M{"file_id": fileID})
	return int(count), err
}

func (r *FileRepository) MarkFileComplete(ctx context.Context, fileID string) error {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
			return fmt.Errorf("invalid file ID format: %w", err)
	}

	// Update the file metadata to mark it as complete
	result, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			bson.M{"$set": bson.M{"complete": true}},
	)

	if err != nil {
			return fmt.Errorf("failed to mark file as complete: %w", err)
	}

	if result.MatchedCount == 0 {
			return fmt.Errorf("no file found with ID: %s", fileID)
	}

	log.Printf("Successfully marked file %s as complete", fileID)
	return nil
}