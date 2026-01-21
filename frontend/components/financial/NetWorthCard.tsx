'use client'

import { TrendingUp, TrendingDown, Wallet, CreditCard } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatCurrency } from '@/utils/currency'
import { useNetWorth } from '@/hooks/queries/useAnalytics'
import { usePrivacyStore } from '@/stores/privacyStore'

export function NetWorthCard() {
  const { data: netWorth, isLoading } = useNetWorth()
  const { isPrivacyMode } = usePrivacyStore()

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Patrimônio Líquido</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center p-4">Carregando...</div>
        </CardContent>
      </Card>
    )
  }

  if (!netWorth) {
    return null
  }

  const displayTotal = isPrivacyMode ? '••••' : formatCurrency(netWorth.total_patrimonio, 'BRL')
  const displayAccounts = isPrivacyMode ? '••••' : formatCurrency(netWorth.total_contas, 'BRL')
  const displayBills = isPrivacyMode ? '••••' : formatCurrency(netWorth.total_faturas, 'BRL')

  const isPositive = netWorth.total_patrimonio >= 0

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Wallet className="h-5 w-5" />
          Patrimônio Líquido
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Total</span>
            <div className="flex items-center gap-2">
              {isPositive ? (
                <TrendingUp className="h-4 w-4 text-green-600" />
              ) : (
                <TrendingDown className="h-4 w-4 text-red-600" />
              )}
              <span className={`text-2xl font-bold ${isPositive ? 'text-green-600' : 'text-red-600'}`}>
                {displayTotal}
              </span>
            </div>
          </div>

          <div className="space-y-2 pt-2 border-t">
            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <Wallet className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">Contas</span>
              </div>
              <span className="font-medium">{displayAccounts}</span>
            </div>

            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <CreditCard className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">Faturas em Aberto</span>
              </div>
              <span className="font-medium text-red-600">-{displayBills}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
