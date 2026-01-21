'use client'

import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
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
import { updateInstallmentSchema, type UpdateInstallmentFormData } from '@/schemas/financialSchema'
import type { Transaction } from '@/types/financial'

interface InstallmentEditDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  installment: Transaction
  onSubmit: (data: UpdateInstallmentFormData) => void
  isLoading?: boolean
}

export function InstallmentEditDialog({
  open,
  onOpenChange,
  installment,
  onSubmit,
  isLoading,
}: InstallmentEditDialogProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<UpdateInstallmentFormData>({
    resolver: zodResolver(updateInstallmentSchema),
    defaultValues: {
      amount: installment.amount / 100, // converter de centavos para reais
      description: installment.description || '',
      date: installment.date,
      status: installment.status,
      scope: 'this',
    },
  })

  const handleFormSubmit = (data: UpdateInstallmentFormData) => {
    const amountInCents = Math.round(data.amount * 100)
    onSubmit({
      ...data,
      amount: amountInCents,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Editar Parcela</DialogTitle>
          <DialogDescription>
            Edite a parcela {installment.installment_number}/{installment.total_installments}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="scope">Escopo da Edição</Label>
            <Controller
              name="scope"
              control={control}
              render={({ field }) => (
                <Select value={field.value} onValueChange={field.onChange}>
                  <SelectTrigger id="scope" className={errors.scope ? 'border-destructive' : ''}>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="this">Apenas esta parcela</SelectItem>
                    <SelectItem value="this_and_future">Esta e as futuras</SelectItem>
                    <SelectItem value="all">Todas as parcelas</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
            {errors.scope && (
              <p className="text-sm text-destructive">{errors.scope.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="amount">Valor</Label>
            <Input
              id="amount"
              type="number"
              step="0.01"
              {...register('amount', { valueAsNumber: true })}
              className={errors.amount ? 'border-destructive' : ''}
            />
            {errors.amount && (
              <p className="text-sm text-destructive">{errors.amount.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Descrição</Label>
            <Textarea
              id="description"
              {...register('description')}
              className={errors.description ? 'border-destructive' : ''}
              rows={3}
            />
            {errors.description && (
              <p className="text-sm text-destructive">{errors.description.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="date">Data</Label>
            <Input
              id="date"
              type="date"
              {...register('date')}
              className={errors.date ? 'border-destructive' : ''}
            />
            {errors.date && (
              <p className="text-sm text-destructive">{errors.date.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <Controller
              name="status"
              control={control}
              render={({ field }) => (
                <Select value={field.value || 'pending'} onValueChange={field.onChange}>
                  <SelectTrigger id="status">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="pending">Pendente</SelectItem>
                    <SelectItem value="paid">Paga</SelectItem>
                    <SelectItem value="cancelled">Cancelada</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={isLoading} className="flex-1">
              {isLoading ? 'Salvando...' : 'Salvar'}
            </Button>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancelar
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
