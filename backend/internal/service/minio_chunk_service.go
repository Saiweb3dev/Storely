// minio_chunk_service.go
package service

import (
    "bytes"
    "context"
    "fmt"
    "io"

    "backend/internal/models"
    
    "github.com/minio/minio-go/v7"
)

type MinIOChunkService struct {
    client  *minio.Client
    bucket  string
}

func NewMinIOChunkService(client *minio.Client, bucket string) *MinIOChunkService {
    return &MinIOChunkService{
        client: client,
        bucket: bucket,
    }
}

func (s *MinIOChunkService) SaveChunk(ctx context.Context, chunk *models.FileChunk) error {
    objectName := fmt.Sprintf("%s/chunk_%d", chunk.FileID, chunk.ChunkIndex)
    
    // Convert chunk data to io.Reader
    reader := bytes.NewReader(chunk.Data)
    
    _, err := s.client.PutObject(ctx, s.bucket, objectName, reader, int64(len(chunk.Data)), minio.PutObjectOptions{
        ContentType: "application/octet-stream",
    })
    
    return err
}

func (s *MinIOChunkService) GetFileFromMinIO(ctx context.Context, fileID string) ([]byte, error) {
    var allData []byte
    
    // List all chunks for this file
    opts := minio.ListObjectsOptions{
        Prefix: fileID + "/",
        Recursive: true,
    }
    
    for object := range s.client.ListObjects(ctx, s.bucket, opts) {
        obj, err := s.client.GetObject(ctx, s.bucket, object.Key, minio.GetObjectOptions{})
        if err != nil {
            return nil, err
        }
        chunkData, err := io.ReadAll(obj)
        if err != nil {
            return nil, err
        }
        
        allData = append(allData, chunkData...)
    }
    
    return allData, nil
}