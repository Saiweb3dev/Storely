// models/minio_file.go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Make sure we use one consistent type name: FileMinIO
type FileMinIO struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FileName    string            `bson:"file_name" json:"fileName"`
    FileType    string            `bson:"file_type" json:"fileType"`
    Size        int64             `bson:"size" json:"size"`
    TotalChunks int               `bson:"total_chunks" json:"totalChunks"`
    CreatedAt   time.Time         `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updatedAt"`
    Complete    bool              `bson:"complete" json:"complete"`
    MinioPath   string           `bson:"minio_path" json:"minioPath"`
    BucketName  string           `bson:"bucket_name" json:"bucketName"`
}