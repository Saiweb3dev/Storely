// internal/repository/user_repository.go
package repository

import (
    "context"
    "time"
    "backend/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *UserRepository) UpdateLoginStats(ctx context.Context, userID primitive.ObjectID, ipAddress string, isSuccessful bool) error {
    now := time.Now()
    updates := bson.M{
        "$set": bson.M{
            "last_attempt_at": now,
            "ip_address":     ipAddress,
        },
    }

    if isSuccessful {
        updates["$set"].(bson.M)["last_login_at"] = now
        updates["$set"].(bson.M)["failed_attempts"] = 0
        updates["$set"].(bson.M)["is_locked"] = false
        updates["$set"].(bson.M)["lock_expires_at"] = nil
        updates["$inc"] = bson.M{"login_count": 1}
    } else {
        updates["$inc"] = bson.M{"failed_attempts": 1}
    }

    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, updates)
    return err
}

func (r *UserRepository) UpdateLockStatus(ctx context.Context, userID primitive.ObjectID, lockExpiry *time.Time) error {
    update := bson.M{
        "$set": bson.M{
            "is_locked":      lockExpiry != nil,
            "lock_expires_at": lockExpiry,
        },
    }
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
    return err
}