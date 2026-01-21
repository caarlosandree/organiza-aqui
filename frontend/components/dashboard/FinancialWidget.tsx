'use client'

import { Wallet, TrendingUp, TrendingDown } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { formatCurrency } from '@/utils/currency'
import { Skeleton } from '@/components/ui/skeleton'

export const FinancialWidget = () => {
  const { data: accounts = [], isLoading } = useAccounts()

  const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0)

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-24" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Financeiro</CardTitle>
        <Wallet className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{formatCurrency(totalBalance)}</div>
        <p className="text-xs text-muted-foreground mt-1">
          {accounts.length} conta{accounts.length !== 1 ? 's' : ''}
        </p>
        <div className="flex items-center gap-4 mt-4">
          <div className="flex items-center gap-1 text-[var(--success)]">
            <TrendingUp className="h-3 w-3" />
            <span className="text-xs">Receitas</span>
          </div>
          <div className="flex items-center gap-1 text-destructive">
            <TrendingDown className="h-3 w-3" />
            <span className="text-xs">Despesas</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
