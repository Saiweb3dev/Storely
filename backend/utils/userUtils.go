package utils

import (
    "crypto/sha256"
    "encoding/hex"
)

// GenerateUserID creates a unique ID based on user data
func GenerateUserID(name, email, password string) string {
    data := name + email + password
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}