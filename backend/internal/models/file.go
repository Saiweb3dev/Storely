package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetadata struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Filename    string            `bson:"filename" json:"filename"`
    Size        int64             `bson:"size" json:"size"`
    ContentType string            `bson:"content_type" json:"content_type"`
    UploadedAt  time.Time         `bson:"uploaded_at" json:"uploaded_at"`
}
