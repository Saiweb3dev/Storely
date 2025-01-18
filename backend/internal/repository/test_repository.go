// repository/test_repository.go
package repository

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type TestRepository struct {
    collection *mongo.Collection
}

func NewTestRepository(client *mongo.Client) *TestRepository {
    return &TestRepository{
        collection: client.Database("Storely").Collection("testData"),
    }
}

func (repo *TestRepository) InsertTestData(ctx context.Context, data interface{}) error {
    _, err := repo.collection.InsertOne(ctx, bson.M{
        "payload":   data,
        "createdAt": time.Now(),
    })
    return err
}

func (repo *TestRepository) GetLastTestData(ctx context.Context) (bson.M, error) {
    // Sort by descending creation time, limit 1
    opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})
    var result bson.M
    err := repo.collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
    return result, err
}