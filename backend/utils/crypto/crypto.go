package crypto

import (
    "encoding/base64"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

var secretKey string

func InitCrypto(key string) error {
    if key == "" {
        return fmt.Errorf("encryption key cannot be empty")
    }
    secretKey = key
    return nil
}

func Decrypt(data string) ([]byte, error) {
    // Split data and hash
    parts := strings.Split(data, ".")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid data format")
    }
    
    base64Data, hash := parts[0], parts[1]
    
    // Verify hash
    h := sha256.New()
    h.Write([]byte(base64Data + secretKey))
    calculatedHash := hex.EncodeToString(h.Sum(nil))
    
    if hash != calculatedHash {
        return nil, fmt.Errorf("invalid hash")
    }
    
    // Decode base64
    decodedData, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return nil, fmt.Errorf("failed to decode data: %v", err)
    }
    
    return decodedData, nil
}

func Encrypt(data []byte) (string, error) {
    // Encode data to base64
    base64Data := base64.StdEncoding.EncodeToString(data)
    
    // Create hash
    h := sha256.New()
    h.Write([]byte(base64Data + secretKey))
    hash := hex.EncodeToString(h.Sum(nil))
    
    // Combine data and hash
    return fmt.Sprintf("%s.%s", base64Data, hash), nil
}