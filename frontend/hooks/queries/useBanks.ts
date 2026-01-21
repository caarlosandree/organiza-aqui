import { useQuery } from '@tanstack/react-query'
import { bankService } from '@/services/bankService'
import type { Bank } from '@/types/bank'

export const useBanks = () => {
  return useQuery<Bank[]>({
    queryKey: ['banks'],
    queryFn: bankService.getBanks,
    staleTime: 24 * 60 * 60 * 1000, // 24 horas - bancos não mudam frequentemente
  })
}

export const useBank = (bankId: string) => {
  return useQuery<Bank>({
    queryKey: ['banks', bankId],
    queryFn: () => bankService.getBank(bankId),
    enabled: !!bankId,
    staleTime: 24 * 60 * 60 * 1000,
  })
}
