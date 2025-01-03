"use client"
import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import { uploadService } from "../service/uploadService"
import { motion, AnimatePresence } from 'framer-motion'
import { Upload, X, AlertCircle, File } from 'lucide-react'
import { RecentUpload, UploadProgressInfo } from '@/types/upload'
import { useUploads } from '@/contexts/UploadsContext';
import MinIODirectUpload from './MinIODirectUpload'
interface FileWithPreview extends File {
  preview: string;
}

export default function FileUpload() {
  const { addUploads } = useUploads();
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
      'text/plain': ['.txt'],
    }
  })

  const removeFile = (file: FileWithPreview) => {
    const newFiles = [...files]
    newFiles.splice(newFiles.indexOf(file), 1)
    setFiles(newFiles)
    setError(null)
  }

  // Update FileUpload.tsx handleUpload method
  const handleUpload = async () => {
    if (files.length === 0) {
      setError('Please select files to upload');
      return;
    }
  
    setIsUploading(true);
    setError(null);
  
    try {
      const updateProgress = (info: UploadProgressInfo) => {
        setUploadProgress(prev => 
          prev.map(p => p.fileName === info.fileName ? info : p)
        );
      };
  
      setUploadProgress(files.map(f => ({ fileName: f.name, progress: 0 })));
      
      const uploadResults = await uploadService.uploadFiles(files, updateProgress);
  
      const newUploads: RecentUpload[] = uploadResults.map(result => ({
        fileId: result.fileId,
        fileName: result.fileName,
        fileType: result.fileType,
        uploadedAt: new Date()
      }));

      addUploads(newUploads);
  
      const stored = localStorage.getItem('recentUploads');
      const existingUploads: RecentUpload[] = stored ? JSON.parse(stored) : [];
      const updatedUploads = [...newUploads, ...existingUploads].slice(0, 10);
      localStorage.setItem('recentUploads', JSON.stringify(updatedUploads));
  
      setFiles([]);
      setUploadProgress([]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setIsUploading(false);
    }
  };

  const handleMinIOUploadComplete = ({ fileId, name, type, size }: {
    fileId: string;
    name: string;
    type: string;
    size: number;
  }) => {
    addUploads([{
      fileId,
      fileName: name,
      fileType: type,
      uploadedAt: new Date(),
      size,
    }]);
  };

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
        <Upload className="mr-2 h-6 w-6 text-blue-500" />
        Upload Files
      </motion.h2>

      <div 
        {...getRootProps()} 
        className={`
          border-2 border-dashed rounded-lg p-6 cursor-pointer transition-colors duration-200
          ${isDragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-blue-400'}
        `}
      >
        <input {...getInputProps()} />
        {isDragActive ? (
          <p className="text-center text-blue-500">Drop the files here...</p>
        ) : (
          <p className="text-center text-gray-500">
            Drag & drop files here, or click to select files
          </p>
        )}
      </div>

      <AnimatePresence>
        {error && (
          <motion.div 
            className="mt-4 p-4 bg-red-100 text-red-700 rounded-md flex items-center"
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.3 }}
          >
            <AlertCircle className="mr-2 h-5 w-5" />
            {error}
          </motion.div>
        )}
      </AnimatePresence>

      <AnimatePresence>
        {files.length > 0 && (
          <motion.div 
            className="mt-4"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <h3 className="font-medium text-gray-700">Selected Files:</h3>
            {files.map(file => (
              <motion.div 
                key={file.name} 
                className="mt-2 flex items-center justify-between text-blue-500 bg-gray-50 p-2 rounded"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.2 }}
              >
                <span className="truncate flex-1 flex items-center">
                  <File className="mr-2 h-4 w-4 text-blue-500" />
                  {file.name}
                </span>
                <motion.button
                  onClick={() => removeFile(file)}
                  className="ml-2 text-red-500 hover:text-red-700"
                  disabled={isUploading}
                  whileHover={{ scale: 1.1 }}
                  whileTap={{ scale: 0.9 }}
                >
                  <X className="h-5 w-5" />
                </motion.button>
              </motion.div>
            ))}
          </motion.div>
        )}
      </AnimatePresence>

      <AnimatePresence>
        {uploadProgress.length > 0 && (
          <motion.div 
            className="mt-4"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            {uploadProgress.map(({ fileName, progress }) => (
              <div key={fileName} className="mt-2">
                <div className="flex justify-between text-sm text-gray-600">
                  <span className="truncate">{fileName}</span>
                  <span>{Math.round(progress)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <motion.div
                    className="bg-blue-500 rounded-full h-2"
                    initial={{ width: 0 }}
                    animate={{ width: `${progress}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>
            ))}
          </motion.div>
        )}
      </AnimatePresence>

      <div className="flex flex-col gap-2 p-4"> {/* Added container for buttons */}
  <motion.button
    onClick={handleUpload}
    disabled={isUploading || files.length === 0}
    className={`
      w-full px-4 py-2 rounded-md font-medium transition-colors duration-200
      ${isUploading || files.length === 0
        ? 'bg-gray-300 cursor-not-allowed'
        : 'bg-blue-500 hover:bg-blue-600 text-white'}
    `}
    whileHover={{ scale: 1.05 }}
    whileTap={{ scale: 0.95 }}
  >
    {isUploading ? 'Uploading...' : 'Upload Files'}
  </motion.button>
  
  <MinIODirectUpload
    onUploadComplete={handleMinIOUploadComplete}
    onError={setError}
  />
</div>
    </motion.div>
  )
}

