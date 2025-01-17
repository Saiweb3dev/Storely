'use client'

import React, { useEffect, useState } from "react"
import axios from "axios"
import { Database } from 'lucide-react'
import { authUtils } from "@/utils/authUtils"
import { useAuth } from "@/contexts/AuthContext"
import { bytesToMB, bytesToGB } from "@/utils/dataSizeUtils"

export default function StorageHealth() {
  const [health, setHealth] = useState<{storageUsed:number, storageLimit:number, availableBalance:number} | null>(null)
  const [error, setError] = useState<string | null>(null)
  const { userData } = useAuth()
  const userID = userData?.userID
  const { token } = authUtils.getAuthTokenAndUserId()

  useEffect(() => {
    if (!token || !userID) return
  
    axios.get(
      "http://localhost:8080/get/user/storageHealth",
      {
        headers: { Authorization: `Bearer ${token}` },
        params: { userID }
      }
    )
    .then(res => setHealth(res.data))
    .catch(err => setError(err.response?.data || "Error fetching storage health"))
  }, [token, userID])

  if (error) return <div className="text-red-500">{error}</div>
  if (!health) return <div className="text-gray-500">Loading...</div>

  const usedMB = bytesToMB(health.storageUsed)
  const limitMB = bytesToMB(health.storageLimit)
  const availableMB = bytesToMB(health.availableBalance)

  return (
    <div className="w-full max-w-md bg-white shadow-lg rounded-lg p-6 border border-gray-200">
      <h2 className="text-2xl font-bold mb-4 flex items-center text-black">
        <Database className="mr-2 text-blue-600" />
        Your Storage
      </h2>
      <div className="mb-4">
        <p className="font-semibold text-black mb-2">Storage Usage</p>
        <div className="flex justify-between mt-1 text-sm">
          <span className="text-blue-600">{usedMB} MB used</span>
          <span className="text-green-600">{availableMB} MB left</span>
        </div>
      </div>
      <div className="flex justify-between text-sm text-black">
        <span>Total Storage: {limitMB} MB</span>
      </div>
    </div>
  )
}