"use client"
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Cloud, User } from 'lucide-react'

const Navbar = () => {
  const pathname = usePathname()

  return (
    <nav className="bg-white shadow-md p-4">
      <div className="container mx-auto flex items-center justify-between">
        <Link href="/" className="flex items-center space-x-2">
          <Cloud className="h-8 w-8 text-blue-500" />
          <span className="text-xl font-bold text-gray-800">Storely</span>
        </Link>
        <div className="flex space-x-4">
          <Link 
            href="/fileUpload" 
            className={`text-gray-600 hover:text-blue-500 ${pathname === '/files' ? 'font-bold' : ''}`}
          >
            Upload
          </Link>
          <Link 
            href="/fileDownload" 
            className={`text-gray-600 hover:text-blue-500 ${pathname === '/files' ? 'font-bold' : ''}`}
          >
            Download
          </Link>
        </div>
        <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
          <Link href="/SignIn">
          <User className="h-5 w-5 text-gray-600" />
          </Link>
        </div>
      </div>
    </nav>
  )
}

export default Navbar

