import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Loader2, Mail, Lock } from 'lucide-react';
import { SignInFormData } from '@/types/auth';

interface SignInFormProps {
  onSubmit: (data: SignInFormData) => void;
  onSwitch: () => void;
  isLoading: boolean;
  error: string | null;
}

export function SignInForm({ onSubmit, onSwitch, isLoading, error }: SignInFormProps) {
  const [formData, setFormData] = useState<SignInFormData>({
    email: '',
    password: ''
  });

  const [focusedField, setFocusedField] = useState<'email' | 'password' | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      transition={{ duration: 0.3 }}
      className="w-full max-w-md mx-auto"
    >
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2, duration: 0.5 }}
        className="mb-8 text-center"
      >
        <motion.h2
          className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-blue-500 bg-300% animate-gradient"
        >
          Login to your account
        </motion.h2>
        <motion.div
          className="mt-2 h-1 w-16 bg-blue-500 mx-auto rounded-full"
          initial={{ scaleX: 0 }}
          animate={{ scaleX: 1 }}
          transition={{ delay: 0.5, duration: 0.5 }}
        />
      </motion.div>

      <motion.form
        onSubmit={handleSubmit}
        className="space-y-6"
        animate={error ? { x: [-10, 10, -10, 10, 0] } : {}}
        transition={{ type: 'spring', stiffness: 500, damping: 10 }}
      >
        <div className="space-y-4">
          <motion.div
            className="relative"
            whileFocus={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              onFocus={() => setFocusedField('email')}
              onBlur={() => setFocusedField(null)}
              placeholder="Email"
              className="w-full text-black pl-10 pr-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={isLoading}
            />
            <AnimatePresence>
              {focusedField === 'email' && (
                <motion.div
                  initial={{ scaleX: 0 }}
                  animate={{ scaleX: 1 }}
                  exit={{ scaleX: 0 }}
                  className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500"
                  style={{ originX: 0 }}
                />
              )}
            </AnimatePresence>
          </motion.div>
          <motion.div
            className="relative"
            whileFocus={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            <input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              onFocus={() => setFocusedField('password')}
              onBlur={() => setFocusedField(null)}
              placeholder="Password"
              className="w-full text-black pl-10 pr-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={isLoading}
            />
            <AnimatePresence>
              {focusedField === 'password' && (
                <motion.div
                  initial={{ scaleX: 0 }}
                  animate={{ scaleX: 1 }}
                  exit={{ scaleX: 0 }}
                  className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500"
                  style={{ originX: 0 }}
                />
              )}
            </AnimatePresence>
          </motion.div>
        </div>

        <AnimatePresence>
          {error && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="text-red-500 text-sm"
            >
              {error}
            </motion.div>
          )}
        </AnimatePresence>

        <motion.button
          type="submit"
          className="w-full bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
          disabled={isLoading}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          {isLoading ? (
            <motion.div
              className="flex items-center justify-center"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
            >
              <Loader2 className="animate-spin mr-2" />
              Signing In...
            </motion.div>
          ) : (
            <motion.span
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              Sign In
            </motion.span>
          )}
        </motion.button>
      </motion.form>

      <motion.div
        className="mt-6 text-center"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
      >
        <button
          onClick={onSwitch}
          disabled={isLoading}
          className="text-sm text-blue-600 hover:text-blue-700 focus:outline-none focus:underline"
        >
          Don't have an account? Create one
        </button>
      </motion.div>

      {/* <motion.div
        className="absolute inset-0 -z-10 bg-gradient-to-br from-blue-900 to-purple-600 opacity-50"
        animate={{
          backgroundPosition: ['0% 0%', '100% 100%'],
          transition: {
            duration: 20,
            ease: 'linear',
            repeat: Infinity,
            repeatType: 'reverse'
          }
        }}
      /> */}
    </motion.div>
  );
}

