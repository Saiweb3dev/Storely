// pages/auth.tsx
"use client"
import { useState } from 'react';
import { AnimatePresence } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { SignUpForm } from './SignUpForm';
import { SignInForm } from './SignInForm';
import { signIn,signUp } from '@/app/api/auth';
import { SignInFormData, SignUpFormData } from '@/types/auth';
import { setAuthToken } from '@/utils/token';
import { useAuth } from '@/contexts/AuthContext';
export default function AuthPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSignUp, setIsSignUp] = useState(false);
  const router = useRouter();
  const { login } = useAuth();

  const handleSignIn = async (data: SignInFormData) => {
    setIsLoading(true);
    try {
      const response = await signIn(data);
      setAuthToken(response.token);
      login(response.token, response.user);
      router.push('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign in failed');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignUp = async (data: SignUpFormData) => {
    setIsLoading(true);
    try {
      const response = await signUp(data);
      login(response.token, response.user);
      router.push('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign up failed');
    } finally {
      setIsLoading(false);
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
              onSwitch={() => setIsSignUp(false)}
              isLoading={isLoading}
              error={error}
            />
          ) : (
            <SignInForm
              key="signin"
              onSubmit={handleSignIn}
              onSwitch={() => setIsSignUp(true)}
              isLoading={isLoading}
              error={error}
            />
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}