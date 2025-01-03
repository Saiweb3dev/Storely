// router.go
package api

import (
    "backend/internal/handlers"
    "backend/internal/service"
    "backend/internal/repository"
    "github.com/gorilla/mux"
    "github.com/minio/minio-go/v7"
    "go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(
    mongoClient *mongo.Client,
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

    // Initialize MinIO repository and handler
    minioRepo := repository.NewMinIOFileRepository(mongoClient)
    minioFileHandler := handlers.NewMinIOFileHandler(minioRepo, minioClient, bucket)
    
    // MinIO specific routes
    router.HandleFunc("/api/minio/files/init", minioFileHandler.InitializeMinIOUpload).Methods("POST")
    router.HandleFunc("/api/minio/files/{fileId}/complete", minioFileHandler.CompleteMinIOUpload).Methods("POST")

    // Chunk Upload Endpoints
    chunkHandler := handlers.NewChunkHandler(chunkRepo,fileRepo, minioClient, bucket)
    router.HandleFunc("/upload-chunk", chunkHandler.HandleChunkUpload).Methods("POST", "OPTIONS")
    router.HandleFunc("/files/{fileId}", chunkHandler.GetCompleteFile).Methods("GET")
    router.HandleFunc("/files/minio/{fileId}", chunkHandler.GetFileFromMinIO).Methods("GET")


    return router
}