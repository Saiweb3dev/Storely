// services/uploadService.ts
import axios from "axios";
import { UploadResponse, UploadProgressInfo, UploadOptions } from "../types/upload";

const DEFAULT_CHUNK_SIZE = 1024 * 1024; // 1MB
const MAX_CHUNK_SIZE = 5 * 1024 * 1024; // 5MB

export type ProgressCallback = (info: UploadProgressInfo) => void;

class UploadService {
  private calculateChunkSize(fileSize: number, requestedSize?: number): number {
    if (!requestedSize) {
      if (fileSize < 5 * 1024 * 1024) {
        return Math.min(DEFAULT_CHUNK_SIZE, fileSize);
      }
      return Math.min(Math.ceil(fileSize / 100), MAX_CHUNK_SIZE);
    }
    return Math.min(requestedSize, MAX_CHUNK_SIZE);
  }

  private async uploadChunk(
    chunk: Blob,
    fileName: string,
    chunkIndex: number,
    totalChunks: number,
    fileId?: string
  ): Promise<UploadResponse> {
    const formData = new FormData();
    formData.append("file", chunk, fileName);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());
    if (fileId) {
      formData.append("fileId", fileId);
    }

    const response = await axios.post<UploadResponse>(
      "http://localhost:8080/upload-chunk",
      formData
    );
    return response.data;
  }

  async uploadFile(file: File, onProgress: ProgressCallback): Promise<UploadResponse> {
    const chunkSize = this.calculateChunkSize(file.size);
    const chunks = Math.ceil(file.size / chunkSize);
    let fileId: string | undefined;
    let lastResponse: UploadResponse;

    for (let i = 0; i < chunks; i++) {
      const chunk = file.slice(
        i * chunkSize,
        Math.min((i + 1) * chunkSize, file.size)
      );

      lastResponse = await this.uploadChunk(chunk, file.name, i, chunks, fileId);
      fileId = lastResponse.fileId;

      onProgress({
        fileName: file.name,
        progress: ((i + 1) / chunks) * 100,
      });
    }

    return lastResponse!;
  }

  async uploadFiles(
    files: File[],
    onProgress: ProgressCallback,
    options?: UploadOptions
  ): Promise<UploadResponse[]> {
    return Promise.all(
      files.map((file) => this.uploadFile(file, onProgress))
    );
  }
}

export const uploadService = new UploadService();