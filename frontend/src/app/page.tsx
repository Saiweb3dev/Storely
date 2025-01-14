'use client'

import { motion } from 'framer-motion'
import { Upload, FileText, Share2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
export default function Home() {
  const features = [
    { icon: Upload, title: 'Easy Upload', description: 'Drag and drop files to upload instantly' },
    { icon: FileText, title: 'File Management', description: 'Organize your files with ease' },
    { icon: Share2, title: 'Quick Sharing', description: 'Share files and folders with a single click' },
  ]

  const { isLoggedIn ,userData } = useAuth()
  console.log(userData)

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gradient-to-br from-blue-50 to-indigo-100">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
        className="text-center mb-12"
      >
        <h1 className="text-5xl font-bold text-gray-800 mb-4">Welcome to MiniDrive {isLoggedIn ? userData.username : ''}</h1>
        <p className="text-xl text-gray-600">Your personal cloud storage solution</p>
      </motion.div>

      <motion.div 
        className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-12"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.4, duration: 0.8 }}
      >
        {features.map((feature, index) => (
          <motion.div
            key={index}
            className="bg-white p-6 rounded-lg shadow-lg"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <feature.icon className="h-12 w-12 text-blue-500 mb-4 mx-auto" />
            <h2 className="text-xl font-semibold text-gray-800 mb-2">{feature.title}</h2>
            <p className="text-gray-600">{feature.description}</p>
          </motion.div>
        ))}
      </motion.div>

      <motion.button
        className="mt-12 px-6 py-3 bg-blue-500 text-white rounded-full font-semibold text-lg shadow-lg"
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.8, duration: 0.5 }}
      >
        Get Started
      </motion.button>
    </main>
  )
}

