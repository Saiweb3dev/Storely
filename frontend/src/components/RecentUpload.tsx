// components/RecentUploads.tsx
"use client"
import { motion } from 'framer-motion';
import { FileIcon, Download, Copy, Check } from 'lucide-react';
import { useUploads } from '@/contexts/UploadsContext';
import ClearOptions from './ClearOptions';
import { useState } from 'react';

export default function RecentUploads() {
  const { recentUploads } = useUploads();
  const [copiedId, setCopiedId] = useState<string | null>(null);

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

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-md mx-auto mt-6 p-6 bg-white rounded-lg shadow-lg"
    >
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl text-black font-bold">Recent Uploads</h2>
        <ClearOptions />
      </div>
      
      {recentUploads.length === 0 ? (
        <p className="text-blue-500 text-center">No recent uploads</p>
      ) : (
        <div className="space-y-3">
          {recentUploads.map((upload) => (
            <motion.div
              key={upload.fileId}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="flex text-black items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div className="flex items-center space-x-3">
                <FileIcon className="h-5 w-5 text-blue-500" />
                <div>
                  <p className="font-medium truncate max-w-[200px]">
                    {upload.fileName}
                  </p>
                  <div className="flex items-center space-x-2 text-sm text-blue-500">
                    <span>{upload.fileType}</span>
                    <span>•</span>
                    <span>{formatDate(new Date(upload.uploadedAt))}</span>
                  </div>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => copyToClipboard(upload.fileId)}
                  className="p-1 hover:bg-gray-200 rounded transition-colors"
                  title="Copy file ID"
                >
                  {copiedId === upload.fileId ? (
                    <Check className="h-4 w-4 text-green-500" />
                  ) : (
                    <Copy className="h-4 w-4 text-gray-500" />
                  )}
                </button>
                <button
                  onClick={() => window.location.href = `/files/${upload.fileId}`}
                  className="p-1 hover:bg-gray-200 rounded transition-colors"
                  title="Download file"
                >
                  <Download className="h-4 w-4 text-gray-500" />
                </button>
              </div>
            </motion.div>
          ))}
        </div>
      )}
    </motion.div>
  );
}