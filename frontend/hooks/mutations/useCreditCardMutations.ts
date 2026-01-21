import { useMutation, useQueryClient } from '@tanstack/react-query'
import { financialService } from '@/services/financialService'
import type { CreateCreditCardRequest, UpdateCreditCardRequest, PayBillRequest } from '@/types/financial'

export function useCreateCreditCard() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateCreditCardRequest) => financialService.createCreditCard(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] })
    },
  })
}

export function useUpdateCreditCard() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCreditCardRequest }) =>
      financialService.updateCreditCard(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] })
      queryClient.invalidateQueries({ queryKey: ['credit-cards', variables.id] })
    },
  })
}

export function useDeleteCreditCard() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => financialService.deleteCreditCard(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['credit-cards'] })
    },
  })
}

export function useCloseBill() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ creditCardId, billId }: { creditCardId: string; billId: string }) =>
      financialService.closeBill(creditCardId, billId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['credit-cards', variables.creditCardId, 'bills'] })
      queryClient.invalidateQueries({ queryKey: ['credit-cards', variables.creditCardId, 'bill-projection'] })
    },
  })
}

export function usePayBill() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ creditCardId, billId, data }: { creditCardId: string; billId: string; data: PayBillRequest }) =>
      financialService.payBill(creditCardId, billId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['credit-cards', variables.creditCardId, 'bills'] })
      queryClient.invalidateQueries({ queryKey: ['credit-cards', variables.creditCardId, 'available-limit'] })
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] }) // Atualiza saldo
    },
  })
}
