import { useMutation, useQueryClient } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'

export function useImportOFX() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ accountId, file }: { accountId: string; file: File }) =>
      financialService.importOFX(accountId, file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
    },
  })
}

export function useImportCSV() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ accountId, file, delimiter }: { accountId: string; file: File; delimiter?: string }) =>
      financialService.importCSV(accountId, file, delimiter),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
    },
  })
}

export function usePreviewOFX() {
  return useMutation({
    mutationFn: ({ accountId, file }: { accountId: string; file: File }) =>
      financialService.previewOFX(accountId, file),
  })
}

export function usePreviewCSV() {
  return useMutation({
    mutationFn: ({ accountId, file, delimiter }: { accountId: string; file: File; delimiter?: string }) =>
      financialService.previewCSV(accountId, file, delimiter),
  })
}
