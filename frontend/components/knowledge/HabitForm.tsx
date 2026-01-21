'use client'

import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  createHabitSchema,
  updateHabitSchema,
  type CreateHabitFormData,
  type UpdateHabitFormData,
} from '@/schemas/knowledgeSchema'
import type { Habit } from '@/types/knowledge'

interface HabitFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  habit?: Habit
  onSubmit: (data: CreateHabitFormData | UpdateHabitFormData) => void
  isSubmitting?: boolean
}

export const HabitForm = ({
  open,
  onOpenChange,
  habit,
  onSubmit,
  isSubmitting = false,
}: HabitFormProps) => {
  const isEditing = !!habit

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
  } = useForm({
    resolver: zodResolver(isEditing ? updateHabitSchema : createHabitSchema),
    defaultValues: habit
      ? {
          name: habit.name,
          description: habit.description || '',
          color: habit.color || '#3b82f6',
          frequency: habit.frequency,
          target_days: habit.target_days,
        }
      : {
          color: '#3b82f6',
          frequency: 'daily',
          target_days: 1,
        },
  })

  const handleFormSubmit = (data: CreateHabitFormData | UpdateHabitFormData) => {
    onSubmit(data)
    if (!isEditing) {
      reset()
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{isEditing ? 'Editar Hábito' : 'Novo Hábito'}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as informações do hábito'
              : 'Preencha os dados para criar um novo hábito'}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Nome</Label>
            <Input
              id="name"
              {...register('name')}
              className={errors.name ? 'border-destructive' : ''}
              placeholder="Digite o nome do hábito"
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Descrição</Label>
            <Textarea
              id="description"
              {...register('description')}
              className={errors.description ? 'border-destructive' : ''}
              placeholder="Digite a descrição do hábito (opcional)"
              rows={3}
            />
            {errors.description && (
              <p className="text-sm text-destructive">
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="frequency">Frequência</Label>
              <Controller
                name="frequency"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger
                      id="frequency"
                      className={errors.frequency ? 'border-destructive' : ''}
                    >
                      <SelectValue placeholder="Selecione a frequência" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="daily">Diário</SelectItem>
                      <SelectItem value="weekly">Semanal</SelectItem>
                      <SelectItem value="monthly">Mensal</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.frequency && (
                <p className="text-sm text-destructive">
                  {errors.frequency.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="target_days">Dias por Período</Label>
              <Input
                id="target_days"
                type="number"
                min="1"
                max="365"
                {...register('target_days', { valueAsNumber: true })}
                className={errors.target_days ? 'border-destructive' : ''}
              />
              {errors.target_days && (
                <p className="text-sm text-destructive">
                  {errors.target_days.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="color">Cor</Label>
            <div className="flex items-center gap-2">
              <Input
                id="color"
                type="color"
                {...register('color')}
                className={errors.color ? 'border-destructive' : 'w-20 h-10'}
              />
              <Input
                type="text"
                {...register('color')}
                className={errors.color ? 'border-destructive' : ''}
                placeholder="#3b82f6"
              />
            </div>
            {errors.color && (
              <p className="text-sm text-destructive">{errors.color.message}</p>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Salvando...' : isEditing ? 'Atualizar' : 'Criar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
