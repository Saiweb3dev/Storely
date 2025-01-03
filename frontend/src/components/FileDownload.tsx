// FileDownload.tsx
"use client"
import { useState } from 'react'
import { motion } from 'framer-motion'
import { Download, AlertCircle } from 'lucide-react'
import { XMLParser } from 'fast-xml-parser'

interface MinIOError {
  Error: {
    Code: string
    Message: string
    Key: string
    BucketName: string
  }
}

const FileDownload = () => {
  const [fileId, setFileId] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [progress, setProgress] = useState<number>(0)

  const parseMinIOError = async (response: Response): Promise<string> => {
    const text = await response.text()
    try {
      const parser = new XMLParser()
      const result = parser.parse(text) as MinIOError
      return `MinIO Error: ${result.Error.Message} (Code: ${result.Error.Code})`
    } catch {
      return text
    }
  }
// FileDownload.tsx
const handleDownload = async () => {
  setLoading(true)
  setError(null)
  setProgress(0)

  try {
    const response = await fetch(`http://localhost:8080/files/minio/${fileId}`)
    
    if (!response.ok) {
      const errorMessage = await parseMinIOError(response)
      throw new Error(errorMessage)
    }

    const data = await response.json()
    console.log('File metadata received:', data)

    if (!data.downloadUrls || !data.downloadUrls.length) {
      throw new Error('No download URLs provided')
    }

    // Create array to store chunks in correct order
    const orderedChunks: Blob[] = new Array(data.downloadUrls.length)
    
    // Download chunks with retries
    const downloadChunk = async (url: string, index: number, retries = 3): Promise<void> => {
      try {
        const chunkResponse = await fetch(url)
        if (!chunkResponse.ok) throw new Error(`Failed to download chunk ${index}`)
        
        const blob = await chunkResponse.blob()
        orderedChunks[index] = blob
        
        setProgress((orderedChunks.filter(Boolean).length / data.downloadUrls.length) * 100)
      } catch (error) {
        if (retries > 0) {
          console.log(`Retrying chunk ${index}, ${retries} attempts remaining`)
          await downloadChunk(url, index, retries - 1)
        } else {
          throw error
        }
      }
    }

    // Process chunks in batches to avoid memory issues
    const BATCH_SIZE = 10
    for (let i = 0; i < data.downloadUrls.length; i += BATCH_SIZE) {
      const batch = data.downloadUrls
        .slice(i, i + BATCH_SIZE)
        .map((url: string, batchIndex: number) => 
          downloadChunk(url, i + batchIndex)
        )
      
      await Promise.all(batch)
    }

    // Verify all chunks downloaded successfully
    if (orderedChunks.some(chunk => !chunk)) {
      throw new Error('Some chunks failed to download')
    }

    // Combine chunks in correct order
    const completeFile = new Blob(orderedChunks, { type: data.fileType })
    
    // Create download link
    const url = window.URL.createObjectURL(completeFile)
    const link = document.createElement('a')
    link.style.display = 'none'
    link.href = url
    link.download = data.fileName
    
    document.body.appendChild(link)
    link.click()
    
    // Cleanup
    window.URL.revokeObjectURL(url)
    link.remove()
    orderedChunks.length = 0 // Clear array to free memory
    
    console.log('Download completed:', data.fileName)

  } catch (err) {
    console.error('Download failed:', err)
    setError(err instanceof Error ? err.message : 'Download failed')
  } finally {
    setLoading(false)
  }
}
  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg"
    >
      <motion.h2 
        className="text-2xl text-gray-800 font-bold mb-6 flex items-center justify-center"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2, duration: 0.5 }}
      >
        <Download className="mr-2 h-6 w-6 text-blue-500" />
        Download File
      </motion.h2>
      
      <motion.div 
        className="flex gap-2"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.4, duration: 0.5 }}
      >
        <input
          type="text"
          value={fileId}
          onChange={(e) => setFileId(e.target.value)}
          placeholder="Enter file ID"
          className="flex-1 text-gray-800 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          disabled={loading}
        />
        <button
          onClick={handleDownload}
          disabled={loading || !fileId}
          className={`px-4 py-2 rounded-md ${
            loading || !fileId 
              ? 'bg-gray-300 cursor-not-allowed'
              : 'bg-blue-500 hover:bg-blue-600 text-white'
          }`}
        >
          {loading ? 'Downloading...' : 'Download'}
        </button>
      </motion.div>

      {error && (
        <motion.div 
          className="mt-4 p-4 bg-red-100 text-red-700 rounded-md flex items-center"
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <AlertCircle className="mr-2 h-5 w-5" />
          {error}
        </motion.div>
      )}

      {loading && progress > 0 && (
        <motion.div className="mt-4">
          <div className="w-full bg-gray-200 rounded-full h-2">
            <motion.div
              className="bg-blue-500 h-2 rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
            />
          </div>
          <p className="text-sm text-gray-600 text-center mt-2">
            {Math.round(progress)}%
          </p>
        </motion.div>
      )}
    </motion.div>
  )
}

export default FileDownload