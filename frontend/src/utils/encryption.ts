import CryptoJS from 'crypto-js';

const SECRET_KEY = process.env.NEXT_PUBLIC_ENCRYPTION_KEY;

export const encryptData = (data: any): string => {
  // Convert data to string and base64 encode
  const jsonString = JSON.stringify(data);
  const base64 = btoa(jsonString);
  
  // Create hash
  const hash = CryptoJS.SHA256(base64 + SECRET_KEY).toString();
  
  // Combine base64 data and hash
  return `${base64}.${hash}`;
};

export const decryptData = (encryptedData: string): any => {
  try {
    // Split data and hash
    const [base64Data, hash] = encryptedData.split('.');
    
    // Verify hash
    const calculatedHash = CryptoJS.SHA256(base64Data + SECRET_KEY).toString();
    if (hash !== calculatedHash) {
      throw new Error('Invalid data');
    }
    
    // Decode data
    const jsonString = atob(base64Data);
    return JSON.parse(jsonString);
  } catch (error) {
    if (error instanceof Error) {
      throw new Error('Decryption failed: ' + error.message);
    } else {
      throw new Error('Decryption failed');
    }
  }
};