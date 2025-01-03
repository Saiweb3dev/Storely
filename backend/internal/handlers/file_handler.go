package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "backend/internal/models"
    "backend/internal/service"
    "log"
)

// FileHandler is the HTTP handler for managing file metadata operations.
type FileHandler struct {
    service *service.FileService // Dependency on the service layer for business logic
}



// NewFileHandler creates a new instance of FileHandler.
// It accepts a FileService as a dependency.
func NewFileHandler(service *service.FileService) *FileHandler {
    return &FileHandler{service: service}
}

// CreateFileMetadata handles HTTP POST requests to create new file metadata.
func (h *FileHandler) CreateFileMetadata(w http.ResponseWriter, r *http.Request) {
    var file models.FileMetadata

    // Decode the JSON request body into the file metadata model
    if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest) // Return a 400 Bad Request response
        log.Printf("Failed to decode request body: %v", err)        // Log the error
        return
    }

    // Set the current time as the uploaded timestamp
    file.UploadedAt = time.Now()

    // Call the service layer to create the file metadata
    err := h.service.CreateFileMetadata(r.Context(), &file)
    if err != nil {
        http.Error(w, "Failed to create file metadata", http.StatusInternalServerError) // Return a 500 Internal Server Error response
        log.Printf("Service error: %v", err)                                           // Log the error
        return
    }

    // Respond with the created file metadata and a 201 Created status
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(file); err != nil {
        log.Printf("Failed to encode response: %v", err) // Log any error in encoding response
    }
}
