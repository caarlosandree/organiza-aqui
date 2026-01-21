'use client'

import { Wallet } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatCurrency } from '@/utils/currency'
import type { Account } from '@/types/financial'

interface AccountCardProps {
  account: Account
  onClick?: () => void
}

const accountTypeLabels: Record<Account['type'], string> = {
  checking: 'Conta Corrente',
  savings: 'Poupança',
  credit: 'Cartão de Crédito',
  investment: 'Investimento',
}

export function AccountCard({ account, onClick }: AccountCardProps) {
  const isCreditCard = account.type === 'credit'
  
  // Para cartão: saldo negativo = dívida, saldo positivo = pagamento excedente
  // Para banco: saldo positivo = dinheiro, saldo negativo = débito
  const balanceColor = isCreditCard
    ? account.balance <= 0 
      ? 'text-muted-foreground'  // Dívida (neutro, não alarmante)
      : 'text-[var(--success)]'  // Pagamento excedente
    : account.balance >= 0 
      ? 'text-[var(--success)]'  // Dinheiro disponível
      : 'text-destructive'  // Débito (alarmante)

  const balanceLabel = isCreditCard
    ? account.balance <= 0
      ? 'Dívida'
      : 'Saldo a haver'  // Não usar "Crédito Disponível" (confunde com limite)
    : 'Saldo'

  return (
    <Card
      className={`cursor-pointer transition-all hover:shadow-md ${onClick ? '' : 'cursor-default'}`}
      onClick={onClick}
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{account.name}</CardTitle>
        <Wallet className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between">
            <Badge variant="outline">{accountTypeLabels[account.type]}</Badge>
            <div className="flex flex-col items-end">
              <span className="text-xs text-muted-foreground">{balanceLabel}</span>
              <span className={`text-2xl font-bold ${balanceColor}`}>
                {formatCurrency(Math.abs(account.balance), account.currency)}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
