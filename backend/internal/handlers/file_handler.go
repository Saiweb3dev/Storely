package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "backend/internal/models"
    "backend/internal/service"
)

type FileHandler struct {
    service *service.FileService
}

func NewFileHandler(service *service.FileService) *FileHandler {
    return &FileHandler{service: service}
}

func (h *FileHandler) CreateFileMetadata(w http.ResponseWriter, r *http.Request) {
    var file models.FileMetadata
    if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    file.UploadedAt = time.Now()
    
    err := h.service.CreateFileMetadata(r.Context(), &file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(file)
}