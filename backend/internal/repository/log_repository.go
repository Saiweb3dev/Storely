package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "time"
)

// LogRepository handles log storage in MongoDB.
type LogRepository struct {
    logCollection      *mongo.Collection
    errorLogCollection *mongo.Collection
}

// NewLogRepository creates a new instance of LogRepository.
// It accepts a MongoDB client and sets up collections for logs and error logs.
func NewLogRepository(client *mongo.Client) *LogRepository {
    return &LogRepository{
        logCollection:      client.Database("Storely").Collection("logs"),
        errorLogCollection: client.Database("Storely").Collection("error_logs"),
    }
}

// InsertLog inserts a normal log into the logs collection.
func (repo *LogRepository) InsertLog(ctx context.Context, logEntry map[string]interface{}) error {
    logEntry["created_at"] = time.Now() // Add a timestamp
    _, err := repo.logCollection.InsertOne(ctx, logEntry)
    return err
}

// InsertErrorLog inserts an error log into the error_logs collection.
func (repo *LogRepository) InsertErrorLog(ctx context.Context, errorLogEntry map[string]interface{}) error {
    errorLogEntry["created_at"] = time.Now() // Add a timestamp
    _, err := repo.errorLogCollection.InsertOne(ctx, errorLogEntry)
    return err
}
