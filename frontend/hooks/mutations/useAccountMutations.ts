import { useMutation, useQueryClient } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { CreateAccountRequest, UpdateAccountRequest } from '@/types/financial'

export function useCreateAccount() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateAccountRequest) => financialService.createAccount(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    },
  })
}

export function useUpdateAccount() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateAccountRequest }) =>
      financialService.updateAccount(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['accounts', variables.id] })
    },
  })
}

export function useDeleteAccount() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => financialService.deleteAccount(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    },
  })
}

export function useUpdateInitialBalance() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: { balance: number; date: string } }) =>
      financialService.updateInitialBalance(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['accounts', variables.id] })
    },
  })
}

export function useRecalculateBalance() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => financialService.recalculateBalance(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['accounts', id] })
    },
  })
}
