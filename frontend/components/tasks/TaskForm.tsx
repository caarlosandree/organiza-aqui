'use client'

import { useForm } from 'react-hook-form'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Controller } from 'react-hook-form'
import {
  createTaskSchema,
  updateTaskSchema,
  type CreateTaskFormData,
  type UpdateTaskFormData,
} from '@/schemas/taskSchema'
import type { Task, TaskStatus } from '@/types/task'
import { useTaskStatuses } from '@/hooks/queries/useTaskStatuses'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useCategories } from '@/hooks/queries/useCategories'

interface TaskFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  task?: Task
  onSubmit: (data: CreateTaskFormData | UpdateTaskFormData) => void
  isSubmitting?: boolean
}

export const TaskForm = ({
  open,
  onOpenChange,
  task,
  onSubmit,
  isSubmitting = false,
}: TaskFormProps) => {
  const { data: statuses = [] } = useTaskStatuses()
  const { data: accounts = [] } = useAccounts()
  const { data: categories = [] } = useCategories()
  const isEditing = !!task

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
  } = useForm<CreateTaskFormData | UpdateTaskFormData>({
    resolver: zodResolver(isEditing ? updateTaskSchema : createTaskSchema),
    defaultValues: task
      ? {
          status_id: task.status_id,
          title: task.title,
          description: task.description || '',
          priority: task.priority,
          due_date: task.due_date || '',
          financial_account_id: task.financial_account_id || '',
          financial_amount: task.financial_amount ? task.financial_amount / 100 : undefined,
          financial_category_id: task.financial_category_id || '',
        }
      : {
          priority: 'medium',
        },
  })

  const handleFormSubmit = (data: CreateTaskFormData | UpdateTaskFormData) => {
    // Converter valor de reais para centavos
    const submitData = { ...data }
    if (submitData.financial_amount !== undefined) {
      submitData.financial_amount = Math.round(submitData.financial_amount * 100)
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
          <DialogTitle>{isEditing ? 'Editar Tarefa' : 'Nova Tarefa'}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as informações da tarefa'
              : 'Preencha os dados para criar uma nova tarefa'}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="status_id">Status</Label>
            <Controller
              name="status_id"
              control={control}
              render={({ field }) => (
                <Select value={field.value} onValueChange={field.onChange}>
                  <SelectTrigger id="status_id" className={errors.status_id ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Selecione um status" />
                  </SelectTrigger>
                  <SelectContent>
                    {statuses.map((status: TaskStatus) => (
                      <SelectItem key={status.id} value={status.id}>
                        <div className="flex items-center gap-2">
                          <div
                            className="h-3 w-3 rounded-full"
                            style={{ backgroundColor: status.color }}
                          />
                          {status.name}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {errors.status_id && (
              <p className="text-sm text-destructive">{errors.status_id.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="title">Título</Label>
            <Input
              id="title"
              {...register('title')}
              className={errors.title ? 'border-destructive' : ''}
              placeholder="Digite o título da tarefa"
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
              placeholder="Digite a descrição da tarefa (opcional)"
              rows={4}
            />
            {errors.description && (
              <p className="text-sm text-destructive">{errors.description.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="priority">Prioridade</Label>
            <Controller
              name="priority"
              control={control}
              render={({ field }) => (
                <Select value={field.value} onValueChange={field.onChange}>
                  <SelectTrigger id="priority" className={errors.priority ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Selecione a prioridade" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="low">Baixa</SelectItem>
                    <SelectItem value="medium">Média</SelectItem>
                    <SelectItem value="high">Alta</SelectItem>
                    <SelectItem value="urgent">Urgente</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
            {errors.priority && (
              <p className="text-sm text-destructive">{errors.priority.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="due_date">Data de Vencimento</Label>
            <Input
              id="due_date"
              type="date"
              {...register('due_date')}
              className={errors.due_date ? 'border-destructive' : ''}
            />
            {errors.due_date && (
              <p className="text-sm text-destructive">{errors.due_date.message}</p>
            )}
          </div>

          <div className="border-t pt-4 space-y-4">
            <p className="text-sm font-medium">Integração Financeira (Opcional)</p>
            <p className="text-xs text-muted-foreground">
              Ao completar esta tarefa, uma transação será criada automaticamente
            </p>

            <div className="space-y-2">
              <Label htmlFor="financial_account_id">Conta</Label>
              <Controller
                name="financial_account_id"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value || ''}
                    onValueChange={(value) => field.onChange(value || undefined)}
                  >
                    <SelectTrigger
                      id="financial_account_id"
                      className={errors.financial_account_id ? 'border-destructive' : ''}
                    >
                      <SelectValue placeholder="Selecione uma conta (opcional)" />
                    </SelectTrigger>
                    <SelectContent>
                      {accounts.map((account) => (
                        <SelectItem key={account.id} value={account.id}>
                          {account.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.financial_account_id && (
                <p className="text-sm text-destructive">
                  {errors.financial_account_id.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="financial_amount">Valor (R$)</Label>
              <Input
                id="financial_amount"
                type="number"
                step="0.01"
                min="0"
                {...register('financial_amount', {
                  setValueAs: (v) => (v === '' ? undefined : parseFloat(v)),
                })}
                className={errors.financial_amount ? 'border-destructive' : ''}
                placeholder="0.00"
              />
              {errors.financial_amount && (
                <p className="text-sm text-destructive">
                  {errors.financial_amount.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="financial_category_id">Categoria</Label>
              <Controller
                name="financial_category_id"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value || ''}
                    onValueChange={(value) => field.onChange(value || undefined)}
                  >
                    <SelectTrigger
                      id="financial_category_id"
                      className={errors.financial_category_id ? 'border-destructive' : ''}
                    >
                      <SelectValue placeholder="Selecione uma categoria (opcional)" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories
                        .filter((cat) => cat.type === 'expense')
                        .map((category) => (
                          <SelectItem key={category.id} value={category.id}>
                            {category.name}
                          </SelectItem>
                        ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.financial_category_id && (
                <p className="text-sm text-destructive">
                  {errors.financial_category_id.message}
                </p>
              )}
            </div>
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
