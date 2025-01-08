package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "backend/config"
    "backend/internal/repository"
    "backend/internal/service"
    "backend/api"
    "backend/middleware"
    "backend/utils/crypto"
    "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    // Load environment variables
    config.LoadEnv()

    encryptionKey := os.Getenv("ENCRYPTION_KEY")
    if err := crypto.InitCrypto(encryptionKey); err != nil {
        log.Fatalf("Failed to initialize crypto: %v", err)
    }

    // Initialize MongoDB client
    client, err := config.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    
    // Initialize MinIO client (kept for future use)
    minioClient, err := config.ConnectMinIO()
    if err != nil {
        log.Fatalf("Failed to connect to MinIO: %v", err)
    }
    
    defer func() {
        if err := disconnectMongo(client); err != nil {
            log.Printf("Error disconnecting MongoDB: %v", err)
        }
    }()

    fileRepo := repository.NewFileRepository(client)
    chunkRepo := repository.NewChunkRepository(client)
    userRepo := repository.NewUserRepository(client) // Add this
    
    // Initialize services
    fileService := service.NewFileService(fileRepo)
    userService := service.NewUserService(userRepo) // Add this
    

    // Create router and register API routes
    bucket := os.Getenv("MINIO_BUCKET_NAME")
    router := api.NewRouter(client,fileService, minioClient, chunkRepo,fileRepo,userService, bucket)

    // Configure CORS middleware
    corsMiddleware := middleware.CORS(router)

    // Start the HTTP server
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    server := startServer(corsMiddleware, port)
    gracefulShutdown(server)
}

func disconnectMongo(client *mongo.Client) error {
    ctx, cancel := createTimeoutContext(10 * time.Second)
    defer cancel()
    return client.Disconnect(ctx)
}

// startServer initializes and starts the HTTP server.
func startServer(handler http.Handler, port string) *http.Server {
    server := &http.Server{
        Addr:    ":" + port,
        Handler: handler,
    }

    go func() {
        log.Printf("🚀 Server starting on http://localhost:%s\n", port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    return server
}

// gracefulShutdown handles server shutdown gracefully.
func gracefulShutdown(server *http.Server) {
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    <-shutdown // Wait for shutdown signal
    log.Println("🔄 Initiating graceful server shutdown...")

    ctx, cancel := createTimeoutContext(10 * time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Error shutting down server: %v", err)
    }

    log.Println("✅ Server stopped gracefully")
}

// createTimeoutContext creates a context with a timeout.
func createTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), timeout)
}
