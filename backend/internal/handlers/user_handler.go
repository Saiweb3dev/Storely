// handlers/user_handler.go
package handlers

import (
    "encoding/json"
    "io"
    "net/http"
    "backend/internal/models"
    "backend/internal/service"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    if r.Method == "OPTIONS" {
        h.handleCORS(w)
        return
    }

    h.setCORSHeaders(w)
    w.Header().Set("Content-Type", "application/json")

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := json.Unmarshal(body, &user); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    if err := h.userService.RegisterUser(user); err != nil {
        http.Error(w, "Unable to register user", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "OPTIONS" {
        h.handleCORS(w)
        return
    }

    h.setCORSHeaders(w)
    w.Header().Set("Content-Type", "application/json")

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.Unmarshal(body, &creds); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    user, err := h.userService.AuthenticateUser(creds.Email, creds.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := h.userService.GenerateToken(user)
    if err != nil {
        http.Error(w, "Token generation failed", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) handleCORS(w http.ResponseWriter) {
    h.setCORSHeaders(w)
    w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) setCORSHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}