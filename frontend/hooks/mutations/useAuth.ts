import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { authService, type LoginRequest, type RegisterRequest } from '@/services/authService'
import { useAuthStore } from '@/stores/authStore'

export function useLogin() {
  const router = useRouter()
  const { login } = useAuthStore()

  return useMutation({
    mutationFn: (data: LoginRequest) => authService.login(data),
    onSuccess: (response) => {
      login(response.user, response.token)
      // Aguardar um pouco mais para garantir que o estado foi atualizado e persistido
      setTimeout(() => {
        router.replace('/tasks')
      }, 200)
    },
    onError: (error: unknown) => {
      console.error('Erro ao fazer login:', error)
    },
  })
}

export function useRegister() {
  const router = useRouter()
  const { login } = useAuthStore()

  return useMutation({
    mutationFn: (data: RegisterRequest) => authService.register(data),
    onSuccess: (response) => {
      login(response.user, response.token)
      // Aguardar um pouco mais para garantir que o estado foi atualizado e persistido
      setTimeout(() => {
        router.replace('/tasks')
      }, 200)
    },
    onError: (error: unknown) => {
      console.error('Erro ao registrar:', error)
    },
  })
}

export function useLogout() {
  const router = useRouter()
  const { logout } = useAuthStore()

  return useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: () => {
      logout()
      router.push('/login')
    },
    onError: (error: unknown) => {
      console.error('Erro ao fazer logout:', error)
      // Mesmo com erro, fazer logout local
      logout()
      router.push('/login')
    },
  })
}
