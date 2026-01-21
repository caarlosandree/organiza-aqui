import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type {
  PatrimonioLiquidoResponse,
  CalendarioVencimentosResponse,
  GastosPorTagResponse,
} from '@/types/financial'

export function useNetWorth() {
  return useQuery<PatrimonioLiquidoResponse>({
    queryKey: ['analytics', 'net-worth'],
    queryFn: () => financialService.getPatrimonioLiquido(),
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export function useUpcomingBills(days: number = 30) {
  return useQuery<CalendarioVencimentosResponse>({
    queryKey: ['analytics', 'upcoming-bills', days],
    queryFn: () => financialService.getCalendarioVencimentos(days),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export function useSpendingByTag(startDate: string, endDate: string) {
  return useQuery<GastosPorTagResponse>({
    queryKey: ['analytics', 'spending-by-tag', startDate, endDate],
    queryFn: () => financialService.getGastosPorTag(startDate, endDate),
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}
