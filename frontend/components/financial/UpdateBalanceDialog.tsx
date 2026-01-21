'use client'

import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { CalendarIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { Field, FieldLabel } from '@/components/ui/field'
import { updateInitialBalanceSchema, type UpdateInitialBalanceFormData } from '@/schemas/financialSchema'
import type { Account } from '@/types/financial'
import { formatCurrency } from '@/utils/currency'
import { useUpdateInitialBalance } from '@/hooks/mutations/useAccountMutations'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface UpdateBalanceDialogProps {
  account: Account
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function UpdateBalanceDialog({ account, open, onOpenChange }: UpdateBalanceDialogProps) {
  const [balanceInput, setBalanceInput] = useState('')
  const [datePickerOpen, setDatePickerOpen] = useState(false)
  const updateMutation = useUpdateInitialBalance()

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    watch,
    setValue,
  } = useForm<UpdateInitialBalanceFormData>({
    resolver: zodResolver(updateInitialBalanceSchema),
    defaultValues: {
      balance: 0,
      date: new Date().toISOString().split('T')[0],
    },
  })

  const selectedDate = watch('date')

  // Reset form when dialog opens/closes
  useEffect(() => {
    if (open) {
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      reset({
        balance: account.balance,
        date: today.toISOString().split('T')[0],
      })
      setBalanceInput(formatCurrency(account.balance / 100, account.currency))
    }
  }, [open, account, reset])

  const handleBalanceChange = (value: string) => {
    // Remove tudo que não é número (mantém apenas dígitos)
    const numbers = value.replace(/\D/g, '')

    // Se não há números, limpar
    if (!numbers) {
      setBalanceInput('')
      setValue('balance', 0, { shouldValidate: true })
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
      currency: account.currency || 'BRL',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(reais)

    setBalanceInput(formatted)

    // Atualizar o campo do formulário com centavos
    setValue('balance', cents, { shouldValidate: true })
  }

  const onSubmit = async (data: UpdateInitialBalanceFormData) => {
    try {
      await updateMutation.mutateAsync({
        id: account.id,
        data: {
          balance: data.balance,
          date: data.date,
        },
      })
      onOpenChange(false)
    } catch (error) {
      console.error('Erro ao atualizar saldo:', error)
    }
  }

  const today = new Date()
  today.setHours(23, 59, 59, 999)
  const maxDate = today.toISOString().split('T')[0]

  // Converter string de data para Date object
  const selectedDateObj = selectedDate ? new Date(selectedDate + 'T00:00:00') : undefined

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Atualizar Saldo Inicial</DialogTitle>
          <DialogDescription>
            Informe o saldo atual da conta e a data de referência. O sistema criará uma transação de ajuste se necessário.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="balance">Saldo Atual</Label>
            <Input
              id="balance"
              type="text"
              value={balanceInput}
              onChange={(e) => handleBalanceChange(e.target.value)}
              placeholder="R$ 0,00"
              className={errors.balance ? 'border-red-500' : ''}
            />
            {errors.balance && (
              <p className="text-sm text-red-500">{errors.balance.message}</p>
            )}
            <input type="hidden" {...register('balance')} />
          </div>

          <Field className="space-y-2">
            <FieldLabel htmlFor="date">Data de Referência</FieldLabel>
            <Controller
              control={control}
              name="date"
              render={({ field }) => {
                const dateValue = field.value ? new Date(field.value + 'T00:00:00') : undefined
                return (
                  <Popover open={datePickerOpen} onOpenChange={setDatePickerOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        id="date"
                        className={cn(
                          'w-full justify-start text-left font-normal',
                          !dateValue && 'text-muted-foreground',
                          errors.date && 'border-red-500'
                        )}
                      >
                        <CalendarIcon className="mr-2 h-4 w-4" />
                        {dateValue ? (
                          format(dateValue, "PPP", { locale: ptBR })
                        ) : (
                          <span>Selecione uma data</span>
                        )}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent 
                      className="w-auto p-0 max-h-[calc(100vh-8rem)] overflow-y-auto" 
                      align="start"
                      side="bottom"
                      sideOffset={4}
                    >
                      <Calendar
                        mode="single"
                        selected={dateValue}
                        defaultMonth={dateValue}
                        captionLayout="dropdown"
                        onSelect={(date) => {
                          if (date) {
                            const dateStr = date.toISOString().split('T')[0]
                            field.onChange(dateStr)
                            setDatePickerOpen(false)
                          }
                        }}
                        disabled={(date) => date > today}
                      />
                    </PopoverContent>
                  </Popover>
                )
              }}
            />
            {errors.date && (
              <p className="text-sm text-red-500">{errors.date.message}</p>
            )}
            <p className="text-xs text-muted-foreground">
              Data em que o saldo informado estava válido
            </p>
          </Field>

          {selectedDate && (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Se houver transações anteriores a esta data, uma transação de ajuste será criada automaticamente.
              </AlertDescription>
            </Alert>
          )}

          {updateMutation.isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Erro ao atualizar saldo. Tente novamente.
              </AlertDescription>
            </Alert>
          )}

          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={updateMutation.isPending}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={updateMutation.isPending}>
              {updateMutation.isPending ? 'Atualizando...' : 'Atualizar Saldo'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
