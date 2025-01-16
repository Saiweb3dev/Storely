// handlers/minio_file_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "fmt"

    "backend/internal/models"
    "backend/internal/repository"
    
    "github.com/gorilla/mux"
    "github.com/minio/minio-go/v7"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type MinIOFileHandler struct {
    minioRepo   *repository.MinIOFileRepository
    minioClient *minio.Client
    bucketName  string
}

func NewMinIOFileHandler(minioRepo *repository.MinIOFileRepository, minioClient *minio.Client, bucketName string) *MinIOFileHandler {
    return &MinIOFileHandler{
        minioRepo:   minioRepo,
        minioClient: minioClient,
        bucketName:  bucketName,
    }
}

func (h *MinIOFileHandler) InitializeMinIOUpload(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID      string `json:"userID"`
        FileName    string `json:"fileName"`
        FileType    string `json:"fileType"`
        FileSize    int64  `json:"fileSize"`
        TotalChunks int    `json:"totalChunks"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    file := &models.FileMinIO{
        ID:          primitive.NewObjectID(),
        UserID:      req.UserID,
        FileName:    req.FileName,
        FileType:    req.FileType,
        Size:        req.FileSize,
        TotalChunks: req.TotalChunks,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Complete:    false,
        BucketName:  h.bucketName,
    }

    if err := h.minioRepo.CreateFile_MinIO(r.Context(), file); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var uploadURLs []map[string]interface{}
    for i := 0; i < req.TotalChunks; i++ {
        objectName := fmt.Sprintf("%s/chunk_%d", file.ID.Hex(), i)
        url, err := h.minioClient.PresignedPutObject(r.Context(), h.bucketName, objectName, time.Hour)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        uploadURLs = append(uploadURLs, map[string]interface{}{
            "chunkIndex": i,
            "uploadUrl":  url.String(),
        })
    }

    response := map[string]interface{}{
        "fileId":      file.ID.Hex(),
        "uploadUrls":  uploadURLs,
        "callbackUrl": fmt.Sprintf("http://localhost:8080/api/minio/files/%s/complete", file.ID.Hex()),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *MinIOFileHandler) CompleteMinIOUpload(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileID := vars["fileId"]

    file, err := h.minioRepo.GetFileByID_MinIO(r.Context(), fileID)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Verify all chunks exist
    for i := 0; i < file.TotalChunks; i++ {
        objectName := fmt.Sprintf("%s/chunk_%d", fileID, i)
        _, err := h.minioClient.StatObject(r.Context(), h.bucketName, objectName, minio.StatObjectOptions{})
        if err != nil {
            http.Error(w, "Missing chunks", http.StatusBadRequest)
            return
        }
    }

    if err := h.minioRepo.MarkFileComplete_MinIO(r.Context(), fileID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "fileId": fileID,
    })
}