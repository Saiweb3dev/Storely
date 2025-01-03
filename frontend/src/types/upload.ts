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
  size?: number; // Add size here
}

export interface UploadOptions {
  chunkSize?: number;
  maxRetries?: number;
}