'use client'

import { NetWorthCard } from '@/components/financial/NetWorthCard'
import { UpcomingBillsCalendar } from '@/components/financial/UpcomingBillsCalendar'
import { SpendingByTagChart } from '@/components/financial/SpendingByTagChart'
import { PrivacyToggle } from '@/components/layout/PrivacyToggle'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Wallet, CreditCard, Upload, TrendingUp } from 'lucide-react'
import Link from 'next/link'

export default function FinancialPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Dashboard Financeiro</h1>
          <p className="text-muted-foreground">
            Visão geral das suas finanças e próximos vencimentos
          </p>
        </div>
        <PrivacyToggle />
      </div>

      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Contas</CardTitle>
            <Wallet className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">-</div>
            <p className="text-xs text-muted-foreground">Gerencie suas contas</p>
            <Link href="/financial/accounts">
              <Button variant="link" className="p-0 h-auto mt-2">
                Ver contas →
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cartões de Crédito</CardTitle>
            <CreditCard className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">-</div>
            <p className="text-xs text-muted-foreground">Gerencie seus cartões</p>
            <Link href="/financial/credit-cards">
              <Button variant="link" className="p-0 h-auto mt-2">
                Ver cartões →
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Transações</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">-</div>
            <p className="text-xs text-muted-foreground">Registre movimentações</p>
            <Link href="/financial/transactions">
              <Button variant="link" className="p-0 h-auto mt-2">
                Ver transações →
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Importar</CardTitle>
            <Upload className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">-</div>
            <p className="text-xs text-muted-foreground">Importe extratos</p>
            <Link href="/financial/import">
              <Button variant="link" className="p-0 h-auto mt-2">
                Importar →
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>

      {/* Patrimônio Líquido e Próximos Vencimentos */}
      <div className="grid gap-4 md:grid-cols-2">
        <NetWorthCard />
        <UpcomingBillsCalendar days={30} />
      </div>

      {/* Gráfico de Gastos por Tag */}
      <SpendingByTagChart />
    </div>
  )
}
