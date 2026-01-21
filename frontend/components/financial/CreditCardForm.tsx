'use client'

import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { creditCardSchema, type CreditCardFormData } from '@/schemas/financialSchema'
import { useAccounts } from '@/hooks/queries/useAccounts'
import type { CreditCard } from '@/types/financial'
import { formatCurrency } from '@/utils/currency'

interface CreditCardFormProps {
  creditCard?: CreditCard
  onSubmit: (data: CreditCardFormData) => void
  onCancel?: () => void
  onDelete?: (id: string) => void
  isLoading?: boolean
}

export function CreditCardForm({ creditCard, onSubmit, onCancel, onDelete, isLoading }: CreditCardFormProps) {
  const { data: accounts } = useAccounts()
  
  // Inicializar o estado com valor formatado se houver creditCard
  const initialLimitInput = creditCard 
    ? formatCurrency(creditCard.limit_amount, 'BRL')
    : ''
  const [limitInput, setLimitInput] = useState(initialLimitInput)

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CreditCardFormData>({
    resolver: zodResolver(creditCardSchema),
    defaultValues: creditCard
      ? {
          name: creditCard.name,
          account_id: creditCard.account_id,
          limit_amount: creditCard.limit_amount / 100, // converter de centavos para reais
          closing_day: creditCard.closing_day,
          due_day: creditCard.due_day,
          color: creditCard.color,
        }
      : {
          limit_amount: 0, // Garantir valor inicial definido
          color: '#3b82f6',
          closing_day: 10,
          due_day: 15,
        },
  })

  // Inicializar o input formatado quando o formulário carregar
  useEffect(() => {
    if (creditCard) {
      setLimitInput(formatCurrency(creditCard.limit_amount, 'BRL'))
      setValue('limit_amount', creditCard.limit_amount / 100, { shouldValidate: false })
    } else {
      setLimitInput('')
      setValue('limit_amount', 0, { shouldValidate: false })
    }
  }, [creditCard, setValue])

  const handleLimitChange = (value: string) => {
    // Remove tudo que não é número (mantém apenas dígitos)
    const numbers = value.replace(/\D/g, '')

    // Se não há números, limpar
    if (!numbers) {
      setLimitInput('')
      setValue('limit_amount', 0, { shouldValidate: true })
      return
    }

    // Converter string de números para centavos
    // Ex: "054" -> 54 centavos
    const cents = parseInt(numbers, 10)

    // Formatar para exibição no padrão brasileiro
    // Dividir por 100 para obter o valor em reais
    const reais = cents / 100

    // Formatar com Intl.NumberFormat para padrão brasileiro
    const formatted = new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(reais)

    setLimitInput(formatted)

    // Atualizar o campo do formulário com valor em reais (para conversão posterior)
    setValue('limit_amount', reais, { shouldValidate: true })
  }

  return (
    <form onSubmit={handleSubmit((data) => {
      // Converter limite de reais para centavos
      onSubmit({
        ...data,
        limit_amount: Math.round(data.limit_amount * 100),
      })
    })} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Nome do Cartão</Label>
        <Input
          id="name"
          {...register('name')}
          className={errors.name ? 'border-destructive' : ''}
          placeholder="Ex: Nubank"
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="account_id">Conta Associada</Label>
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
        <Label htmlFor="limit_amount">Limite (R$)</Label>
        <Input
          id="limit_amount"
          type="text"
          value={limitInput || ''}
          onChange={(e) => handleLimitChange(e.target.value)}
          placeholder="R$ 0,00"
          className={errors.limit_amount ? 'border-destructive' : ''}
        />
        {errors.limit_amount && (
          <p className="text-sm text-destructive">{errors.limit_amount.message}</p>
        )}
        <input type="hidden" {...register('limit_amount', { valueAsNumber: true })} />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="closing_day">Dia de Fechamento</Label>
          <Input
            id="closing_day"
            type="number"
            min="1"
            max="31"
            {...register('closing_day', { valueAsNumber: true })}
            className={errors.closing_day ? 'border-destructive' : ''}
          />
          {errors.closing_day && (
            <p className="text-sm text-destructive">{errors.closing_day.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="due_day">Dia de Vencimento</Label>
          <Input
            id="due_day"
            type="number"
            min="1"
            max="31"
            {...register('due_day', { valueAsNumber: true })}
            className={errors.due_day ? 'border-destructive' : ''}
          />
          {errors.due_day && (
            <p className="text-sm text-destructive">{errors.due_day.message}</p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="color">Cor</Label>
        <div className="flex gap-2">
          <Input
            id="color"
            type="color"
            {...register('color')}
            className={errors.color ? 'border-destructive' : ''}
          />
          <Input
            {...register('color')}
            className={errors.color ? 'border-destructive' : ''}
            placeholder="#3b82f6"
          />
        </div>
        {errors.color && (
          <p className="text-sm text-destructive">{errors.color.message}</p>
        )}
      </div>

      <div className="flex gap-2">
        {creditCard && onDelete && (
          <Button
            type="button"
            variant="destructive"
            onClick={() => {
              if (confirm('Tem certeza que deseja excluir este cartão de crédito?')) {
                onDelete(creditCard.id)
              }
            }}
            disabled={isLoading}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Excluir
          </Button>
        )}
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? 'Salvando...' : creditCard ? 'Atualizar' : 'Criar'}
        </Button>
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
            Cancelar
          </Button>
        )}
      </div>
    </form>
  )
}
