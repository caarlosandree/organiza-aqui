import { z } from 'zod'

// Calendar Event Schemas
export const createCalendarEventSchema = z.object({
  title: z.string().min(1, 'Título é obrigatório').max(255, 'Título muito longo'),
  description: z.string().max(5000, 'Descrição muito longa').optional(),
  start_date: z.string().min(1, 'Data de início é obrigatória'),
  end_date: z.string().optional(),
  all_day: z.boolean().default(false),
  location: z.string().max(255, 'Localização muito longa').optional(),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor deve ser um código hexadecimal').default('#3b82f6'),
})

export const updateCalendarEventSchema = createCalendarEventSchema

export type CreateCalendarEventFormData = z.infer<typeof createCalendarEventSchema>
export type UpdateCalendarEventFormData = z.infer<typeof updateCalendarEventSchema>

// Note Schemas
export const createNoteSchema = z.object({
  title: z.string().min(1, 'Título é obrigatório').max(255, 'Título muito longo'),
  content: z.string().min(1, 'Conteúdo é obrigatório'),
  tags: z.array(z.string()).optional().default([]),
  is_pinned: z.boolean().optional().default(false),
})

export const updateNoteSchema = createNoteSchema

export type CreateNoteFormData = z.infer<typeof createNoteSchema>
export type UpdateNoteFormData = z.infer<typeof updateNoteSchema>

// Habit Schemas
export const createHabitSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome muito longo'),
  description: z.string().max(5000, 'Descrição muito longa').optional(),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor deve ser um código hexadecimal').default('#3b82f6'),
  frequency: z.enum(['daily', 'weekly', 'monthly'], {
    message: 'Frequência inválida',
  }),
  target_days: z.number().int().min(1, 'Deve ser pelo menos 1').max(365, 'Máximo 365 dias'),
})

export const updateHabitSchema = createHabitSchema

export type CreateHabitFormData = z.infer<typeof createHabitSchema>
export type UpdateHabitFormData = z.infer<typeof updateHabitSchema>

// Habit Tracking Schemas
export const createHabitTrackingSchema = z.object({
  habit_id: z.string().uuid('ID de hábito inválido'),
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
  completed: z.boolean().optional().default(true),
  notes: z.string().max(1000, 'Notas muito longas').optional(),
})

export const updateHabitTrackingSchema = z.object({
  completed: z.boolean(),
  notes: z.string().max(1000, 'Notas muito longas').optional(),
})

export type CreateHabitTrackingFormData = z.infer<typeof createHabitTrackingSchema>
export type UpdateHabitTrackingFormData = z.infer<typeof updateHabitTrackingSchema>
