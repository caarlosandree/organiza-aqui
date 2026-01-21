import { useMutation, useQueryClient } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { CreateTransactionRequest, UpdateInstallmentRequest } from '@/types/financial'

export function useCreateInstallments() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTransactionRequest) => financialService.createInstallments(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
    },
  })
}

export function useUpdateInstallment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateInstallmentRequest }) =>
      financialService.updateInstallment(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['transactions', variables.id] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
    },
  })
}

export function useCancelFutureInstallments() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => financialService.cancelFutureInstallments(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
    },
  })
}
