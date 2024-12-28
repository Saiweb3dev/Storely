// FileUpload.tsx
"use client"
import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import { uploadService, UploadProgressInfo } from "../service/uploadService";

interface FileWithPreview extends File {
  preview: string;
}

export default function FileUpload() {
  const [files, setFiles] = useState<FileWithPreview[]>([])
  const [uploadProgress, setUploadProgress] = useState<UploadProgressInfo[]>([])
  const [isUploading, setIsUploading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const onDrop = useCallback((acceptedFiles: File[]) => {
    setFiles(acceptedFiles.map(file => Object.assign(file, {
      preview: URL.createObjectURL(file)
    })))
    setError(null)
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({ 
    onDrop,
    accept: {
      'image/*': [],
      'application/pdf': ['.pdf'],
      'application/msword': ['.doc'],
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document': ['.docx'],
    }
  })

  const removeFile = (file: FileWithPreview) => {
    const newFiles = [...files]
    newFiles.splice(newFiles.indexOf(file), 1)
    setFiles(newFiles)
    setError(null)
  }

  const handleUpload = async () => {
    if (files.length === 0) {
      setError('Please select files to upload')
      return
    }

    setIsUploading(true)
    setError(null)

    try {
      const updateProgress = (info: UploadProgressInfo) => {
        setUploadProgress(prev => 
          prev.map(p => p.fileName === info.fileName ? info : p)
        )
      }

      setUploadProgress(files.map(f => ({ fileName: f.name, progress: 0 })))
      
      await uploadService.uploadFiles(files, updateProgress, {
        chunkSize: 2 * 1024 * 1024 // 2MB chunks
      })

      setFiles([])
      setUploadProgress([])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setIsUploading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
      <div {...getRootProps()} className={`
        border-2 border-dashed rounded-lg p-6 cursor-pointer
        ${isDragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}
      `}>
        <input {...getInputProps()} />
        {isDragActive ? (
          <p className="text-center text-blue-500">Drop the files here...</p>
        ) : (
          <p className="text-center text-gray-500">
            Drag & drop files here, or click to select files
          </p>
        )}
      </div>

      {error && (
        <div className="mt-4 p-3 bg-red-100 text-red-700 rounded">
          {error}
        </div>
      )}

      {files.length > 0 && (
        <div className="mt-4">
          <h3 className="font-medium">Selected Files:</h3>
          {files.map(file => (
            <div key={file.name} className="mt-2 flex items-center justify-between">
              <span className="truncate flex-1">{file.name}</span>
              <button
                onClick={() => removeFile(file)}
                className="ml-2 text-red-500 hover:text-red-700"
                disabled={isUploading}
              >
                Remove
              </button>
            </div>
          ))}
        </div>
      )}

      {uploadProgress.length > 0 && (
        <div className="mt-4">
          {uploadProgress.map(({ fileName, progress }) => (
            <div key={fileName} className="mt-2">
              <div className="flex justify-between text-sm">
                <span className="truncate">{fileName}</span>
                <span>{Math.round(progress)}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-blue-500 rounded-full h-2 transition-all duration-300"
                  style={{ width: `${progress}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      )}

      <button
        onClick={handleUpload}
        disabled={isUploading || files.length === 0}
        className={`
          mt-4 w-full px-4 py-2 rounded font-medium
          ${isUploading || files.length === 0
            ? 'bg-gray-300 cursor-not-allowed'
            : 'bg-blue-500 hover:bg-blue-600 text-white'}
        `}
      >
        {isUploading ? 'Uploading...' : 'Upload Files'}
      </button>
    </div>
  )
}