package config

import (
    "context"
    "log"
    "time"
    "os"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

// LoadEnv loads environment variables from a .env file or system environment.
func LoadEnv() {
    // Specify the relative path to the .env file
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    requiredVars := []string{"MONGODB_URI", "DB_NAME", "SERVER_PORT"}
    for _, v := range requiredVars {
        if os.Getenv(v) == "" {
            log.Fatalf("%s is not set in the environment variables", v)
        }
    }
}

// ConnectDB establishes a connection to MongoDB using URI from environment variables.
func ConnectDB() (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Get MongoDB URI from environment variables
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        log.Fatal("MONGODB_URI is not set in the environment variables")
    }

    log.Printf("Connecting to Database: %s", uri)
    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // Ping the database to verify connection
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    return client, nil
}

//Connection to MinIO client
func ConnectMinIO() (*minio.Client,error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_SECRET_KEY")

    if endpoint == "" || accessKey == "" || secretKey == "" {
        log.Fatal("MINIO_ENDPOINT, MINIO_ACCESS_KEY, and MINIO_SECRET_KEY are not set in the environment variables")
    }

    client, err := minio.New(endpoint, &minio.Options{
        Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: true,
    })

    if err != nil {
        return nil, err
    }

    log.Printf("Connected to MinIO at %s", endpoint)
    return client, nil
}