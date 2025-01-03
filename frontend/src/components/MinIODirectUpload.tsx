// components/MinIODirectUpload.tsx
// Handles direct upload to MinIO. Reusable on any page or parent component.
"use client"
import axios from "axios"
import { useState } from "react"
import { UploadProgressInfo } from "@/types/minio"
import type { RecentUpload } from "@/types/upload"

interface MinIODirectUploadProps {
  file: File          // The file to upload
  onProgress: (info: UploadProgressInfo) => void
  onComplete: (upload: RecentUpload) => void
  onError: (message: string) => void
}

const CHUNK_SIZE = 5 * 1024 * 1024

export default function MinIODirectUpload({
  file,
  onProgress,
  onComplete,
  onError,
}: MinIODirectUploadProps) {
  const [uploading, setUploading] = useState(false)

  const uploadToMinIO = async () => {
    if (!file) return
    setUploading(true)
    try {
      // Initialize upload
      const initRes = await axios.post("http://localhost:8080/api/minio/files/init", {
        fileName: file.name,
        fileType: file.type,
        fileSize: file.size,
        totalChunks: Math.ceil(file.size / CHUNK_SIZE),
      })

      const { fileId, uploadUrls, callbackUrl } = initRes.data
      let completedChunks = 0

      // Upload all chunks
      for (let i = 0; i < uploadUrls.length; i++) {
        const { uploadUrl } = uploadUrls[i]
        const start = i * CHUNK_SIZE
        const end = Math.min(start + CHUNK_SIZE, file.size)
        const chunk = file.slice(start, end)
        await axios.put(uploadUrl, chunk, { headers: { "Content-Type": file.type } })
        completedChunks++
        onProgress({
          fileName: file.name,
          progress: (completedChunks / uploadUrls.length) * 100,
        })
      }

      // Complete upload
      await axios.post(callbackUrl)

      onComplete({
        fileId,
        fileName: file.name,
        fileType: file.type,
        uploadedAt: new Date(),
        size: file.size,
      })
    } catch (err: any) {
      onError(err.message || "MinIO upload failed")
    } finally {
      setUploading(false)
    }
  }

  return (
    <button
      onClick={uploadToMinIO}
      disabled={!file || uploading}
      className={` text-white px-4 py-2 rounded ${!file || uploading
        ? 'bg-gray-300 cursor-not-allowed'
        : 'bg-blue-500 hover:bg-blue-600 text-white'}
    `}
    >
      {uploading ? "Uploading..." : "Upload Files Direct"}
    </button>
  )
}