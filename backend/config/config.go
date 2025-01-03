package config

import (
    "context"
    "log"
    "time"
    "os"
    "fmt"

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
// In config/config.go
func ConnectMinIO() (*minio.Client, error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_SECRET_KEY")
    bucketName := os.Getenv("MINIO_BUCKET_NAME")

    if endpoint == "" || accessKey == "" || secretKey == "" {
        log.Fatal("MinIO environment variables not set")
    }

    // Initialize MinIO client with HTTP
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: false,  // Set to false for HTTP
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create MinIO client: %w", err)
    }

    // Create bucket if it doesn't exist
    exists, err := client.BucketExists(context.Background(), bucketName)
    if err != nil {
        return nil, fmt.Errorf("failed to check bucket existence: %w", err)
    }

    if !exists {
        err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to create bucket: %w", err)
        }
        log.Printf("Created new bucket: %s", bucketName)
        
        // Set bucket policy
        if err := setPublicBucketPolicy(client, bucketName); err != nil {
            log.Printf("Warning: Failed to set bucket policy: %v", err)
        }
    }

    log.Printf("Successfully connected to MinIO at %s", endpoint)
    return client, nil
}

// Add this function to config.go
func setPublicBucketPolicy(client *minio.Client, bucketName string) error {
    policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {"AWS": ["*"]},
                "Action": ["s3:GetBucketLocation", "s3:ListBucket", "s3:ListBucketMultipartUploads"],
                "Resource": ["arn:aws:s3:::` + bucketName + `"]
            },
            {
                "Effect": "Allow",
                "Principal": {"AWS": ["*"]},
                "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
                "Resource": ["arn:aws:s3:::` + bucketName + `/*"]
            }
        ]
    }`

    return client.SetBucketPolicy(context.Background(), bucketName, policy)
}