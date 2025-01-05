// components/SignInPage.tsx
'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { User, Lock, ArrowRight, AlertCircle, UserPlus, LogIn } from 'lucide-react'
import { useSignInForm } from '../hooks/useSignInForm'
import { signIn, register } from '../app/api/auth'
import { setAuthToken } from '../utils/token'
import { useRouter } from 'next/navigation'

export default function SignInPage() {
  const { formData, handleChange, resetForm } = useSignInForm()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [isSignUp, setIsSignUp] = useState(false)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      if (isSignUp) {
        const response = await register(formData)
        console.log('Registration successful')
        // Automatically sign in after registration
        const signInResponse = await signIn(formData)
        setAuthToken(signInResponse.token)
      } else {
        const response = await signIn(formData)
        setAuthToken(response.token)
      }
      
      console.log(`${isSignUp ? 'Sign-up' : 'Sign-in'} successful`)
      resetForm()
      router.push('/')
    } catch (error) {
      setError(error instanceof Error ? error.message : `${isSignUp ? 'Registration' : 'Sign-in'} failed`)
      console.error(`${isSignUp ? 'Registration' : 'Sign-in'} failed:`, error)
    } finally {
      setIsLoading(false)
    }
  }

  const toggleMode = () => {
    setIsSignUp(!isSignUp)
    setError(null)
    resetForm()
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="bg-white p-8 rounded-lg shadow-lg w-full max-w-md"
      >
        <motion.h1 
          key={isSignUp ? 'signup' : 'signin'}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-3xl font-bold text-center text-blue-600 mb-6"
        >
          {isSignUp ? 'Create Account' : 'Sign In'}
        </motion.h1>

        <AnimatePresence mode='wait'>
          {error && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="mb-4 p-3 bg-red-100 text-red-700 rounded-md flex items-center"
            >
              <AlertCircle className="mr-2" size={20} />
              {error}
            </motion.div>
          )}
        </AnimatePresence>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="relative">
            <User className="absolute top-3 left-3 text-gray-600" size={20} />
            <input
              type="email"
              name="nameOrEmail"
              value={formData.nameOrEmail}
              onChange={handleChange}
              placeholder="Email"
              className="w-full pl-10 pr-4 py-2 border text-black border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div className="relative">
            <Lock className="absolute top-3 left-3 text-gray-400" size={20} />
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Password"
              className="w-full pl-10 pr-4 py-2 border text-black border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            type="submit"
            disabled={isLoading}
            className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-300 flex items-center justify-center"
          >
            {isLoading ? (
              `${isSignUp ? 'Creating Account...' : 'Signing In...'}`
            ) : (
              <>
                {isSignUp ? (
                  <>
                    Create Account
                    <UserPlus className="ml-2" size={20} />
                  </>
                ) : (
                  <>
                    Sign In
                    <LogIn className="ml-2" size={20} />
                  </>
                )}
              </>
            )}
          </motion.button>
        </form>

        <div className="mt-4 text-center">
          <button
            onClick={toggleMode}
            className="text-blue-600 hover:text-blue-700 focus:outline-none"
          >
            {isSignUp ? 'Already have an account? Sign In' : "Don't have an account? Sign Up"}
          </button>
        </div>
      </motion.div>
    </div>
  )
}