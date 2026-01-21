import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { Category } from '@/types/financial'

export function useCategories(params?: { type?: string; tree?: boolean }) {
  return useQuery<Category[]>({
    queryKey: ['categories', params],
    queryFn: () => financialService.getCategories(params),
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export function useCategory(id: string) {
  return useQuery<Category>({
    queryKey: ['categories', id],
    queryFn: () => financialService.getCategory(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}
