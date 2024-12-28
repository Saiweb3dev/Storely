// src/services/uploadService.ts
import axios from "axios";

const DEFAULT_CHUNK_SIZE = 1024 * 1024; // 1MB
const MAX_CHUNK_SIZE = 5 * 1024 * 1024; // 5MB

export interface UploadProgressInfo {
  fileName: string;
  progress: number;
}

export interface UploadOptions {
  chunkSize?: number;
  maxRetries?: number;
}

export type ProgressCallback = (info: UploadProgressInfo) => void;

class UploadService {
  private calculateChunkSize(fileSize: number, requestedSize?: number): number {
    if (!requestedSize) {
      // For small files, use smaller chunks
      if (fileSize < 5 * 1024 * 1024) {  // < 5MB
        return Math.min(DEFAULT_CHUNK_SIZE, fileSize);
      }
      // For larger files, scale chunk size with file size
      return Math.min(Math.ceil(fileSize / 100), MAX_CHUNK_SIZE);
    }
    
    return Math.min(requestedSize, MAX_CHUNK_SIZE);
  }

  private async uploadChunk(
    chunk: Blob,
    fileName: string,
    chunkIndex: number,
    totalChunks: number
  ) {
    const formData = new FormData();
    formData.append("file", chunk, fileName);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());

    return axios.post("http://localhost:8080/upload-chunk", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  }

  async uploadFile(
    file: File, 
    onProgress: ProgressCallback, 
    options: UploadOptions = {}
  ) {
    const chunkSize = this.calculateChunkSize(file.size, options.chunkSize);
    const chunks = Math.ceil(file.size / chunkSize);

    for (let i = 0; i < chunks; i++) {
      const chunk = file.slice(
        i * chunkSize,
        Math.min((i + 1) * chunkSize, file.size)
      );

      await this.uploadChunk(chunk, file.name, i, chunks);

      onProgress({
        fileName: file.name,
        progress: ((i + 1) / chunks) * 100
      });
    }
  }

  async uploadFiles(
    files: File[], 
    onProgress: ProgressCallback,
    options?: UploadOptions
  ) {
    for (const file of files) {
      await this.uploadFile(file, onProgress, options);
    }
  }
}

export const uploadService = new UploadService();