// handlers/minio_file_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils"
    "backend/utils/logger"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.uber.org/zap"
)

type MinIOFileHandler struct {
    minioRepo   *repository.MinIOFileRepository
    userRepo   *repository.UserRepository
    minioClient *minio.Client
    bucketName  string
}

func NewMinIOFileHandler(minioRepo *repository.MinIOFileRepository,userRepo *repository.UserRepository, minioClient *minio.Client, bucketName string) *MinIOFileHandler {
    return &MinIOFileHandler{
        minioRepo:   minioRepo,
        userRepo:   userRepo,
        minioClient: minioClient,
        bucketName:  bucketName,
    }
}

func (h *MinIOFileHandler) InitializeMinIOUpload(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID      string `json:"userID"`
        FileName    string `json:"fileName"`
        FileType    string `json:"fileType"`
        FileSize    float64  `json:"fileSize"`
        TotalChunks int    `json:"totalChunks"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.userRepo.CheckUserStorageLimit(r.Context(), req.UserID, req.FileSize); err != nil {
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
        logger.L().Error("File Creation failed",
        zap.String("userID",file.UserID),
                zap.String("File Name",file.FileName),
                zap.String("File Type",file.FileType),
                zap.Float64("File Size",file.Size),
            zap.Error(err))
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

    logger.L().Info("File Upload Initialized",
     zap.String("File ID",file.ID.String()),
     zap.String("userID",file.UserID),
     zap.String("File Name",file.FileName),
     zap.String("File Type",file.FileType),
     zap.Float64("File Size",file.Size),
    )

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *MinIOFileHandler) CompleteMinIOUpload(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileID := vars["fileId"]

    file, err := h.minioRepo.GetFileByID_MinIO(r.Context(), fileID)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        logger.L().Error("File Not found in MinIO",
        zap.String("userID",file.UserID),
                zap.String("File Name",file.FileName),
                zap.String("File Type",file.FileType),
                zap.Float64("File Size",file.Size),
            zap.Error(err))
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
        logger.L().Error("Failed to mark file complete",
        zap.String("userID",file.UserID),
                zap.String("File Name",file.FileName),
                zap.String("File Type",file.FileType),
                zap.Float64("File Size",file.Size),
            zap.Error(err))
        return
    }

    logger.L().Info("File Upload Completed",
     zap.String("File ID",file.ID.String()),
     zap.String("userID",file.UserID),
     zap.String("File Name",file.FileName),
     zap.String("File Type",file.FileType),
     zap.Float64("File Size",file.Size),
    )

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "fileId": fileID,
    })
}

func (h *MinIOFileHandler) GetUserStorageHealth(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to get user storage health")
    
    // Validate JWT token
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        http.Error(w, "Authorization header missing", http.StatusUnauthorized)
        return
    }
    tokenString = strings.TrimPrefix(tokenString, "Bearer ")
    token, err := utils.ValidateJWT(tokenString)
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Get userID from query parameters instead of body
    userID := r.URL.Query().Get("userID")
    if userID == "" {
        http.Error(w, "Missing userID parameter", http.StatusBadRequest)
        return
    }

    used, limit, err := h.userRepo.GetStorageUsedAndLimit(r.Context(), userID)
    if err != nil {
        http.Error(w, "Unable to retrieve storage data", http.StatusInternalServerError)
        return
    }
    balance := limit - used

    logger.L().Info("User Calling for Storage Health Data",
     zap.String("userID",userID),
     zap.String("Storage Balance",fmt.Sprintf("%f",balance)),
    )

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "storageUsed":      used,
        "storageLimit":     limit,
        "availableBalance": balance,
    })
}

// In handlers/minio_file_handler.go
func (h *MinIOFileHandler) DeleteFileFromMinIO(w http.ResponseWriter, r *http.Request) {

    log.Println("Received request to delete file from MinIO")

    // Validate JWT token
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        http.Error(w, "Authorization header missing", http.StatusUnauthorized)
        return
    }
    tokenString = strings.TrimPrefix(tokenString, "Bearer ")
    token, err := utils.ValidateJWT(tokenString)
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Parse request body for userId and fileId
    var req struct {
        FileID string `json:"fileId"`
        UserID string `json:"userId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    log.Println("FileID:",req.FileID,"UserID:",req.UserID)

    // Retrieve file from DB
    file, err := h.minioRepo.GetFileByID_MinIO(r.Context(), req.FileID)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Check user mismatch
    if file.UserID != req.UserID {
        http.Error(w, "Not authorized to delete this file", http.StatusForbidden)
        return
    }

    // Remove all chunks
    for i := 0; i < file.TotalChunks; i++ {
        objectName := fmt.Sprintf("%s/chunk_%d", req.FileID, i)
        removeErr := h.minioClient.RemoveObject(r.Context(), h.bucketName, objectName, minio.RemoveObjectOptions{})
        if removeErr != nil {
            http.Error(w, "Failed removing chunk(s)", http.StatusInternalServerError)
            return
        }
    }

    // Remove metadata
    err = h.minioRepo.DeleteMinIOFile(r.Context(), req.FileID)
    if err != nil {
        http.Error(w, "Failed to remove metadata", http.StatusInternalServerError)
        return
    }

    if err := h.userRepo.DecreaseUsedStorage(r.Context(), file.UserID,file.Size); err != nil {
        log.Println("Failed to decrement user storage:", err)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "deleted",
        "fileId": req.FileID,
    })
}