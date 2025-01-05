// app/api/auth.ts
import { SignInFormData, AuthResponse } from '../../types/auth'

export async function signIn(data: SignInFormData): Promise<AuthResponse> {
  try {
    const response = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: data.nameOrEmail,
        password: data.password,
      }),
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Login failed')
    }

    const authData = await response.json()
    return authData
  } catch (error) {
    throw error
  }
}

export async function register(data: SignInFormData): Promise<AuthResponse> {
  try {
    const response = await fetch('http://localhost:8080/api/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: data.nameOrEmail,
        password: data.password,
      }),
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Registration failed')
    }

    return await response.json()
  } catch (error) {
    throw error
  }
}