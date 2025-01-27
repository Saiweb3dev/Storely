// contexts/UploadsContext.tsx
"use client"
import React, { createContext, useContext, useState, useEffect } from 'react';
import { RecentUpload } from '@/types/upload';

interface UploadsContextType {
  recentUploads: RecentUpload[];
  addUploads: (uploads: RecentUpload[], userID: string) => void;
  clearAll: () => void;
  keepLastN: (n: number) => void;
  removeUpload: (fileId: string) => void;
}

const UploadsContext = createContext<UploadsContextType | undefined>(undefined);

export function UploadsProvider({ children }: { children: React.ReactNode }) {
  const [recentUploads, setRecentUploads] = useState<RecentUpload[]>([]);

  useEffect(() => {
    const stored = localStorage.getItem('recentUploads');
    if (stored) {
      try {
        const uploads = JSON.parse(stored);
        setRecentUploads(uploads.map((upload: any) => ({
          ...upload,
          uploadedAt: new Date(upload.uploadedAt)
        })));
      } catch (error) {
        console.error('Failed to parse recent uploads:', error);
      }
    }
  }, []);

   const addUploads = (newUploads: RecentUpload[], userId: string) => {
    setRecentUploads(prev => {
      // assign userId to each new upload
      const uploadsWithUser = newUploads.map(u => ({ ...u, userId }));
      const updated = [...uploadsWithUser, ...prev];
      localStorage.setItem('recentUploads', JSON.stringify(updated));
      return updated;
    });
  }

  const clearAll = () => {
    setRecentUploads([]);
    localStorage.removeItem('recentUploads');
  };

  const keepLastN = (n: number) => {
    setRecentUploads(prev => {
      const updated = prev.slice(0, n);
      localStorage.setItem('recentUploads', JSON.stringify(updated));
      return updated;
    });
  };

  const removeUpload = (fileId : string) => {
    setRecentUploads((prev) => {
      const updated = prev.filter((item) => item.fileId !== fileId);
      localStorage.setItem('recentUploads', JSON.stringify(updated));
      return updated;
    })
  }

  return (
    <UploadsContext.Provider value={{ recentUploads, addUploads, clearAll, keepLastN,removeUpload }}>
      {children}
    </UploadsContext.Provider>
  );
}

export const useUploads = () => {
  const context = useContext(UploadsContext);
  if (context === undefined) {
    throw new Error('useUploads must be used within a UploadsProvider');
  }
  return context;
};