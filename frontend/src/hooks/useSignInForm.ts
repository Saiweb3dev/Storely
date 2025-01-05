import { useState } from 'react'
import { SignInFormData } from '../types/auth'

export function useSignInForm() {
  const [formData, setFormData] = useState<SignInFormData>({
    nameOrEmail: '',
    password: '',
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prevData => ({ ...prevData, [name]: value }))
  }

  const resetForm = () => {
    setFormData({ nameOrEmail: '', password: '' })
  }

  return { formData, handleChange, resetForm }
}

