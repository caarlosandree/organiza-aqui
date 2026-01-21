import { z } from 'zod'

export const createTaskStatusSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome muito longo'),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor inválida'),
  order_index: z.number().int().min(0, 'Índice deve ser positivo'),
  is_default: z.boolean().default(false),
})

export type CreateTaskStatusFormData = z.infer<typeof createTaskStatusSchema>

export const updateTaskStatusSchema = createTaskStatusSchema

export type UpdateTaskStatusFormData = z.infer<typeof updateTaskStatusSchema>

export const createTaskSchema = z.object({
  status_id: z.string().uuid('ID de status inválido'),
  title: z.string().min(1, 'Título é obrigatório').max(255, 'Título muito longo'),
  description: z.string().max(5000, 'Descrição muito longa').optional(),
  priority: z.enum(['low', 'medium', 'high', 'urgent'], {
    message: 'Prioridade inválida',
  }),
  due_date: z.string().optional(),
  financial_account_id: z.string().uuid('ID de conta inválido').optional(),
  financial_amount: z.number().int().positive('Valor deve ser positivo').optional(),
  financial_category_id: z.string().uuid('ID de categoria inválido').optional(),
})

export type CreateTaskFormData = z.infer<typeof createTaskSchema>

export const updateTaskSchema = createTaskSchema

export type UpdateTaskFormData = z.infer<typeof updateTaskSchema>
