'use client'

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { Cloud, User } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import UserProfileModal from './Modal/UserProfileModal';
import { UnauthorizedModal } from './Modal/UnauthorizedModal';

// You might want to move these interfaces to a separate types file
interface UserData {
  userID?: string;
  username?: string;
  storageUsed?: number;
  storageLimit?: number;
}

interface AuthContextType {
  isLoggedIn: boolean;
  userData: UserData;
  logout: () => void;
}

const Navbar: React.FC = () => {
  const pathname = usePathname();
  const router = useRouter();
  const { isLoggedIn, userData, logout } = useAuth() as AuthContextType;
  
  const [isProfileModalOpen, setIsProfileModalOpen] = useState<boolean>(false);
  const [isUnauthorizedModalOpen, setIsUnauthorizedModalOpen] = useState<boolean>(false);

  const handleNavigation = (path: string): void => {
    if (!isLoggedIn) {
      setIsUnauthorizedModalOpen(true);
      return;
    }
    router.push(path);
  };

  const handleUserIconClick = (): void => {
    if (isLoggedIn) {
      setIsProfileModalOpen(true);
    } else {
      router.push('/auth');
    }
  };

  const handleDisconnect = (): void => {
    logout();
    setIsProfileModalOpen(false);
  };

  return (
    <nav className="bg-white shadow-md p-4">
      <div className="container mx-auto flex items-center justify-between">
        <Link href="/" className="flex items-center space-x-2">
          <Cloud className="h-8 w-8 text-blue-500" />
          <span className="text-xl font-bold text-gray-800">Storely</span>
        </Link>
        <div className="flex space-x-4">
          <button
            onClick={() => handleNavigation('/fileUpload')}
            type="button"
            className={`text-gray-600 hover:text-blue-500 ${
              pathname === '/fileUpload' ? 'font-bold' : ''
            }`}
          >
            Upload
          </button>
          <button
            onClick={() => handleNavigation('/fileDownload')}
            type="button"
            className={`text-gray-600 hover:text-blue-500 ${
              pathname === '/fileDownload' ? 'font-bold' : ''
            }`}
          >
            Download
          </button>
        </div>
        <div
          className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center cursor-pointer"
          onClick={handleUserIconClick}
          role="button"
          tabIndex={0}
        >
          <User className="h-5 w-5 text-gray-600" />
        </div>
      </div>
      
      <UserProfileModal
        isOpen={isProfileModalOpen}
        onClose={() => setIsProfileModalOpen(false)}
        username={userData?.username || 'No username'}
        onDisconnect={handleDisconnect}
        storageUsed={userData?.storageUsed || 0}
        storageLimit={userData?.storageLimit || 10}
      />
      
      <UnauthorizedModal
        isOpen={isUnauthorizedModalOpen}
        onClose={() => setIsUnauthorizedModalOpen(false)}
      />
    </nav>
  );
};

export default Navbar;