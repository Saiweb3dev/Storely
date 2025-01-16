import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, AlertCircle } from 'lucide-react';
import { useRouter } from 'next/navigation';

interface CustomModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}

interface UnauthorizedModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CustomModal: React.FC<CustomModalProps> = ({ 
  isOpen, 
  onClose, 
  title, 
  children 
}) => {
  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed text-black inset-0 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 bg-black bg-opacity-50"
            onClick={onClose}
          />
          <motion.div
            initial={{ scale: 0.95, opacity: 0, y: 20 }}
            animate={{ scale: 1, opacity: 1, y: 0 }}
            exit={{ scale: 0.95, opacity: 0, y: 20 }}
            transition={{ duration: 0.2 }}
            className="relative bg-white rounded-lg shadow-lg w-full max-w-md mx-4 z-50"
          >
            <div className="p-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold">{title}</h2>
                <button
                  onClick={onClose}
                  className="p-1 hover:bg-gray-100 rounded-full transition-colors"
                  type="button"
                >
                  <X className="h-5 w-5 text-gray-500" />
                </button>
              </div>
              {children}
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
};

export const UnauthorizedModal: React.FC<UnauthorizedModalProps> = ({ 
  isOpen, 
  onClose 
}) => {
  const router = useRouter();

  const handleLogin = async (): Promise<void> => {
    // Close the modal first
    onClose();
    
    // Add a small delay to allow the modal animation to complete
    setTimeout(() => {
      router.push('/auth');
    }, 200); // This matches the duration of our modal animation
  };

  return (
    <CustomModal
      isOpen={isOpen}
      onClose={onClose}
      title="Not Logged In"
    >
      <div className="flex flex-col items-center gap-4 py-4">
        <div className="flex items-center justify-center w-12 h-12 bg-red-100 rounded-full">
          <AlertCircle className="h-6 w-6 text-red-500" />
        </div>
        <p className="text-center text-gray-600">
          Please log in to access this feature
        </p>
        <button
          onClick={handleLogin}
          type="button"
          className="w-full bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-lg transition-colors"
        >
          Go to Login
        </button>
      </div>
    </CustomModal>
  );
};