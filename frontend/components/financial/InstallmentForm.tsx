'use client'

import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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
import { z } from 'zod'
import { transactionSchema, type TransactionFormData } from '@/schemas/financialSchema'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useCategories } from '@/hooks/queries/useCategories'

interface InstallmentFormProps {
  onSubmit: (data: TransactionFormData & { total_installments: number }) => void
  onCancel?: () => void
  isLoading?: boolean
}

export function InstallmentForm({ onSubmit, onCancel, isLoading }: InstallmentFormProps) {
  const { data: accounts } = useAccounts()

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<TransactionFormData & { total_installments: number }>({
    resolver: zodResolver(transactionSchema.safeExtend({
      total_installments: z.number().int().min(2).max(999),
    })),
    defaultValues: {
      date: new Date().toISOString().split('T')[0],
      total_installments: 2,
      type: 'expense',
    },
  })

  const transactionType = watch('type')
  const { data: filteredCategories } = useCategories({
    type: transactionType,
  })

  const handleFormSubmit = (data: TransactionFormData & { total_installments: number }) => {
    const amountInCents = Math.round(data.amount * 100)
    onSubmit({
      ...data,
      amount: amountInCents,
    })
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="account_id">Conta</Label>
        <Controller
          name="account_id"
          control={control}
          render={({ field }) => (
            <Select value={field.value} onValueChange={field.onChange}>
              <SelectTrigger id="account_id" className={errors.account_id ? 'border-destructive' : ''}>
                <SelectValue placeholder="Selecione a conta" />
              </SelectTrigger>
              <SelectContent>
                {accounts?.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.account_id && (
          <p className="text-sm text-destructive">{errors.account_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="type">Tipo</Label>
        <Controller
          name="type"
          control={control}
          render={({ field }) => (
            <Select
              value={field.value}
              onValueChange={(value) => {
                field.onChange(value)
                setValue('category_id', undefined)
              }}
            >
              <SelectTrigger id="type" className={errors.type ? 'border-destructive' : ''}>
                <SelectValue placeholder="Selecione o tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="income">Receita</SelectItem>
                <SelectItem value="expense">Despesa</SelectItem>
              </SelectContent>
            </Select>
          )}
        />
        {errors.type && (
          <p className="text-sm text-destructive">{errors.type.message}</p>
        )}
      </div>

      {(transactionType === 'income' || transactionType === 'expense') && (
        <div className="space-y-2">
          <Label htmlFor="category_id">Categoria (Opcional)</Label>
          <Controller
            name="category_id"
            control={control}
            render={({ field }) => (
              <Select
                value={field.value || ''}
                onValueChange={(value) => field.onChange(value || undefined)}
              >
                <SelectTrigger
                  id="category_id"
                  className={errors.category_id ? 'border-destructive' : ''}
                >
                  <SelectValue placeholder="Selecione a categoria" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Nenhuma</SelectItem>
                  {filteredCategories?.map((category) => (
                    <SelectItem key={category.id} value={category.id}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {errors.category_id && (
            <p className="text-sm text-destructive">{errors.category_id.message}</p>
          )}
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor="amount">Valor Total</Label>
        <Input
          id="amount"
          type="number"
          step="0.01"
          {...register('amount', { valueAsNumber: true })}
          className={errors.amount ? 'border-destructive' : ''}
          placeholder="0.00"
        />
        {errors.amount && (
          <p className="text-sm text-destructive">{errors.amount.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="total_installments">Número de Parcelas</Label>
        <Input
          id="total_installments"
          type="number"
          min="2"
          max="999"
          {...register('total_installments', { valueAsNumber: true })}
          className={errors.total_installments ? 'border-destructive' : ''}
        />
        {errors.total_installments && (
          <p className="text-sm text-destructive">{errors.total_installments.message}</p>
        )}
        <p className="text-xs text-muted-foreground">
          O valor será dividido igualmente entre as parcelas
        </p>
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Descrição</Label>
        <Textarea
          id="description"
          {...register('description')}
          className={errors.description ? 'border-destructive' : ''}
          placeholder="Descrição da transação"
          rows={3}
        />
        {errors.description && (
          <p className="text-sm text-destructive">{errors.description.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="date">Data da Primeira Parcela</Label>
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

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? 'Criando...' : 'Criar Parcelas'}
        </Button>
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
        )}
      </div>
    </form>
  )
}
