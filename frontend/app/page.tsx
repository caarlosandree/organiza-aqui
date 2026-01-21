'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'

export default function HomePage() {
  const router = useRouter()
  const { isAuthenticated, token } = useAuthStore()

  useEffect(() => {
    if (isAuthenticated || token) {
      router.replace('/tasks')
    } else {
      router.replace('/login')
    }
  }, [isAuthenticated, token, router])

  return null
}
