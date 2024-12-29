"use client"
import { useState } from 'react'
import { motion } from 'framer-motion'
import { Download, AlertCircle } from 'lucide-react'

const FileDownload = () => {
  const [fileId, setFileId] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleDownload = async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`http://localhost:8080/files/${fileId}`, {
        headers: {
          'Accept': '*/*' // Accept any content type
        }
      })
      if (!response.ok) throw new Error('Download failed')

      // Get content type and filename from headers
      const contentType = response.headers.get('Content-Type') || 'application/octet-stream'
      const disposition = response.headers.get('Content-Disposition')
      
      // Extract filename with extension from Content-Disposition
      let filename = 'downloaded-file'
      if (disposition) {
        const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(disposition)
        if (matches != null && matches[1]) {
          filename = matches[1].replace(/['"]/g, '')
        }
      }
  
      const blob = await response.blob()
      const blobWithType = new Blob([blob], { type: contentType })
      const url = window.URL.createObjectURL(blobWithType)
  
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename) // Use original filename with extension
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (err) {
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
        
        <motion.button
          onClick={handleDownload}
          disabled={loading || !fileId.trim()}
          className={`px-4 py-2 rounded-md font-medium transition-colors duration-200
            ${loading || !fileId.trim()
              ? 'bg-gray-300 cursor-not-allowed'
              : 'bg-blue-500 hover:bg-blue-600 text-white'
            }`}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          {loading ? 'Downloading...' : 'Download'}
        </motion.button>
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
    </motion.div>
  )
}

export default FileDownload

