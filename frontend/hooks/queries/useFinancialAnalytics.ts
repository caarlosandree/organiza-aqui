import { useQuery } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type {
  PatrimonioLiquidoResponse,
  CalendarioVencimentosResponse,
  GastosPorTagResponse,
} from '@/types/financial'

export function usePatrimonioLiquido() {
  return useQuery<PatrimonioLiquidoResponse>({
    queryKey: ['analytics', 'patrimonio-liquido'],
    queryFn: () => financialService.getPatrimonioLiquido(),
    staleTime: 5 * 60 * 1000, // 5 minutos (cálculo pesado)
  })
}

export function useCalendarioVencimentos(days?: number) {
  return useQuery<CalendarioVencimentosResponse>({
    queryKey: ['analytics', 'calendario-vencimentos', days],
    queryFn: () => financialService.getCalendarioVencimentos(days),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export function useGastosPorTag(startDate: string, endDate: string) {
  return useQuery<GastosPorTagResponse>({
    queryKey: ['analytics', 'gastos-por-tag', startDate, endDate],
    queryFn: () => financialService.getGastosPorTag(startDate, endDate),
    enabled: !!startDate && !!endDate,
    staleTime: 5 * 60 * 1000, // 5 minutos (cálculo pesado)
  })
}
