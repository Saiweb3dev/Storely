// handlers/user_handler.go
package handlers

import (
    "backend/internal/models"
    "backend/internal/service"
    "encoding/json"
    "net/http"
    "log"
    "time"
    "fmt"
    "backend/utils/crypto"
    "backend/middleware"
    "backend/utils"

    // "go.mongodb.org/mongo-driver/bson"
    // "go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
    userService *service.UserService
}

type LoginCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Parse incoming encrypted request
    var encryptedData struct {
        Data string `json:"data"`
    }
    if err := json.NewDecoder(r.Body).Decode(&encryptedData); err != nil {
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    // Decrypt the request data
    decryptedData, err := crypto.Decrypt(encryptedData.Data)
    if err != nil {
        http.Error(w, "Failed to decrypt data", http.StatusBadRequest)
        return
    }

    // Parse user registration data
    var userData struct {
        UserName string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.Unmarshal(decryptedData, &userData); err != nil {
        http.Error(w, "Invalid data format", http.StatusBadRequest)
        return
    }

    // Create new user model
    user := &models.User{
        Name:         userData.UserName,
        Email:        userData.Email,
        Password:     userData.Password,
        IPAddress:    middleware.GetIP(r),
        StorageUsed:  float64(utils.MBToBytes(0)),
        StorageLimit: float64(utils.MBToBytes(10240)),
        CreatedAt:    time.Now(),
    }
    user.UserID = utils.GenerateUserID(user.Name, user.Email, user.Password)

    // Register user in database
    if err := h.userService.RegisterUser(*user); err != nil {
        log.Printf("Registration failed: %v", err)
        http.Error(w, "Registration failed", http.StatusInternalServerError)
        return
    }

    // Generate JWT token for new user
    token, err := h.userService.GenerateToken(user)
    if err != nil {
        log.Printf("Token generation failed: %v", err)
        http.Error(w, "Registration failed", http.StatusInternalServerError)
        return
    }

    // Prepare response data
    responseData := map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "username":     user.Name,
            "email":       user.Email,
            "storageUsed": user.StorageUsed,
            "storageLimit": user.StorageLimit,
        },
    }

    // Encrypt response
    jsonData, err := json.Marshal(responseData)
    if err != nil {
        http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
        return
    }

    encryptedResponse, err := crypto.Encrypt(jsonData)
    if err != nil {
        http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
        return
    }

    // Send encrypted response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "data": encryptedResponse,
    })
    
    log.Printf("User registered successfully: %s", userData.Email)
}


func (h *UserHandler) parseLoginRequest(r *http.Request) (*LoginCredentials, error) {
    // Change to match registration format
    var encryptedData struct {
        Data string `json:"data"`
    }

    if err := json.NewDecoder(r.Body).Decode(&encryptedData); err != nil {
        return nil, fmt.Errorf("failed to read request body: %w", err)
    }

    // Use the Data field instead of FormValue
    decryptedData, err := crypto.Decrypt(encryptedData.Data)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }

    var creds LoginCredentials
    if err := json.Unmarshal(decryptedData, &creds); err != nil {
        return nil, fmt.Errorf("failed to parse credentials: %w", err)
    }

    return &creds, nil
}

func (h *UserHandler) sendEncryptedResponse(w http.ResponseWriter, user *models.User, token string) error {
    responseData := map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "username": user.Name,
            "email": user.Email,
            "storageUsed": 0, // Need to make this dynamic
            "storageLimit": 10, // Default storage limit
            "createdAt": user.CreatedAt,
        },
        "message": "Login successful",
    }

    jsonData, err := json.Marshal(responseData)
    if err != nil {
        return fmt.Errorf("failed to marshal response: %w", err)
    }

    encryptedResponse, err := crypto.Encrypt(jsonData)
    if err != nil {
        return fmt.Errorf("failed to encrypt response: %w", err)
    }

    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(map[string]string{
        "data": encryptedResponse,
    })
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "OPTIONS" {
        h.handleCORS(w)
        return
    }

    log.Printf("Received login request from IP: %s", middleware.GetIP(r))

    creds, err := h.parseLoginRequest(r)
    if err != nil {
        log.Printf("Login request parsing failed: %v", err)
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    // Handle login attempt using service
    user, err := h.userService.HandleLoginAttempt(
        r.Context(),
        creds.Email,
        creds.Password,
        middleware.GetIP(r),
    )

    // Handle different types of errors
    if err != nil {
        switch err {
        case service.ErrAccountLocked:
            log.Printf("Login attempted on locked account: %s", creds.Email)
            http.Error(w, "Account is locked", http.StatusForbidden)
        case service.ErrInvalidCredentials:
            log.Printf("Invalid credentials for: %s", creds.Email)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        default:
            log.Printf("Login error for %s: %v", creds.Email, err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Generate JWT token
    token, err := h.userService.GenerateToken(user)
    if err != nil {
        log.Printf("Token generation failed for %s: %v", creds.Email, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Send encrypted response
    if err := h.sendEncryptedResponse(w, user, token); err != nil {
        log.Printf("Failed to send response for %s: %v", creds.Email, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Printf("Login successful for user: %s", creds.Email)
}

func (h *UserHandler) handleCORS(w http.ResponseWriter) {
    h.setCORSHeaders(w)
    w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) setCORSHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Content-Encrypted")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}