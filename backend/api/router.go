package api

import (
    "backend/internal/handlers"
    "backend/internal/service"
    "github.com/gorilla/mux"
    "github.com/minio/minio-go/v7"
)

func NewRouter(fileService *service.FileService, minioClient *minio.Client, bucket string) *mux.Router {
    router := mux.NewRouter()

    // File Metadata Endpoint
    fileHandler := handlers.NewFileHandler(fileService)
    router.HandleFunc("/file_metadata", fileHandler.CreateFileMetadata).Methods("POST")

    // MinIO File Upload Endpoint
		endpoint := "http://localhost:9000"
    minioService := service.NewMinIOService(minioClient, "storely-test", endpoint)
    minioHandler := handlers.NewMinIOHandler(minioService)
    router.HandleFunc("/upload", minioHandler.UploadFile).Methods("POST")

    chunkHandler := handlers.NewChunkHandler()
    router.HandleFunc("/upload-chunk", chunkHandler.HandleChunkUpload).Methods("POST", "OPTIONS")

    return router
}
