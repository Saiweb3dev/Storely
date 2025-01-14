// internal/models/user.go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    UserID       string            `bson:"user_id"`
    Name         string            `bson:"name"`
    Email        string            `bson:"email"`
    Password     string            `bson:"password"`
    StorageUsed  float64          `bson:"storage_used"`
    StorageLimit float64          `bson:"storage_limit"`
    CreatedAt    time.Time        `bson:"created_at"`
    IPAddress    string           `bson:"ip_address"`
    IsLocked     bool             `bson:"is_locked"`
    LockExpiresAt *time.Time      `bson:"lock_expires_at,omitempty"`
    FailedAttempts int            `bson:"failed_attempts"`
}