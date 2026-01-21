'use client'

import { CreditCard as CreditCardIcon } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatCurrency } from '@/utils/currency'
import { useAvailableLimit } from '@/hooks/queries/useCreditCards'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { usePrivacyStore } from '@/stores/privacyStore'
import type { CreditCard } from '@/types/financial'

interface CreditCardCardProps {
  creditCard: CreditCard
  onClick?: () => void
}

export function CreditCardCard({ creditCard, onClick }: CreditCardCardProps) {
  const { data: availableLimitData } = useAvailableLimit(creditCard.id)
  const { data: accounts } = useAccounts()
  const { isPrivacyMode } = usePrivacyStore()
  const availableLimit = availableLimitData?.available_limit || 0
  const usedLimit = creditCard.limit_amount - availableLimit
  const usagePercentage = creditCard.limit_amount > 0 
    ? (usedLimit / creditCard.limit_amount) * 100 
    : 0

  const linkedAccount = accounts?.find((acc) => acc.id === creditCard.account_id)

  const displayLimit = isPrivacyMode ? '••••' : formatCurrency(creditCard.limit_amount, 'BRL')
  const displayAvailable = isPrivacyMode ? '••••' : formatCurrency(availableLimit, 'BRL')
  const displayUsed = isPrivacyMode ? '••••' : formatCurrency(usedLimit, 'BRL')

  // Estilo do AccountBadge usado em Transações
  const accountBadgeColors = linkedAccount
    ? (() => {
        const colorMap: Record<string, string> = {
          checking: 'bg-purple-100 text-purple-700 border-purple-200',
          savings: 'bg-blue-100 text-blue-700 border-blue-200',
          credit: 'bg-orange-100 text-orange-700 border-orange-200',
          investment: 'bg-emerald-100 text-emerald-700 border-emerald-200',
        }
        return colorMap[linkedAccount.type] || 'bg-gray-100 text-gray-700 border-gray-200'
      })()
    : ''

  return (
    <Card
      className={`cursor-pointer transition-all hover:shadow-md ${onClick ? '' : 'cursor-default'}`}
      onClick={onClick}
      style={{ borderLeftColor: creditCard.color, borderLeftWidth: '4px' }}
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="flex items-center gap-2">
          <CardTitle className="text-sm font-medium">{creditCard.name}</CardTitle>
          {linkedAccount && (
            <span className={`text-[10px] font-bold px-2 py-0.5 rounded border ${accountBadgeColors}`}>
              {linkedAccount.name}
            </span>
          )}
        </div>
        <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Limite Total</span>
            <span className="text-lg font-semibold">{displayLimit}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Disponível</span>
            <span className={`text-lg font-semibold ${availableLimit > 0 ? 'text-green-600' : 'text-red-600'}`}>
              {displayAvailable}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Utilizado</span>
            <span className="text-lg font-semibold">{displayUsed}</span>
          </div>
          <div className="space-y-1">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>Uso do limite</span>
              <span>{isPrivacyMode ? '••' : `${usagePercentage.toFixed(1)}%`}</span>
            </div>
            <div className="h-2 bg-secondary rounded-full overflow-hidden">
              <div
                className={`h-full transition-all ${
                  usagePercentage > 80 ? 'bg-red-500' : usagePercentage > 50 ? 'bg-yellow-500' : 'bg-green-500'
                }`}
                style={{ width: `${Math.min(usagePercentage, 100)}%` }}
              />
            </div>
          </div>
          <div className="flex items-center gap-4 text-xs text-muted-foreground">
            <div>
              <span>Fechamento: </span>
              <span className="font-medium">{creditCard.closing_day}</span>
            </div>
            <div>
              <span>Vencimento: </span>
              <span className="font-medium">{creditCard.due_day}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
