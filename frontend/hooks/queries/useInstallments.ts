import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { Transaction } from '@/types/financial'

export function useInstallments(parentId: string) {
  return useQuery<Transaction[]>({
    queryKey: ['installments', parentId],
    queryFn: () => financialService.getInstallments(parentId),
    enabled: !!parentId,
    staleTime: 2 * 60 * 1000, // 2 minutos
  })
}
