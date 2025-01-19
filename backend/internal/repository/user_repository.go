// internal/repository/user_repository.go
package repository

import (
    "context"
    "time"
    "fmt"
    "backend/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
    collection := client.Database("Storely").Collection("users")

    // Create unique indexes
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys:    bson.D{{Key: "name", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }
    
    _, err := collection.Indexes().CreateMany(context.Background(), indexes)
    if err != nil {
        fmt.Printf("failed to create indexes: %v\n", err)
        return nil
    }
    
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

func (r *UserRepository) CheckDuplicate(ctx context.Context, email, username string) (bool, string, error) {
    // Check email
    count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
    if err != nil {
        return false, "", err
    }
    if count > 0 {
        return true, "Email already registered", nil
    }

    // Check username
    count, err = r.collection.CountDocuments(ctx, bson.M{"name": username})
    if err != nil {
        return false, "", err
    }
    if count > 0 {
        return true, "Username already taken", nil
    }

    return false, "", nil
}

// In user_repository.go
func (r *UserRepository) GetStorageUsedAndLimit(ctx context.Context, userID string) (float64, float64, error) {
    filter := bson.M{"user_id": userID}
    var user models.User
    if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
        return 0, 0, fmt.Errorf("failed to find user: %w", err)
    }
    return user.StorageUsed, user.StorageLimit, nil
}

func (r *UserRepository) CheckUserStorageLimit(ctx context.Context, userID string, fileSize float64) error {
    used, limit, err := r.GetStorageUsedAndLimit(ctx, userID)
    if err != nil {
        return err
    }
    if used+fileSize > limit {
        return fmt.Errorf("storage limit exceeded")
    }
    _, err = r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$inc": bson.M{"storage_used": fileSize}})
    if err != nil {
        return fmt.Errorf("failed to update storage: %w", err)
    }
    return nil
}

func (r *UserRepository) DecreaseUsedStorage(ctx context.Context, userID string, size float64) error {
    // Optionally check current storage:
    current, _, err := r.GetStorageUsedAndLimit(ctx, userID)
    if err != nil {
        return err
    }

    // Clamp at 0
    decrement := -size
    if current < size {
        decrement = -current
    }

    _, err = r.collection.UpdateOne(
        ctx,
        bson.M{"user_id": userID},
        bson.M{"$inc": bson.M{"storage_used": decrement}},
    )
    if err != nil {
        return fmt.Errorf("failed to update storage: %w", err)
    }
    return nil
}