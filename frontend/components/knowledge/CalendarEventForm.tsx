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
import { Checkbox } from '@/components/ui/checkbox'
import {
  createCalendarEventSchema,
  updateCalendarEventSchema,
  type CreateCalendarEventFormData,
  type UpdateCalendarEventFormData,
} from '@/schemas/knowledgeSchema'
import type { CalendarEvent } from '@/types/knowledge'

interface CalendarEventFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  event?: CalendarEvent
  onSubmit: (data: CreateCalendarEventFormData | UpdateCalendarEventFormData) => void
  isSubmitting?: boolean
}

export const CalendarEventForm = ({
  open,
  onOpenChange,
  event,
  onSubmit,
  isSubmitting = false,
}: CalendarEventFormProps) => {
  const isEditing = !!event

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    watch,
  } = useForm({
    resolver: zodResolver(isEditing ? updateCalendarEventSchema : createCalendarEventSchema),
    defaultValues: event
      ? {
          title: event.title,
          description: event.description || '',
          start_date: event.start_date,
          end_date: event.end_date || '',
          all_day: event.all_day ?? false,
          location: event.location || '',
          color: event.color || '#3b82f6',
        }
      : {
          all_day: false,
          color: '#3b82f6',
        },
  })

  const allDay = watch('all_day')

  const handleFormSubmit = (data: CreateCalendarEventFormData | UpdateCalendarEventFormData) => {
    // Converter datetime-local para RFC3339 se necessário
    const submitData = { ...data }
    
    const startDate = submitData.start_date as string | undefined
    const endDate = submitData.end_date as string | undefined
    
    if (startDate) {
      if (!startDate.includes('T')) {
        // Se não tem T, é date, adicionar hora
        submitData.start_date = allDay 
          ? `${startDate}T00:00:00Z`
          : `${startDate}T00:00:00`
      } else if (!startDate.includes('Z') && !startDate.includes('+')) {
        // Adicionar timezone se não tiver
        const date = new Date(startDate)
        submitData.start_date = date.toISOString()
      }
    }
    if (endDate) {
      if (!endDate.includes('T')) {
        submitData.end_date = allDay
          ? `${endDate}T23:59:59Z`
          : `${endDate}T23:59:59`
      } else if (!endDate.includes('Z') && !endDate.includes('+')) {
        const date = new Date(endDate)
        submitData.end_date = date.toISOString()
      }
    }
    onSubmit(submitData)
    if (!isEditing) {
      reset()
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Evento' : 'Novo Evento'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as informações do evento'
              : 'Preencha os dados para criar um novo evento'}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Título</Label>
            <Input
              id="title"
              {...register('title')}
              className={errors.title ? 'border-destructive' : ''}
              placeholder="Digite o título do evento"
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Descrição</Label>
            <Textarea
              id="description"
              {...register('description')}
              className={errors.description ? 'border-destructive' : ''}
              placeholder="Digite a descrição do evento (opcional)"
              rows={3}
            />
            {errors.description && (
              <p className="text-sm text-destructive">
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Controller
              name="all_day"
              control={control}
              render={({ field }) => (
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="all_day"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                  <Label htmlFor="all_day" className="cursor-pointer">
                    Dia inteiro
                  </Label>
                </div>
              )}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="start_date">
                {allDay ? 'Data' : 'Data e Hora de Início'}
              </Label>
              <Input
                id="start_date"
                type={allDay ? 'date' : 'datetime-local'}
                {...register('start_date')}
                className={errors.start_date ? 'border-destructive' : ''}
              />
              {errors.start_date && (
                <p className="text-sm text-destructive">
                  {errors.start_date.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="end_date">
                {allDay ? 'Data Final' : 'Data e Hora Final'}
              </Label>
              <Input
                id="end_date"
                type={allDay ? 'date' : 'datetime-local'}
                {...register('end_date')}
                className={errors.end_date ? 'border-destructive' : ''}
              />
              {errors.end_date && (
                <p className="text-sm text-destructive">
                  {errors.end_date.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="location">Localização</Label>
            <Input
              id="location"
              {...register('location')}
              className={errors.location ? 'border-destructive' : ''}
              placeholder="Digite a localização (opcional)"
            />
            {errors.location && (
              <p className="text-sm text-destructive">{errors.location.message}</p>
            )}
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
