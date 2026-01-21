import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { Account } from '@/types/financial'

export function useAccounts() {
  return useQuery<Account[]>({
    queryKey: ['accounts'],
    queryFn: financialService.getAccounts,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export function useAccount(id: string) {
  return useQuery<Account>({
    queryKey: ['accounts', id],
    queryFn: () => financialService.getAccount(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}
