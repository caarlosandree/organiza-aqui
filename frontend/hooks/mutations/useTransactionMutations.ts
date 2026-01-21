import { useMutation, useQueryClient } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { CreateTransactionRequest, UpdateTransactionRequest } from '@/types/financial'

export function useCreateTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTransactionRequest) => financialService.createTransaction(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] }) // Atualiza limite utilizado/disponível
    },
  })
}

export function useUpdateTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTransactionRequest }) =>
      financialService.updateTransaction(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['transactions', variables.id] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] }) // Atualiza limite utilizado/disponível
    },
  })
}

export function useDeleteTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => financialService.deleteTransaction(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] }) // Atualiza limite utilizado/disponível
    },
  })
}

export function useUpdateTransactionStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: 'pending' | 'paid' | 'cancelled' }) =>
      financialService.updateTransactionStatus(id, status),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['transactions', variables.id] })
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] }) // Atualiza limite utilizado/disponível (status pode afetar limite)
    },
  })
}
