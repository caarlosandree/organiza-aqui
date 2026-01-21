import api from '@/lib/axios'

export interface LoginRequest {
  identifier: string
  password: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
  name: string
}

export interface User {
  id: string
  email: string
  username: string
  name: string
}

export interface LoginResponse {
  token: string
  user: User
}

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data)
    return response.data
  },

  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/register', data)
    return response.data
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout')
  },

  me: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me')
    return response.data
  },
}
