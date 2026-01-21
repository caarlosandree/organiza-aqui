// Zustand stores
// Adicione suas stores aqui conforme necessário

import { create } from 'zustand'
import { persist } from 'zustand/middleware'

// Exemplo de store - ajuste conforme necessário
interface AuthState {
  user: unknown | null
  token: string | null
  isAuthenticated: boolean
  login: (user: unknown, token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) => {
        set({ user, token, isAuthenticated: true })
      },
      logout: () => {
        set({ user: null, token: null, isAuthenticated: false })
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
