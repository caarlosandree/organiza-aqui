import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { CreditCard, CreditCardBill, ProjecaoFaturaResponse } from '@/types/financial'

export function useCreditCards() {
  return useQuery<CreditCard[]>({
    queryKey: ['credit-cards'],
    queryFn: () => financialService.getCreditCards(),
    staleTime: 2 * 60 * 1000, // 2 minutos
  })
}

export function useCreditCard(id: string) {
  return useQuery<CreditCard>({
    queryKey: ['credit-cards', id],
    queryFn: () => financialService.getCreditCard(id),
    enabled: !!id,
    staleTime: 2 * 60 * 1000,
  })
}

export function useCreditCardBills(creditCardId: string) {
  return useQuery<CreditCardBill[]>({
    queryKey: ['credit-cards', creditCardId, 'bills'],
    queryFn: () => financialService.getBills(creditCardId),
    enabled: !!creditCardId,
    staleTime: 1 * 60 * 1000, // 1 minuto (faturas podem mudar mais frequentemente)
  })
}

export function useAvailableLimit(creditCardId: string) {
  return useQuery<{ available_limit: number }>({
    queryKey: ['credit-cards', creditCardId, 'available-limit'],
    queryFn: () => financialService.getAvailableLimit(creditCardId),
    enabled: !!creditCardId,
    staleTime: 30 * 1000, // 30 segundos (limite muda com transações)
  })
}

export function useBillProjection(creditCardId: string) {
  return useQuery<ProjecaoFaturaResponse>({
    queryKey: ['credit-cards', creditCardId, 'bill-projection'],
    queryFn: () => financialService.getBillProjection(creditCardId),
    enabled: !!creditCardId,
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}
