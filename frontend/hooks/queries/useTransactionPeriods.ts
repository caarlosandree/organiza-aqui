import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { TransactionPeriod, TransactionPeriodFilters, PeriodWithTransactions } from '@/types/financial'

export function useTransactionPeriods(filters?: TransactionPeriodFilters) {
  return useQuery<TransactionPeriod[]>({
    queryKey: ['transaction-periods', filters],
    queryFn: () => financialService.getTransactionPeriods(filters),
    staleTime: 30 * 1000, // 30 segundos
  })
}

export function useTransactionPeriod(id: string) {
  return useQuery<PeriodWithTransactions>({
    queryKey: ['transaction-periods', id],
    queryFn: () => financialService.getTransactionPeriod(id),
    enabled: !!id,
    staleTime: 30 * 1000, // 30 segundos
  })
}
