// app/api/auth.ts
import { SignInFormData,SignUpFormData, AuthResponse } from '../../types/auth'
import { encryptData,decryptData } from '@/utils/encryption';

export async function signUp(data: SignUpFormData): Promise<AuthResponse> {
  try {
    const encryptedData = encryptData(data);
    
    const response = await fetch('http://localhost:8080/api/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ data: encryptedData })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Registration failed');
    }
    console.log("Response: ", response)

    const encrypted = await response.json();
    if (encrypted.data) {
      const decrypted = decryptData(encrypted.data);

      console.log("Decrypted: ", decrypted)
      
      // Ensure the response matches the AuthResponse interface
      return {
        token: decrypted.token,
        user: {
          username: decrypted.user?.username || data.username,
          email: decrypted.user?.email || data.email,
          storageUsed: decrypted.user?.storageUsed || 0,
          storageLimit: decrypted.user?.storageLimit || 10
        }
      };
    }
    throw new Error('Invalid response format');
  } catch (error) {
    console.error('SignUp error:', error);
    throw error;
  }
}

export async function signIn(data: Pick<SignInFormData, 'email' | 'password'>): Promise<AuthResponse> {
  try {
    const encryptedData = encryptData(data);
    
    const response = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Content-Encrypted': 'true'
      },
      body: JSON.stringify({ data: encryptedData })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Login failed');
    }

    const encrypted = await response.json();
    // Decrypt the response data to get the token
    return decryptData(encrypted.data);
  } catch (error) {
    throw error;
  }
}