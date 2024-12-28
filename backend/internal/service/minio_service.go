package service

import (
    "context"
    "fmt"
    "io"
    "net/url"
    "path"

		"github.com/minio/minio-go/v7"
)

type MinIOService struct {
    client *minio.Client
    bucket string
    endpoint string
}

func NewMinIOService(client *minio.Client, bucket string, endpoint string) *MinIOService {
    return &MinIOService{
        client: client,
        bucket: bucket,
        endpoint: endpoint,
    }
}

func (s *MinIOService) UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64) (string, error) {
    _, err := s.client.PutObject(ctx, s.bucket, fileName, file, fileSize, minio.PutObjectOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
    }

    // Parse the endpoint URL to ensure we use the correct protocol
    endpointURL, err := url.Parse(s.endpoint)
    if err != nil {
        return "", fmt.Errorf("invalid endpoint URL: %w", err)
    }

    // Construct the file URL using the same protocol as the endpoint
    fileURL := &url.URL{
        Scheme: endpointURL.Scheme,
        Host:   endpointURL.Host,
        Path:   path.Join(s.bucket, fileName),
    }

    return fileURL.String(), nil
}