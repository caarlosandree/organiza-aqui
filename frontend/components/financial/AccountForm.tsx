'use client'

import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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
import { accountSchema, type AccountFormData } from '@/schemas/financialSchema'
import type { Account } from '@/types/financial'
import { useBanks } from '@/hooks/queries/useBanks'
import { currencies, getDefaultCurrency } from '@/utils/currencies'
import { BankCombobox } from '@/components/financial/BankCombobox'

interface AccountFormProps {
  account?: Account
  onSubmit: (data: AccountFormData) => void
  onCancel?: () => void
  isLoading?: boolean
}

export function AccountForm({ account, onSubmit, onCancel, isLoading }: AccountFormProps) {
  const { data: banks, isLoading: isLoadingBanks } = useBanks()
  const {
    register,
    handleSubmit,
    control,
    setValue,
    watch,
    formState: { errors },
  } = useForm<AccountFormData>({
    resolver: zodResolver(accountSchema),
    defaultValues: account
      ? {
          name: account.name,
          type: account.type,
          currency: account.currency,
          bank_id: account.bank_id,
        }
      : {
          currency: getDefaultCurrency().code,
        },
  })

  const accountType = watch('type')

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Nome da Conta</Label>
        <Input
          id="name"
          {...register('name')}
          className={errors.name ? 'border-destructive' : ''}
          placeholder="Ex: Conta Principal"
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="type">Tipo</Label>
        <Select
          value={accountType}
          onValueChange={(value) => setValue('type', value as AccountFormData['type'])}
        >
          <SelectTrigger id="type" className={errors.type ? 'border-destructive' : ''}>
            <SelectValue placeholder="Selecione o tipo" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="checking">Conta Corrente</SelectItem>
            <SelectItem value="savings">Poupança</SelectItem>
            <SelectItem value="credit">Cartão de Crédito</SelectItem>
            <SelectItem value="investment">Investimento</SelectItem>
          </SelectContent>
        </Select>
        {errors.type && (
          <p className="text-sm text-destructive">{errors.type.message}</p>
        )}
      </div>

      <BankCombobox
        control={control}
        name="bank_id"
        banks={banks}
        isLoading={isLoadingBanks}
      />

      <div className="space-y-2">
        <Label htmlFor="currency">Moeda</Label>
        <Controller
          name="currency"
          control={control}
          render={({ field, fieldState }) => (
            <>
              <Select
                value={field.value || getDefaultCurrency().code}
                onValueChange={field.onChange}
              >
                <SelectTrigger
                  id="currency"
                  className={fieldState.error ? 'border-destructive' : ''}
                >
                  <SelectValue placeholder="Selecione a moeda" />
                </SelectTrigger>
                <SelectContent>
                  {currencies.map((currency) => (
                    <SelectItem key={currency.code} value={currency.code}>
                      {currency.code} - {currency.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {fieldState.error && (
                <p className="text-sm text-destructive">{fieldState.error.message}</p>
              )}
            </>
          )}
        />
      </div>

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? 'Salvando...' : account ? 'Atualizar' : 'Criar'}
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
