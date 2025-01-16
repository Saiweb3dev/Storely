// utils/authUtils.ts
interface StoredUserData {
  userID: string;
  username: string;
  email: string;
  storageUsed: number;
  storageLimit: number;
}

export const authUtils = {
  TOKEN_KEY: 'authToken',
  USER_KEY: 'userData',
  EXPIRY_KEY: 'tokenExpiry',

  setAuth(token: string, userData: StoredUserData) {
    const expiry = new Date();
    expiry.setHours(expiry.getHours() + 24);

    // Ensure consistent casing when storing
    const normalizedUserData = {
      userID: userData.userID || '',
      username: userData.username || '',
      email: userData.email || '',
      storageUsed: userData.storageUsed || 0,
      storageLimit: userData.storageLimit || 10
    };
    
    localStorage.setItem(this.TOKEN_KEY, token);
    localStorage.setItem(this.EXPIRY_KEY, expiry.toISOString());
    localStorage.setItem(this.USER_KEY, JSON.stringify(normalizedUserData));
  },

  getAuth() {
    try {
      const token = localStorage.getItem(this.TOKEN_KEY);
      const expiry = localStorage.getItem(this.EXPIRY_KEY);
      const userDataString = localStorage.getItem(this.USER_KEY);

      if (!token || !expiry || !userDataString) {
        return null;
      }

      if (new Date() > new Date(expiry)) {
        this.clearAuth();
        return null;
      }

      const userData = JSON.parse(userDataString);
      return {
        token,
        userData: {
          userID: userData.userID || '',
          username: userData.username || '',
          email: userData.email || '',
          storageUsed: userData.storageUsed || 0,
          storageLimit: userData.storageLimit || 10
        }
      };
    } catch (error) {
      console.error('Error parsing auth data:', error);
      this.clearAuth();
      return null;
    }
  },

  clearAuth() {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.EXPIRY_KEY);
    localStorage.removeItem(this.USER_KEY);
  },

  isLoggedIn() {
    return !!this.getAuth();
  }
};