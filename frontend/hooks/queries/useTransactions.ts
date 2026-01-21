import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { TransactionFilters, TransactionsResponse } from '@/types/financial'

export function useTransactions(filters?: TransactionFilters) {
  return useQuery<TransactionsResponse>({
    queryKey: ['transactions', filters],
    queryFn: () => financialService.getTransactions(filters),
    staleTime: 2 * 60 * 1000, // 2 minutos
  })
}

export function useTransaction(id: string) {
  return useQuery({
    queryKey: ['transactions', id],
    queryFn: () => financialService.getTransaction(id),
    enabled: !!id,
    staleTime: 2 * 60 * 1000,
  })
}
