package api

import (
    "net/http"
    "backend/internal/service"
    "backend/internal/handlers"
    "github.com/gorilla/mux"
)

// NewRouter sets up all API routes.
func NewRouter(fileService *service.FileService) http.Handler {
    router := mux.NewRouter()

    // Initialize handlers
    fileHandler := handlers.NewFileHandler(fileService)

    // Register routes
    router.HandleFunc("/file_metadata", fileHandler.CreateFileMetadata).Methods("POST")

    return router
}
