"use client"
import { useState } from 'react';
import { AnimatePresence } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { SignUpForm } from './SignUpForm';
import { SignInForm } from './SignInForm';
import { signIn, signUp } from '@/app/api/auth';
import { SignInFormData, SignUpFormData } from '@/types/auth';
import { setAuthToken } from '@/utils/token';
import { useAuth } from '@/contexts/AuthContext';

interface AuthMessage {
  type: 'success' | 'error';
  content: string;
}

export default function AuthPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSignUp, setIsSignUp] = useState(false);
  const [message, setMessage] = useState<AuthMessage | null>(null);
  const router = useRouter();
  const { login } = useAuth();

  const handleSignIn = async (data: SignInFormData) => {
    setIsLoading(true);
    setError(null);
    setMessage(null);
    try {
      const response = await signIn(data);
      setAuthToken(response.token);
      login(response.token, response.user);
      router.push('/');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Sign in failed';
      setError(errorMessage);
      setMessage({ type: 'error', content: errorMessage });
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignUp = async (data: SignUpFormData) => {
    setIsLoading(true);
    setError(null);
    setMessage(null);
    try {
      const response = await signUp(data);
      console.log("Sign Up Response In Auth Main : ", response);
      
      if (response.message === "Registration successful") {
        setMessage({
          type: 'success',
          content: 'Registration successful! Please login to your account.'
        });
        setIsSignUp(false); // Switch to login form
      } else {
        throw new Error('Unexpected response from server');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Sign up failed';
      setError(errorMessage);
      setMessage({ type: 'error', content: errorMessage });
    } finally {
      setIsLoading(false);
    }
  };

  const handleFormSwitch = (toSignUp: boolean) => {
    setIsSignUp(toSignUp);
    setError(null);
    // Preserve success message when switching to login after registration
    if (!toSignUp && message?.type !== 'success') {
      setMessage(null);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md">
        <AnimatePresence mode="wait">
          {isSignUp ? (
            <SignUpForm
              key="signup"
              onSubmit={handleSignUp}
              onSwitch={() => handleFormSwitch(false)}
              isLoading={isLoading}
              error={message?.type === 'error' ? message.content : null}
            />
          ) : (
            <SignInForm
              key="signin"
              onSubmit={handleSignIn}
              onSwitch={() => handleFormSwitch(true)}
              isLoading={isLoading}
              error={message?.type === 'error' ? message.content : null}
              successMessage={message?.type === 'success' ? message.content : null}
            />
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}