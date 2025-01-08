// handlers/user_handler.go
package handlers

import (
    "backend/internal/models"
    "backend/internal/service"
    "encoding/json"
    "net/http"
    "log"
    "backend/utils/crypto"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var encryptedData struct {
        Data string `json:"data"`
    }

    if err := json.NewDecoder(r.Body).Decode(&encryptedData); err != nil {
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    decryptedData, err := crypto.Decrypt(encryptedData.Data)
    if err != nil {
        http.Error(w, "Failed to decrypt data", http.StatusBadRequest)
        return
    }

    var userData struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.Unmarshal(decryptedData, &userData); err != nil {
        http.Error(w, "Invalid data format", http.StatusBadRequest)
        return
    }


    log.Printf("Registering user: %s (%s)", userData.Name, userData.Email)

    user := &models.User{
        Name:     userData.Name,
        Email:    userData.Email,
        Password: userData.Password,
    }

    if err := h.userService.RegisterUser(*user); err != nil {
        log.Printf("Registration failed: %v", err)
        http.Error(w, "Registration failed", http.StatusInternalServerError)
        return
    }

    log.Printf("User registered successfully: %s", userData.Email)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User registered successfully",
    })
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "OPTIONS" {
        h.handleCORS(w)
        return
    }
    h.setCORSHeaders(w)
    log.Printf("Received login request")

    var encryptedData struct {
        Data string `json:"data"`
    }

    if err := json.NewDecoder(r.Body).Decode(&encryptedData); err != nil {
        log.Printf("Failed to decode login request: %v", err)
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    log.Printf("Encrypted login data received")

    decryptedData, err := crypto.Decrypt(encryptedData.Data)
    if err != nil {
        log.Printf("Login decryption failed: %v", err)
        http.Error(w, "Failed to decrypt data", http.StatusBadRequest)
        return
    }

    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.Unmarshal(decryptedData, &creds); err != nil {
        log.Printf("Failed to parse login credentials: %v", err)
        http.Error(w, "Invalid decrypted data format", http.StatusBadRequest)
        return
    }

    
    log.Printf("Attempting login for user: %s", creds.Email)

    user, err := h.userService.AuthenticateUser(creds.Email, creds.Password)
    if err != nil {
        log.Printf("Authentication failed for %s: %v", creds.Email, err)
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := h.userService.GenerateToken(user)
    if err != nil {
        log.Printf("Token generation failed for %s: %v", creds.Email, err)
        http.Error(w, "Token generation failed", http.StatusInternalServerError)
        return
    }

    // Create response data
    responseData := map[string]string{
        "token": token,
        "message": "Login successful",
    }
    jsonData, err := json.Marshal(responseData)
    if err != nil {
        log.Printf("Failed to marshal response for %s: %v", creds.Email, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Encrypt response
    encryptedResponse, err := crypto.Encrypt(jsonData)
    if err != nil {
        log.Printf("Failed to encrypt response for %s: %v", creds.Email, err)
        http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
        return
    }

    log.Printf("Login successful for user: %s", creds.Email)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "data": encryptedResponse,
    })
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