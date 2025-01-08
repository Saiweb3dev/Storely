// app/api/auth.ts
import { SignInFormData,SignUpFormData, AuthResponse } from '../../types/auth'
import { encryptData,decryptData } from '@/utils/encryption';
import { UserData } from '../../types/auth';

export async function register(data: UserData): Promise<AuthResponse> {
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
      throw new Error('Registration failed');
    }

    return await response.json();
  } catch (error) {
    throw error;
  }
}

// app/api/auth.ts
export async function signIn(data: Pick<UserData, 'email' | 'password'>): Promise<AuthResponse> {
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