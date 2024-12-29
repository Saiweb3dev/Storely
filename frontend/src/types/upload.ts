// types/upload.ts
export interface UploadProgressInfo {
  fileName: string;
  progress: number;
}

export interface UploadResponse {
  fileId: string;
  fileName: string;
  fileType: string;
  totalChunks: number;
  chunksReceived: number;
}

export interface RecentUpload {
  fileId: string;
  fileName: string;
  fileType: string;
  uploadedAt: Date;
}

export interface UploadOptions {
  chunkSize?: number;
  maxRetries?: number;
}