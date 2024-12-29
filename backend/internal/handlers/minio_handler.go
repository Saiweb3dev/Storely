package handlers

import (
    "net/http"
    "backend/internal/service"
	"log"
)

type MinIOHandler struct {
    service *service.MinIOService
}

func NewMinIOHandler(service *service.MinIOService) *MinIOHandler {
    return &MinIOHandler{service: service}
}

// UploadFile handles file upload to MinIO.
func (h *MinIOHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file upload handler")
	log.Printf("Request Headers: %v", r.Header)
	log.Printf("Request Content Length: %d", r.ContentLength)

    file, header, err := r.FormFile("file")
    if err != nil {
			log.Printf("Error reading file: %v", err)
        http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    fileSize := header.Size
    fileName := header.Filename
		log.Printf("File received: %s, Size: %d bytes", fileName, fileSize)

    fileURL, err := h.service.UploadFile(r.Context(), fileName, file, fileSize)
    if err != nil {
			log.Printf("Error uploading file to MinIO: %v", err)
        http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("File uploaded successfully. URL: %s", fileURL)
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("File uploaded successfully: " + fileURL))
}
