// internal/handlers/chunk_handler.go
package handlers

import (
    "fmt"
    "io"
    "encoding/json"
    "encoding/base64"
	"log"
    "net/http"
    "strconv"
    "time"
    "context"
    "bytes"

    "backend/internal/models"
    "backend/internal/repository"
    
    "github.com/gorilla/mux"
    "github.com/minio/minio-go/v7"
    "go.mongodb.org/mongo-driver/bson/primitive"

)
type ChunkHandler struct {
    chunkRepo    *repository.ChunkRepository
    fileRepo     *repository.FileRepository
    minioRepo    *repository.MinIOFileRepository
    minioClient  *minio.Client
    bucketName   string
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details"`
}
type ChunkUploadResponse struct {
    FileID      string `json:"fileId"`
    FileName    string `json:"fileName"`
    FileType    string `json:"fileType"`
    TotalChunks int    `json:"totalChunks"`
    ChunksReceived int `json:"chunksReceived"`
}

func NewChunkHandler(
    chunkRepo *repository.ChunkRepository,
    fileRepo *repository.FileRepository,
    minioRepo *repository.MinIOFileRepository, 
    minioClient *minio.Client,
    bucketName string,
) *ChunkHandler {
    return &ChunkHandler{
        chunkRepo:   chunkRepo,
        fileRepo:    fileRepo,
        minioRepo:   minioRepo,                  
        minioClient: minioClient,
        bucketName:  bucketName,
    }
}

func (h *ChunkHandler) HandleChunkUpload(w http.ResponseWriter, r *http.Request) {
    var fileID string
    isFirstChunk := false
    
    // Get or create fileID
    fileID = r.FormValue("fileId")
    if fileID == "" {
        fileID = primitive.NewObjectID().Hex()
        isFirstChunk = true
    }

    chunkIndex, _ := strconv.Atoi(r.FormValue("chunkIndex"))
    totalChunks, _ := strconv.Atoi(r.FormValue("totalChunks"))
    
    file, header, err := r.FormFile("file")
    if err != nil {
        log.Printf("Error getting file: %v", err)
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        log.Printf("Error reading file: %v", err)
        http.Error(w, "Failed to read file", http.StatusInternalServerError)
        return
    }

    // If first chunk, create file metadata
    if isFirstChunk {
        fileMetadata := &models.File{
            ID:        func() primitive.ObjectID {
                id, err := primitive.ObjectIDFromHex(fileID)
                if err != nil {
                    log.Printf("Error converting fileID to ObjectID: %v", err)
                    http.Error(w, "Invalid file ID", http.StatusBadRequest)
                    return primitive.NilObjectID
                }
                return id
            }(),
            FileName:  header.Filename,
            FileType:  header.Header.Get("Content-Type"),
            Size:      r.ContentLength * int64(totalChunks),
            CreatedAt: time.Now(),
            Complete:  false,
        }
        
        if err := h.fileRepo.CreateFile(r.Context(), fileMetadata); err != nil {
            log.Printf("Error creating file metadata: %v", err)
            http.Error(w, "Failed to create file metadata", http.StatusInternalServerError)
            return
        }
    }

    chunk := &models.FileChunk{
        FileID:      fileID,
        ChunkIndex:  chunkIndex,
        TotalChunks: totalChunks,
        Data:        data,
        FileName:    header.Filename,
        FileType:    header.Header.Get("Content-Type"),
        UploadedAt:  time.Now(),
    }

    // Save to MongoDB
    if err := h.chunkRepo.SaveChunk(r.Context(), chunk); err != nil {
        log.Printf("Error saving chunk: %v", err)
        http.Error(w, "Failed to save chunk", http.StatusInternalServerError)
        return
    }

      // Additional MinIO storage
    chunkData, _ := base64.StdEncoding.DecodeString(string(chunk.Data))
    objectName := fmt.Sprintf("%s/chunk_%d", chunk.FileID, chunk.ChunkIndex)
    
    _, err = h.minioClient.PutObject(
        context.Background(),
        h.bucketName,
        objectName,
        bytes.NewReader(chunkData),
        int64(len(chunkData)),
        minio.PutObjectOptions{},
    )
    if err != nil {
        log.Printf("MinIO storage error: %v", err)
        errorResponse := ErrorResponse{
            Error:   "Failed to store chunk in MinIO",
            Details: err.Error(),
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    // Count received chunks
    chunksReceived, err := h.chunkRepo.CountChunks(r.Context(), fileID)
    if err != nil {
        log.Printf("Error counting chunks: %v", err)
    }

    // Update file completion status if all chunks received
    if chunksReceived == totalChunks {
        if err := h.fileRepo.MarkFileComplete(r.Context(), fileID); err != nil {
            log.Printf("Error marking file complete: %v", err)
        } else {
            log.Printf("Successfully marked file %s as complete", fileID)
        }
    }

    response := ChunkUploadResponse{
        FileID:         fileID,
        FileName:       header.Filename,
        FileType:       header.Header.Get("Content-Type"),
        TotalChunks:    totalChunks,
        ChunksReceived: chunksReceived,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}


func (h *ChunkHandler) GetFileFromMinIO(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileID := vars["fileId"]

    fileMetadata, err := h.minioRepo.GetFileByID_MinIO(r.Context(), fileID)
    if err != nil {
        log.Printf("Error getting file metadata: %v", err)
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    var downloadUrls []string
    // Generate presigned URLs for each chunk
    for i := 0; i < fileMetadata.TotalChunks; i++ {
        objectName := fmt.Sprintf("%s/chunk_%d", fileID, i)
        presignedURL, err := h.minioClient.PresignedGetObject(
            r.Context(),
            h.bucketName,
            objectName,
            time.Hour,
            nil,
        )
        if err != nil {
            log.Printf("Error generating presigned URL for chunk %d: %v", i, err)
            http.Error(w, "Failed to generate download URLs", http.StatusInternalServerError)
            return
        }
        downloadUrls = append(downloadUrls, presignedURL.String())
    }

    response := map[string]interface{}{
        "downloadUrls": downloadUrls,
        "fileName":    fileMetadata.FileName,
        "fileType":    fileMetadata.FileType,
        "totalChunks": fileMetadata.TotalChunks,
        "expiresIn":   "1 hour",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}