"use client"
import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FileIcon, Copy, Check, Trash2, X, Image, FileText, File } from 'lucide-react';
import { authUtils } from '@/utils/authUtils'; // adjust the import path
import { useUploads } from '@/contexts/UploadsContext';

type UploadItem = {
  fileId: string;
  fileName: string;
  fileType: string;
  uploadedAt: Date;
  size?: number;
};

export default function DeletableUploads() {
  const { recentUploads, removeUpload } = useUploads();
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [confirmFileName, setConfirmFileName] = useState('');
  const [fileToDelete, setFileToDelete] = useState<UploadItem | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const copyToClipboard = async (fileId: string) => {
    try {
      await navigator.clipboard.writeText(fileId);
      setCopiedId(fileId);
      setTimeout(() => setCopiedId(null), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat('en-US', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  };

  const formatSize = (size: number) => {
    const units = ['B', 'KB', 'MB', 'GB'];
    let formattedSize = size;
    let unitIndex = 0;
    while (formattedSize >= 1024 && unitIndex < units.length - 1) {
      formattedSize /= 1024;
      unitIndex++;
    }
    return `${formattedSize.toFixed(2)} ${units[unitIndex]}`;
  };

  const handleDeleteClick = (upload: UploadItem) => {
    setFileToDelete(upload);
    setConfirmFileName('');
    setIsModalOpen(true);
  };

  const confirmDelete = async () => {
    if (!fileToDelete) return;
    const { fileId, fileName } = fileToDelete;

    if (confirmFileName !== fileName) {
      alert("Entered file name doesn't match.");
      return;
    }

    const { token, userID } = authUtils.getAuthTokenAndUserId();

    try {
      const response = await fetch('http://localhost:8080/api/minio/files/delete', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          userId: userID,
          fileId: fileId,
        }),
      });

      if (!response.ok) {
        const msg = await response.text();
        throw new Error(msg);
      }

      removeUpload(fileId);
      alert('File deleted successfully.');
    } catch (err) {
      console.error('Delete error:', err);
      alert('Failed to delete file.');
    } finally {
      setFileToDelete(null);
      setConfirmFileName('');
      setIsModalOpen(false);
    }
  };

  const getFileTypeIcon = (fileType: string) => {
    if (fileType.startsWith('image/')) {
      return <Image className="h-8 w-8 text-blue-500" />;
    } else if (fileType === 'application/pdf') {
      return <FileText className="h-8 w-8 text-red-500" />;
    } else {
      return <File className="h-8 w-8 text-gray-500" />;
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-4xl mx-auto mt-6 p-6 bg-white rounded-lg shadow-lg relative"
    >
      <h2 className="text-2xl text-black font-bold mb-6">Recent Uploads</h2>
      
      <AnimatePresence>
        {recentUploads.length === 0 ? (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="text-blue-500 text-center text-lg"
          >
            No recent uploads
          </motion.p>
        ) : (
          <motion.div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {recentUploads.map((upload) => (
              <motion.div
                key={upload.fileId}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                whileHover={{ scale: 1.03, transition: { duration: 0.2 } }}
                className="bg-gray-50 rounded-lg shadow-md overflow-hidden"
              >
                <div className="p-4">
                  <div className="flex items-center justify-between mb-2">
                    {getFileTypeIcon(upload.fileType)}
                    <div className="flex items-center space-x-2">
                      <motion.button
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.9 }}
                        onClick={() => copyToClipboard(upload.fileId)}
                        className="p-1 hover:bg-gray-200 rounded transition-colors"
                        title="Copy file ID"
                      >
                        {copiedId === upload.fileId ? (
                          <Check className="h-4 w-4 text-green-500" />
                        ) : (
                          <Copy className="h-4 w-4 text-gray-500" />
                        )}
                      </motion.button>
                      <motion.button
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.9 }}
                        onClick={() => handleDeleteClick(upload)}
                        className="p-1 hover:bg-gray-200 rounded transition-colors"
                        title="Delete file"
                      >
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </motion.button>
                    </div>
                  </div>
                  <h3 className="font-semibold text-black text-lg mb-2 truncate" title={upload.fileName}>
                    {upload.fileName}
                  </h3>
                  <div className="text-sm text-gray-600 space-y-1">
                    <p>{upload.fileType}</p>
                    <p>{formatDate(new Date(upload.uploadedAt))}</p>
                    {upload.size && <p>{formatSize(upload.size)}</p>}
                  </div>
                </div>
              </motion.div>
            ))}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete confirmation modal */}
      <AnimatePresence>
        {isModalOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-white rounded-lg p-6 w-96 relative"
            >
              <button
                onClick={() => setIsModalOpen(false)}
                className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
              >
                <X className="h-6 w-6" />
              </button>
              <h3 className="text-lg font-semibold mb-4">Confirm Deletion</h3>
              <p className="text-red-600 font-semibold mb-2">
                Type the file name <strong>{fileToDelete?.fileName}</strong> to confirm deletion:
              </p>
              <input
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded mb-4"
                value={confirmFileName}
                onChange={(e) => setConfirmFileName(e.target.value)}
                placeholder="Enter file name"
              />
              <div className="flex justify-end space-x-2">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 border text-black border-gray-300 rounded"
                >
                  Cancel
                </motion.button>
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={confirmDelete}
                  className="px-4 py-2 bg-red-600 text-white rounded disabled:opacity-50"
                  disabled={confirmFileName !== fileToDelete?.fileName}
                >
                  Confirm Delete
                </motion.button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
}

