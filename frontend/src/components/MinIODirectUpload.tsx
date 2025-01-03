import React, { useState } from 'react';
import axios from 'axios';

interface MinIODirectUploadProps {
  onUploadComplete: (info: {
    fileId: string;
    name: string;
    type: string;
    size: number;
  }) => void;
  onError: (error: string) => void;
}
 const MinIODirectUpload: React.FC<MinIODirectUploadProps> = ({ onUploadComplete, onError }) => {
  const [isUploading, setIsUploading] = useState(false);
  const CHUNK_SIZE = 5 * 1024 * 1024; // 5MB chunks

  const handleDirectUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    try {
      // Initialize upload
      const response = await axios.post<MinIOUploadResponse>('http://localhost:8080/api/minio/files/init', {
        fileName: file.name,
        fileType: file.type,
        fileSize: file.size,
        totalChunks: Math.ceil(file.size / CHUNK_SIZE)
      });

      const { fileId, uploadUrls, callbackUrl } = response.data;

// if (callbackUrl) {
//   const fullCallbackUrl = new URL(callbackUrl, "http://localhost:8080").toString();
//   console.log(fullCallbackUrl);
// } else {
//   console.error("callbackUrl is undefined or missing in response.data");
// }



      // Split file into chunks and upload directly to MinIO
      const uploadPromises = uploadUrls.map(async ({ chunkIndex, uploadUrl }, index) => {
        const start = index * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);

        await axios.put(uploadUrl, chunk, {
          headers: { 'Content-Type': file.type }
        });
      });
      console.log("uploadPromises",uploadPromises);

      await Promise.all(uploadPromises);

      // Notify backend of completion
      console.log("Calling callbackUrl",callbackUrl);
      await axios.post(callbackUrl);
      onUploadComplete({
        fileId,
        name: file.name,
        type: file.type,
        size: file.size,
      });

    } catch (error) {
      onError(error instanceof Error ? error.message : 'Upload failed');
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="minio-direct-upload mt-4"> {/* Added margin-top */}
    <input
      type="file"
      onChange={handleDirectUpload}
      disabled={isUploading}
      style={{ display: 'none' }}
      id="minio-file-input"
    />
    <label 
      htmlFor="minio-file-input"
      className={`
        w-full px-4 py-2 rounded-md font-medium transition-colors duration-200 block text-center
        ${isUploading 
          ? 'bg-gray-300 cursor-not-allowed' 
          : 'bg-green-500 hover:bg-green-600 text-white cursor-pointer'
        }
      `}
    >
      {isUploading ? 'Uploading...' : 'Direct Upload to MinIO'}
    </label>
  </div>
  );
};


export default MinIODirectUpload;