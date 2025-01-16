import React from 'react';
import { X } from 'lucide-react';
import { bytesToMB,bytesToGB } from '@/utils/dataSizeUtils';
interface UserProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  username: string;
  onDisconnect: () => void;
  storageUsed: number;
  storageLimit: number;
}

const UserProfileModal: React.FC<UserProfileModalProps> = ({
  isOpen,
  onClose,
  username,
  onDisconnect,
  storageUsed,
  storageLimit
}) => {
  if (!isOpen) return null;

  const usagePercentage = (storageUsed / storageLimit) * 100;
  const storageLeft = storageLimit - storageUsed;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-80 text-gray-400">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl text-black font-bold">User Profile</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <X size={24} />
          </button>
        </div>
        <div className="mb-4">
          <p className="font-semibold">Username: <span className='text-black'>{username}</span></p>
        </div>
        <div className="mb-4">
          <p className="font-semibold mb-2">Storage Usage</p>
          <div className="w-full bg-gray-200 rounded-full h-2.5">
            <div
              className="bg-blue-600 h-2.5 rounded-full"
              style={{ width: `${usagePercentage}%` }}
            ></div>
          </div>
          <div className="flex justify-between mt-1 text-sm">
            <span className='text-blue-600'>{bytesToMB(storageUsed)} used</span>
            <span>{bytesToMB(storageLeft)} left</span>
          </div>
          <p className="text-center mt-2">
            {usagePercentage.toFixed(1)}% of {bytesToGB(storageLimit)} used
          </p>
        </div>
        <button
          onClick={onDisconnect}
          className="w-full bg-red-500 text-white py-2 rounded hover:bg-red-600 transition duration-300"
        >
          Disconnect
        </button>
      </div>
    </div>
  );
};

export default UserProfileModal;

