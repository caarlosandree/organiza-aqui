'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const { isAuthenticated, token, user } = useAuthStore()
  const [isChecking, setIsChecking] = useState(true)

  useEffect(() => {
    // Aguardar um pouco para o Zustand restaurar o estado do localStorage
    const checkAuth = () => {
      if (typeof window !== 'undefined') {
        const storedToken = localStorage.getItem('token')
        // Verificar se há token e usuário (autenticação completa)
        const hasToken = token || storedToken
        const hasUser = !!user
        const hasAuth = (isAuthenticated || (hasToken && hasUser))

        if (!hasAuth) {
          router.replace('/login')
        } else {
          setIsChecking(false)
        }
      } else {
        setIsChecking(false)
      }
    }

    // Aguardar um tick para garantir que o Zustand restaurou o estado
    const timer = setTimeout(checkAuth, 100)
    return () => clearTimeout(timer)
  }, [isAuthenticated, token, user, router])

  // Aguardar verificação inicial antes de renderizar
  if (isChecking) {
    return null
  }

  // Verificar novamente antes de renderizar
  if (typeof window !== 'undefined') {
    const storedToken = localStorage.getItem('token')
    const hasToken = token || storedToken
    const hasUser = !!user
    const hasAuth = isAuthenticated || (hasToken && hasUser)

    if (!hasAuth) {
      return null
    }
  } else if (!isAuthenticated && !token) {
    return null
  }

  return <>{children}</>
}
