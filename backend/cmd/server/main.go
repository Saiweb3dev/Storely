package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "backend/config"
    "backend/internal/handlers"
    "backend/internal/repository"
    "backend/internal/service"
    "github.com/gorilla/mux"
)

func main() {
    // Connect to MongoDB
    client, err := config.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Ensure disconnection when the program exits
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := client.Disconnect(ctx); err != nil {
            log.Printf("Error disconnecting from MongoDB: %v", err)
        }
    }()

    fileRepo := repository.NewFileRepository(client)
    fileService := service.NewFileService(fileRepo)
    fileHandler := handlers.NewFileHandler(fileService)

    router := mux.NewRouter()
    router.HandleFunc("/file_metadata", fileHandler.CreateFileMetadata).Methods("POST")

    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Channel to listen for errors coming from the server
    serverErrors := make(chan error, 1)

    // Start server
    go func() {
        fmt.Printf("🚀 Server starting on http://localhost%s\n", server.Addr)
        fmt.Println("Press Ctrl+C to stop the server")
        serverErrors <- server.ListenAndServe()
    }()

    // Channel to listen for interrupt signal
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    // Block until either an error or interrupt occurs
    select {
    case err := <-serverErrors:
        log.Fatalf("Server error: %v", err)
    case <-shutdown:
        fmt.Println("\nServer is shutting down...")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // Gracefully shutdown the server
        if err := server.Shutdown(ctx); err != nil {
            log.Printf("Could not gracefully shutdown the server: %v\n", err)
            return
        }
        fmt.Println("Server stopped gracefully")
    }
}