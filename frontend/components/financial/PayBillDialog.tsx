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
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { payBillSchema, type PayBillFormData } from '@/schemas/financialSchema'
import { useAccounts } from '@/hooks/queries/useAccounts'

interface PayBillDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (data: PayBillFormData) => void
  isLoading?: boolean
}

export function PayBillDialog({ open, onOpenChange, onSubmit, isLoading }: PayBillDialogProps) {
  const { data: accounts } = useAccounts()

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<PayBillFormData>({
    resolver: zodResolver(payBillSchema),
    defaultValues: {
      date: new Date().toISOString().split('T')[0],
    },
  })

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Pagar Fatura</DialogTitle>
          <DialogDescription>
            Selecione a conta bancária e a data do pagamento
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="account_id">Conta de Pagamento</Label>
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
            <Label htmlFor="date">Data do Pagamento</Label>
            <Controller
              name="date"
              control={control}
              render={({ field }) => (
                <Input
                  id="date"
                  type="date"
                  value={field.value}
                  onChange={field.onChange}
                  className={errors.date ? 'border-destructive' : ''}
                />
              )}
            />
            {errors.date && (
              <p className="text-sm text-destructive">{errors.date.message}</p>
            )}
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={isLoading} className="flex-1">
              {isLoading ? 'Processando...' : 'Pagar'}
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
