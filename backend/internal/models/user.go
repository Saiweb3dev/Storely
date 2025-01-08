// internal/models/user.go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    UserID            string            `json:"user_id" bson:"user_id"`
    Name              string            `json:"name" bson:"name"`
    Email             string            `json:"email" bson:"email"`
    Password          string            `json:"password" bson:"password"`
    CreatedAt         time.Time         `json:"created_at" bson:"created_at"`
    LastLoginAt       time.Time         `json:"last_login_at" bson:"last_login_at"`
    LoginCount        int               `json:"login_count" bson:"login_count"`
    FailedAttempts    int              `json:"failed_attempts" bson:"failed_attempts"`
    LastAttemptAt     time.Time         `json:"last_attempt_at" bson:"last_attempt_at"`
    IsLocked          bool              `json:"is_locked" bson:"is_locked"`
    LockExpiresAt     *time.Time        `json:"lock_expires_at,omitempty" bson:"lock_expires_at,omitempty"`
    IPAddress         string            `json:"ip_address" bson:"ip_address"`
}