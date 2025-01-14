"use client"
import React, { useState, useEffect } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Cloud, User } from 'lucide-react'
import UserProfileModal from './User/UserProfileModal'
import { useAuth } from '@/contexts/AuthContext'

const Navbar = () => {
  const pathname = usePathname()
  const { isLoggedIn, userData, logout } = useAuth();
  const [isModalOpen, setIsModalOpen] = useState(false)

  const handleUserIconClick = () => {
    if (isLoggedIn) {
      setIsModalOpen(true);
    } else {
      window.location.href = '/auth';
    }
  };

  const handleDisconnect = () => {
    logout();
    setIsModalOpen(false);
  };

  return (
    <nav className="bg-white shadow-md p-4">
      <div className="container mx-auto flex items-center justify-between">
        <Link href="/" className="flex items-center space-x-2">
          <Cloud className="h-8 w-8 text-blue-500" />
          <span className="text-xl font-bold text-gray-800">Storely</span>
        </Link>
        <div className="flex space-x-4">
          <Link 
            href="/fileUpload" 
            className={`text-gray-600 hover:text-blue-500 ${pathname === '/fileUpload' ? 'font-bold' : ''}`}
          >
            Upload
          </Link>
          <Link 
            href="/fileDownload" 
            className={`text-gray-600 hover:text-blue-500 ${pathname === '/fileDownload' ? 'font-bold' : ''}`}
          >
            Download
          </Link>
        </div>
        <div 
          className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center cursor-pointer"
          onClick={handleUserIconClick}
        >
          <User className="h-5 w-5 text-gray-600" />
        </div>
      </div>
      <UserProfileModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        username={userData.username || 'No username'}
        onDisconnect={handleDisconnect}
        storageUsed={userData?.storageUsed || 0}
        storageLimit={userData?.storageLimit || 10}
      />
    </nav>
  )
}

export default Navbar

