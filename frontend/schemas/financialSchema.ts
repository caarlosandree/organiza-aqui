import { z } from 'zod'

export const accountSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome deve ter no máximo 255 caracteres'),
  type: z.enum(['checking', 'savings', 'credit', 'investment'], {
    message: 'Tipo inválido',
  }),
  currency: z.string().length(3, 'Moeda deve ter 3 caracteres'),
  bank_id: z.string().uuid('ID do banco inválido'),
})

export const categorySchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome deve ter no máximo 255 caracteres'),
  parent_id: z.string().uuid('ID do pai inválido').optional().nullable(),
  type: z.enum(['income', 'expense'], {
    message: 'Tipo deve ser income ou expense',
  }),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor deve ser um código hexadecimal válido'),
})

export const transactionSchema = z.object({
  account_id: z.string().uuid('ID da conta inválido'),
  category_id: z.string().uuid('ID da categoria inválido').optional().nullable(),
  type: z.enum(['income', 'expense', 'transfer'], {
    message: 'Tipo deve ser income, expense ou transfer',
  }),
  amount: z.number().positive('Valor deve ser positivo'),
  description: z.string().max(1000, 'Descrição deve ter no máximo 1000 caracteres').optional(),
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
  reference_month: z.string().regex(/^\d{4}-\d{2}$/, 'Mês de referência deve estar no formato YYYY-MM').optional(),
  status: z.enum(['pending', 'paid', 'cancelled']).optional(),
  tags: z.array(z.string()).optional(),
  to_account_id: z.string().uuid('ID da conta destino inválido').optional(),
  total_installments: z.number().int().min(1).max(999).optional(),
}).refine((data) => {
  // Se for transfer, to_account_id é obrigatório
  if (data.type === 'transfer' && !data.to_account_id) {
    return false
  }
  return true
}, {
  message: 'to_account_id é obrigatório para transferências',
  path: ['to_account_id'],
})

export const creditCardSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome deve ter no máximo 255 caracteres'),
  account_id: z.string().uuid('ID da conta inválido'),
  limit_amount: z.number().positive('Limite deve ser positivo'),
  closing_day: z.number().int().min(1).max(31, 'Dia de fechamento deve ser entre 1 e 31'),
  due_day: z.number().int().min(1).max(31, 'Dia de vencimento deve ser entre 1 e 31'),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor deve ser um código hexadecimal válido'),
})

export const payBillSchema = z.object({
  account_id: z.string().uuid('ID da conta inválido'),
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
})

export const updateInstallmentSchema = z.object({
  amount: z.number().positive('Valor deve ser positivo'),
  description: z.string().max(1000, 'Descrição deve ter no máximo 1000 caracteres'),
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
  status: z.enum(['pending', 'paid', 'cancelled']).optional(),
  scope: z.enum(['this', 'this_and_future', 'all'], {
    message: 'Escopo deve ser this, this_and_future ou all',
  }),
})

export const updateInitialBalanceSchema = z.object({
  balance: z.number().int('Saldo deve ser um número inteiro (em centavos)'),
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
})

export type AccountFormData = z.infer<typeof accountSchema>
export type CategoryFormData = z.infer<typeof categorySchema>
export type TransactionFormData = z.infer<typeof transactionSchema>
export type CreditCardFormData = z.infer<typeof creditCardSchema>
export type PayBillFormData = z.infer<typeof payBillSchema>
export type UpdateInstallmentFormData = z.infer<typeof updateInstallmentSchema>
export type UpdateInitialBalanceFormData = z.infer<typeof updateInitialBalanceSchema>