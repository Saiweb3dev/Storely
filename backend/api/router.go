// router.go
package api

import (
    "backend/internal/handlers"
    "backend/internal/service"
    "backend/internal/repository"
    "github.com/gorilla/mux"
    "github.com/minio/minio-go/v7"
)

func NewRouter(
    fileService *service.FileService,
    minioClient *minio.Client,
    chunkRepo *repository.ChunkRepository,
    fileRepo *repository.FileRepository,
    bucket string,
) *mux.Router {
    router := mux.NewRouter()

    // File Metadata Endpoint
    fileHandler := handlers.NewFileHandler(fileService)
    router.HandleFunc("/file_metadata", fileHandler.CreateFileMetadata).Methods("POST")

    // MinIO File Upload Endpoint (kept for future use)
    endpoint := "http://localhost:9000"
    minioService := service.NewMinIOService(minioClient, bucket, endpoint)
    minioHandler := handlers.NewMinIOHandler(minioService)
    router.HandleFunc("/upload", minioHandler.UploadFile).Methods("POST")

    // Chunk Upload Endpoints
    chunkHandler := handlers.NewChunkHandler(chunkRepo,fileRepo)
    router.HandleFunc("/upload-chunk", chunkHandler.HandleChunkUpload).Methods("POST", "OPTIONS")
    router.HandleFunc("/files/{fileId}", chunkHandler.GetCompleteFile).Methods("GET")

    return router
}