// router.go
package api

import (
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/service"

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


	// Initialize MinIO repository and handler
	minioRepo := repository.NewMinIOFileRepository(mongoClient)
	minioFileHandler := handlers.NewMinIOFileHandler(minioRepo, minioClient, bucket)
    chunkHandler := handlers.NewChunkHandler(chunkRepo, fileRepo, minioRepo, minioClient, bucket)

	//used externally
	router.HandleFunc("/upload-chunk", chunkHandler.HandleChunkUpload).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/minio/files/init", minioFileHandler.InitializeMinIOUpload).Methods("POST")
	router.HandleFunc("/files/minio/{fileId}", chunkHandler.GetFileFromMinIO).Methods("GET")

    //used internally
    router.HandleFunc("/api/minio/files/{fileId}/complete", minioFileHandler.CompleteMinIOUpload).Methods("POST")

	return router
}
