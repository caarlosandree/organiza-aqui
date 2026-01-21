import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  email: string
  username: string
  name: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  setUser: (user: User) => void
  setToken: (token: string) => void
  login: (user: User, token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setUser: (user) => set({ user, isAuthenticated: !!get().token }),
      setToken: (token) => {
        set({ token, isAuthenticated: !!token })
        if (typeof window !== 'undefined') {
          localStorage.setItem('token', token)
        }
      },
      login: (user, token) => {
        set({ user, token, isAuthenticated: true })
        if (typeof window !== 'undefined') {
          localStorage.setItem('token', token)
        }
      },
      logout: () => {
        set({ user: null, token: null, isAuthenticated: false })
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token')
        }
      },
    }),
    {
      name: 'auth-storage',
      // Garantir que isAuthenticated seja calculado após restaurar
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.isAuthenticated = !!(state.token && state.user)
        }
      },
    }
  )
)
