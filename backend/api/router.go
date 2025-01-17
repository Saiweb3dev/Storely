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
	userRepo *repository.UserRepository,
	userService *service.UserService,
	bucket string,
) *mux.Router {
	router := mux.NewRouter()

	minioRepo := repository.NewMinIOFileRepository(mongoClient)
	minioFileHandler := handlers.NewMinIOFileHandler(minioRepo,userRepo, minioClient, bucket)
	chunkHandler := handlers.NewChunkHandler(chunkRepo, fileRepo, minioRepo, minioClient, bucket)
	userHandler := handlers.NewUserHandler(userService)

	router.HandleFunc("/upload-chunk", chunkHandler.HandleChunkUpload).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/minio/files/init", minioFileHandler.InitializeMinIOUpload).Methods("POST")
	router.HandleFunc("/files/minio/{fileId}", chunkHandler.GetFileFromMinIO).Methods("GET")
	
	router.HandleFunc("/api/auth/register", userHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", userHandler.Login).Methods("POST", "OPTIONS")
	
	router.HandleFunc("/api/minio/files/{fileId}/complete", minioFileHandler.CompleteMinIOUpload).Methods("POST")

	router.HandleFunc("/get/user/storageHealth", minioFileHandler.GetUserStorageHealth).Methods("GET")

	return router
}
