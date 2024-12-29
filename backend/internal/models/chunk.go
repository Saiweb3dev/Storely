package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type FileChunk struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    FileID      string            `bson:"file_id"`
    ChunkIndex  int               `bson:"chunk_index"`
    TotalChunks int              `bson:"total_chunks"`
    Data        []byte            `bson:"data"`
    FileName    string            `bson:"file_name"`
    FileType    string            `bson:"file_type"`
    UploadedAt  time.Time         `bson:"uploaded_at"`
}

type File struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    FileName  string            `bson:"file_name"`
    FileType  string            `bson:"file_type"`
    Size      int64             `bson:"size"`
    CreatedAt time.Time         `bson:"created_at"`
    Complete  bool              `bson:"complete"`
}