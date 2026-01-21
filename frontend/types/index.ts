// Tipos globais da aplicação
// Adicione seus tipos aqui conforme necessário

export interface ApiResponse<T> {
  data: T
  message?: string
  error?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export * from './financial'
export * from './task'
export * from './timeline'
export * from './knowledge'