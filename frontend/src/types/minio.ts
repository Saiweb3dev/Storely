export interface MinIOUploadResponse {
  fileId: string;
  uploadUrls: {
    chunkIndex: number;
    uploadUrl: string;
  }[];
  callbackUrl: string;
}

export interface MinIOError {
  Error: {
    Code: string
    Message: string
    Key: string
    BucketName: string
  }
}

export interface RecentUpload {
  fileId: string;
  fileName: string;
  fileType: string;
  uploadedAt: Date;
  size?: number;
}

export interface MinIODirectUploadProps {
  file: File          // The file to upload
  onProgress: (info: UploadProgressInfo) => void
  onComplete: (upload: RecentUpload) => void
  onError: (message: string) => void
}


export interface UploadProgressInfo {
  fileName: string
  progress: number
}

export interface UploadResponse {
  fileId: string
  fileName: string
  fileType: string
}