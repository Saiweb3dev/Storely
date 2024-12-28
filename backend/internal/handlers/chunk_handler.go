// internal/handlers/chunk_handler.go
package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    // "github.com/gorilla/mux"
)

type ChunkHandler struct {}

func NewChunkHandler() *ChunkHandler {
    return &ChunkHandler{}
}

func (h *ChunkHandler) HandleChunkUpload(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10MB max memory
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get file chunk
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Get chunk metadata
    chunkIndex := r.FormValue("chunkIndex")
    totalChunks := r.FormValue("totalChunks")

    // Convert to integers
    chunkIndexInt, err := strconv.Atoi(chunkIndex)
    if err != nil {
        http.Error(w, "Invalid chunk index", http.StatusBadRequest)
        return
    }

    totalChunksInt, err := strconv.Atoi(totalChunks)
    if err != nil {
        http.Error(w, "Invalid total chunks", http.StatusBadRequest)
        return
    }

    // Log chunk information (for testing)
    fmt.Printf("Received chunk %d of %d for file %s (size: %d bytes)\n",
        chunkIndexInt+1, totalChunksInt, header.Filename, header.Size)

    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf(`{
        "status": "success",
        "message": "Chunk received",
        "chunk": %d,
        "total": %d,
        "filename": "%s"
    }`, chunkIndexInt+1, totalChunksInt, header.Filename)))
}