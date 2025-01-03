// internal/models/file.go

package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetadata struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FileName    string            `bson:"file_name" json:"fileName"`
    FileType    string            `bson:"file_type" json:"fileType"`
    Size        int64             `bson:"size" json:"size"`
    UploadedAt  time.Time         `bson:"uploaded_at" json:"uploadedAt"`
    Complete    bool              `bson:"complete" json:"complete"`
}

