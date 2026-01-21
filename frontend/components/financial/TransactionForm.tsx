'use client'

import { useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { X, Landmark, CreditCard } from 'lucide-react'
import { transactionSchema, type TransactionFormData } from '@/schemas/financialSchema'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useCreditCards } from '@/hooks/queries/useCreditCards'
import { useCategories } from '@/hooks/queries/useCategories'
import { cn } from '@/lib/utils'
import type { Transaction } from '@/types/financial'

interface TransactionFormProps {
  transaction?: Transaction
  onSubmit: (data: TransactionFormData) => void
  onCancel?: () => void
  isLoading?: boolean
}

export function TransactionForm({
  transaction,
  onSubmit,
  onCancel,
  isLoading,
}: TransactionFormProps) {
  const { data: accounts } = useAccounts()
  const { data: creditCards } = useCreditCards()
  const [tagInput, setTagInput] = useState('')
  const [tags, setTags] = useState<string[]>(transaction?.tags || [])
  
  // Determinar contexto inicial baseado na transação existente ou default
  const initialContext = transaction 
    ? (accounts?.find(a => a.id === transaction.account_id)?.type === 'credit' ? 'credit_card' : 'bank')
    : 'bank'

  const [context, setContext] = useState<'bank' | 'credit_card'>(initialContext)

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<TransactionFormData>({
    resolver: zodResolver(transactionSchema),
    defaultValues: transaction
      ? {
          account_id: transaction.account_id,
          category_id: transaction.category_id || undefined,
          type: transaction.type,
          amount: transaction.amount / 100, // Converter centavos para reais
          description: transaction.description || '',
          date: transaction.date,
          reference_month: transaction.reference_month || undefined,
          status: transaction.status || 'pending',
          tags: transaction.tags || [],
          to_account_id: transaction.to_account_id || undefined,
          total_installments: transaction.total_installments || undefined,
        }
      : {
          date: new Date().toISOString().split('T')[0],
          reference_month: new Date().toISOString().slice(0, 7), // YYYY-MM
          status: 'pending',
          tags: [],
          type: 'expense',
        },
  })

  const transactionType = watch('type')
  const watchAccountId = watch('account_id')
  const watchTotalInstallments = watch('total_installments')
  const watchDate = watch('date')

  // Quando o contexto muda, ajustar o tipo e limpar conta selecionada
  const handleContextChange = (newContext: 'bank' | 'credit_card') => {
    setContext(newContext)
    if (newContext === 'credit_card') {
      setValue('type', 'expense')
      setValue('account_id', '')
    }
  }

  // Filtrar contas baseado no contexto
  const availableAccounts = context === 'credit_card'
    ? accounts?.filter(acc => creditCards?.some(cc => cc.account_id === acc.id))
    : accounts?.filter(acc => acc.type !== 'credit')

  // Atualizar reference_month quando date mudar (se não foi definido manualmente)
  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newDate = e.target.value
    register('date').onChange(e)
    // Se reference_month não foi definido manualmente, atualizar baseado na data
    const currentRefMonth = watch('reference_month')
    if (!currentRefMonth || currentRefMonth === transaction?.reference_month) {
      if (newDate) {
        setValue('reference_month', newDate.slice(0, 7)) // YYYY-MM
      }
    }
  }

  // Atualizar categorias quando o tipo mudar
  const { data: filteredCategories } = useCategories({
    type: transactionType === 'transfer' ? undefined : transactionType,
  })

  const handleAddTag = () => {
    const trimmedTag = tagInput.trim()
    if (trimmedTag && !tags.includes(trimmedTag)) {
      const newTags = [...tags, trimmedTag]
      setTags(newTags)
      setValue('tags', newTags)
      setTagInput('')
    }
  }

  const handleRemoveTag = (tagToRemove: string) => {
    const newTags = tags.filter((tag) => tag !== tagToRemove)
    setTags(newTags)
    setValue('tags', newTags)
  }

  const handleFormSubmit = (data: TransactionFormData) => {
    // Converter valor para centavos (amount já vem como number em reais)
    const amountInCents = Math.round(data.amount * 100)
    onSubmit({
      ...data,
      amount: amountInCents,
      tags: tags.length > 0 ? tags : undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {/* Seleção de Contexto: Banco ou Cartão */}
      {!transaction && (
        <div className="grid grid-cols-2 gap-3 p-1 bg-muted rounded-lg">
          <button
            type="button"
            className={cn(
              'py-2 px-3 rounded-md text-sm font-medium transition-all flex items-center justify-center gap-2',
              context === 'bank'
                ? 'bg-background shadow text-foreground'
                : 'text-muted-foreground hover:text-foreground'
            )}
            onClick={() => handleContextChange('bank')}
          >
            <Landmark size={16} /> Conta Banco
          </button>
          <button
            type="button"
            className={cn(
              'py-2 px-3 rounded-md text-sm font-medium transition-all flex items-center justify-center gap-2',
              context === 'credit_card'
                ? 'bg-background shadow text-foreground'
                : 'text-muted-foreground hover:text-foreground'
            )}
            onClick={() => handleContextChange('credit_card')}
          >
            <CreditCard size={16} /> Cartão Crédito
          </button>
        </div>
      )}

      {/* Seleção de Conta/Cartão */}
      <div className="space-y-2">
        <Label htmlFor="account_id">
          {context === 'credit_card' ? 'Cartão vinculado à conta' : 'Conta Bancária'}
        </Label>
        <Controller
          name="account_id"
          control={control}
          render={({ field }) => (
            <Select value={field.value || ''} onValueChange={field.onChange}>
              <SelectTrigger id="account_id" className={errors.account_id ? 'border-destructive' : ''}>
                <SelectValue placeholder="Selecione a conta" />
              </SelectTrigger>
              <SelectContent>
                {availableAccounts?.map((account) => (
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

      {/* Seleção de Receita/Despesa (apenas para Banco) */}
      {context === 'bank' && !transaction && (
        <div className="flex gap-4 p-2 border rounded-lg">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="type"
              className="accent-destructive w-4 h-4"
              checked={transactionType === 'expense'}
              onChange={() => {
                setValue('type', 'expense')
                setValue('category_id', undefined)
              }}
            />
            <span className="text-sm text-foreground font-medium">Saída</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="type"
              className="accent-emerald-500 w-4 h-4"
              checked={transactionType === 'income'}
              onChange={() => {
                setValue('type', 'income')
                setValue('category_id', undefined)
              }}
            />
            <span className="text-sm text-foreground font-medium">Entrada</span>
          </label>
        </div>
      )}

      {/* Tipo (para edição ou transferências) */}
      {transaction && (
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
                  <SelectItem value="transfer">Transferência</SelectItem>
                </SelectContent>
              </Select>
            )}
          />
          {errors.type && (
            <p className="text-sm text-destructive">{errors.type.message}</p>
          )}
        </div>
      )}

      {transactionType === 'transfer' && (
        <div className="space-y-2">
          <Label htmlFor="to_account_id">Conta Destino</Label>
          <Controller
            name="to_account_id"
            control={control}
            render={({ field }) => (
              <Select value={field.value || ''} onValueChange={field.onChange}>
                <SelectTrigger
                  id="to_account_id"
                  className={errors.to_account_id ? 'border-destructive' : ''}
                >
                  <SelectValue placeholder="Selecione a conta destino" />
                </SelectTrigger>
                <SelectContent>
                  {accounts?.filter((acc) => acc.id !== watch('account_id')).map((account) => (
                    <SelectItem key={account.id} value={account.id}>
                      {account.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {errors.to_account_id && (
            <p className="text-sm text-destructive">{errors.to_account_id.message}</p>
          )}
        </div>
      )}

      {(transactionType === 'income' || transactionType === 'expense') && context === 'bank' && (
        <div className="space-y-2">
          <Label htmlFor="total_installments">Parcelamento (Opcional)</Label>
          <Input
            id="total_installments"
            type="number"
            min="1"
            max="999"
            {...register('total_installments', { valueAsNumber: true })}
            className={errors.total_installments ? 'border-destructive' : ''}
            placeholder="Número de parcelas (ex: 12)"
          />
          {errors.total_installments && (
            <p className="text-sm text-destructive">{errors.total_installments.message}</p>
          )}
          {watchTotalInstallments && watchTotalInstallments > 1 && (
            <p className="text-xs text-muted-foreground">
              Esta transação será dividida em {watchTotalInstallments} parcelas mensais
            </p>
          )}
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor="description">Descrição</Label>
        <Input
          id="description"
          type="text"
          {...register('description')}
          className={errors.description ? 'border-destructive' : ''}
          placeholder="Ex: Supermercado..."
        />
        {errors.description && (
          <p className="text-sm text-destructive">{errors.description.message}</p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="amount">Valor</Label>
          <Input
            id="amount"
            type="number"
            step="0.01"
            {...register('amount', { valueAsNumber: true })}
            className={errors.amount ? 'border-destructive' : ''}
            placeholder="0,00"
          />
          {errors.amount && (
            <p className="text-sm text-destructive">{errors.amount.message}</p>
          )}
        </div>
        {(transactionType === 'income' || transactionType === 'expense') && (
          <div className="space-y-2">
            <Label htmlFor="category_id_secondary">Categoria</Label>
            <Controller
              name="category_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value || 'none'}
                  onValueChange={(value) => field.onChange(value === 'none' ? undefined : value)}
                >
                  <SelectTrigger
                    id="category_id_secondary"
                    className={errors.category_id ? 'border-destructive' : ''}
                  >
                    <SelectValue placeholder="Selecione" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">Outros</SelectItem>
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
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="date">Data Real</Label>
          <Input
            id="date"
            type="date"
            {...register('date')}
            onChange={handleDateChange}
            className={errors.date ? 'border-destructive' : ''}
          />
          {errors.date && (
            <p className="text-sm text-destructive">{errors.date.message}</p>
          )}
        </div>
        <div className="space-y-2">
          <Label htmlFor="reference_month" className="font-semibold text-primary">
            Mês Referência
          </Label>
          <Input
            id="reference_month"
            type="month"
            {...register('reference_month')}
            className={cn(
              'border-2 bg-primary/5 border-primary/20 focus:border-primary',
              errors.reference_month ? 'border-destructive' : ''
            )}
          />
          {errors.reference_month && (
            <p className="text-sm text-destructive">{errors.reference_month.message}</p>
          )}
        </div>
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
                <SelectItem value="paid">Paga/Concluída</SelectItem>
                <SelectItem value="cancelled">Cancelada</SelectItem>
              </SelectContent>
            </Select>
          )}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="tags">Tags</Label>
        <div className="flex gap-2">
          <Input
            id="tags"
            type="text"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault()
                handleAddTag()
              }
            }}
            placeholder="Digite uma tag e pressione Enter"
          />
          <Button type="button" variant="outline" onClick={handleAddTag}>
            Adicionar
          </Button>
        </div>
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {tags.map((tag) => (
              <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                {tag}
                <button
                  type="button"
                  onClick={() => handleRemoveTag(tag)}
                  className="ml-1 hover:bg-destructive/20 rounded-full p-0.5"
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
            ))}
          </div>
        )}
      </div>

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? 'Salvando...' : transaction ? 'Atualizar' : 'Criar'}
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
