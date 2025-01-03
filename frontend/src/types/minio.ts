interface MinIOUploadResponse {
  fileId: string;
  uploadUrls: {
    chunkIndex: number;
    uploadUrl: string;
  }[];
  callbackUrl: string;
}