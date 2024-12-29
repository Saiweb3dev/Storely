// components/ClearOptions.tsx
import { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Trash2 } from 'lucide-react';
import { useUploads } from '@/contexts/UploadsContext';

export default function ClearOptions() {
  const [isOpen, setIsOpen] = useState(false);
  const { clearAll, keepLastN } = useUploads();
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const menuVariants = {
    hidden: { opacity: 0, scale: 0.95, y: -10 },
    visible: { opacity: 1, scale: 1, y: 0 }
  };

  return (
    <div className="relative" ref={menuRef}>
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => setIsOpen(!isOpen)}
        className="p-2 hover:bg-gray-100 rounded-full transition-colors"
      >
        <Trash2 className="h-5 w-5 text-gray-500" />
      </motion.button>

      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial="hidden"
            animate="visible"
            exit="hidden"
            variants={menuVariants}
            transition={{ duration: 0.2 }}
            className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 z-50"
          >
            <div className="py-1">
              <motion.button
                whileHover={{ backgroundColor: '#F3F4F6' }}
                onClick={() => {
                  keepLastN(5);
                  setIsOpen(false);
                }}
                className="block w-full text-left px-4 py-2 text-sm text-gray-700"
              >
                Keep Last 5 Uploads
              </motion.button>
              <motion.button
                whileHover={{ backgroundColor: '#F3F4F6' }}
                onClick={() => {
                  clearAll();
                  setIsOpen(false);
                }}
                className="block w-full text-left px-4 py-2 text-sm text-red-600"
              >
                Clear All
              </motion.button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}