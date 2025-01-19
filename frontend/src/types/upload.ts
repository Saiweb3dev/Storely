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
  userId?: string;
  fileId: string;
  fileName: string;
  fileType: string;
  uploadedAt: Date;
  size?: number;
}

export interface UploadOptions {
  chunkSize?: number;
  maxRetries?: number;
}

export interface FileWithPreview extends File {
  preview: string;
}