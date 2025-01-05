// internal/repository/user_repository.go
package repository

import (
    "context"
    "backend/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
    collection := client.Database("Storely").Collection("users")
    return &UserRepository{collection: collection}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}