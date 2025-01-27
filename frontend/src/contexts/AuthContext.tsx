// contexts/AuthContext.tsx
"use client"
import { createContext, useContext, useState, useEffect } from 'react';
import { authUtils } from '@/utils/authUtils';

interface AuthContextType {
  isLoggedIn: boolean;
  userData: any;
  login: (token: string, userData: any) => void;
  logout: () => void;
}

interface UserData {
  userID: string;
  username: string;
  email: string;
  storageUsed: number;
  storageLimit: number;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userData, setUserData] = useState({
    userID : '',
    username: '',
    email: '',
    storageUsed: 0,
    storageLimit: 10,
  });

  useEffect(() => {
    const auth = authUtils.getAuth();
    if (auth?.userData) {
      setIsLoggedIn(true);
      setUserData(auth.userData);
    }
  }, []);

  const login = (token: string, userData: any) => {
    const sanitizedUserData: UserData = {
      userID : userData.userID || userData.userId || '',
      username: userData.username || '',
      email: userData.email || '',
      storageUsed: userData.storageUsed || 0,
      storageLimit: userData.storageLimit || 10,
    };
    console.log("Storing user data:", sanitizedUserData)
    authUtils.setAuth(token, sanitizedUserData);
    setIsLoggedIn(true);
    setUserData(sanitizedUserData);
  };

  const logout = () => {
    authUtils.clearAuth();
    setIsLoggedIn(false);
    setUserData({
      userID: '',
      username: '',
      email:'',
      storageUsed: 0,
      storageLimit: 10,
    });
  };

  return (
    <AuthContext.Provider value={{ isLoggedIn, userData, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};