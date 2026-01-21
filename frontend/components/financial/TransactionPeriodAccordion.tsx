'use client'

import { useState } from 'react'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { Calendar, ChevronDown, Wallet, CreditCard } from 'lucide-react'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
// Tabs simples usando estado
import { formatCurrency } from '@/utils/currency'
import { TransactionTable } from './TransactionTable'
import type { PeriodWithTransactions, Transaction } from '@/types/financial'

interface TransactionPeriodAccordionProps {
  periods: PeriodWithTransactions[]
  onEdit?: (transaction: Transaction) => void
  onDelete?: (id: string) => void
}

interface GroupedPeriod {
  year: number
  month: number
  periods: PeriodWithTransactions[]
}

// Agrupar períodos por ano e mês
function groupPeriodsByMonth(periods: PeriodWithTransactions[]): GroupedPeriod[] {
  const grouped = new Map<string, GroupedPeriod>()

  for (const period of periods) {
    const key = `${period.period.year}-${period.period.month}`
    if (!grouped.has(key)) {
      grouped.set(key, {
        year: period.period.year,
        month: period.period.month,
        periods: [],
      })
    }
    grouped.get(key)!.periods.push(period)
  }

  return Array.from(grouped.values()).sort((a, b) => {
    if (a.year !== b.year) return b.year - a.year
    return b.month - a.month
  })
}

export function TransactionPeriodAccordion({
  periods,
  onEdit,
  onDelete,
}: TransactionPeriodAccordionProps) {
  const groupedPeriods = groupPeriodsByMonth(periods)
  const [openMonths, setOpenMonths] = useState<Set<string>>(
    new Set([groupedPeriods[0] ? `${groupedPeriods[0].year}-${groupedPeriods[0].month}` : ''])
  )
  const [activeTabs, setActiveTabs] = useState<Map<string, string>>(
    new Map(groupedPeriods.map((g) => [`${g.year}-${g.month}`, 'overview']))
  )

  const toggleMonth = (key: string) => {
    setOpenMonths((prev) => {
      const next = new Set(prev)
      if (next.has(key)) {
        next.delete(key)
      } else {
        next.add(key)
      }
      return next
    })
  }

  const setActiveTab = (monthKey: string, tab: string) => {
    setActiveTabs((prev) => {
      const next = new Map(prev)
      next.set(monthKey, tab)
      return next
    })
  }

  const getActiveTab = (monthKey: string) => {
    return activeTabs.get(monthKey) || 'overview'
  }

  if (groupedPeriods.length === 0) {
    return (
      <div className="text-center p-8 text-muted-foreground">
        <p>Nenhum período encontrado.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {groupedPeriods.map((group) => {
        const monthKey = `${group.year}-${group.month}`
        const isOpen = openMonths.has(monthKey)
        const monthDate = new Date(group.year, group.month - 1, 1)

        // Calcular totais do mês
        let totalIncome = 0
        let totalBankExpense = 0
        let totalCreditCardExpense = 0
        let totalTransactions = 0

        for (const period of group.periods) {
          totalIncome += period.total_income
          totalBankExpense += period.total_bank_expense
          totalCreditCardExpense += period.total_credit_card_expense
          totalTransactions += period.transactions.length
        }

        const balance = totalIncome - totalBankExpense

        // Separar transações por tipo de período
        const bankTransactions: Transaction[] = []
        const creditCardTransactions: Transaction[] = []
        const allTransactions: Transaction[] = []

        for (const period of group.periods) {
          if (period.period.period_type === 'bank') {
            bankTransactions.push(...period.transactions)
          } else {
            creditCardTransactions.push(...period.transactions)
          }
          allTransactions.push(...period.transactions)
        }

        // Ordenar transações por data
        const sortByDate = (a: Transaction, b: Transaction) =>
          new Date(b.date).getTime() - new Date(a.date).getTime()
        allTransactions.sort(sortByDate)
        bankTransactions.sort(sortByDate)
        creditCardTransactions.sort(sortByDate)

        return (
          <Collapsible key={monthKey} open={isOpen} onOpenChange={() => toggleMonth(monthKey)}>
            <Card>
              <CollapsibleTrigger className="w-full">
                <CardHeader className="cursor-pointer hover:bg-muted/50 transition-colors">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Calendar className="h-5 w-5 text-muted-foreground" />
                      <CardTitle className="text-lg">
                        {format(monthDate, "MMMM 'de' yyyy", { locale: ptBR })}
                      </CardTitle>
                      <Badge variant="secondary">{totalTransactions} lançamentos</Badge>
                    </div>
                    <div className="flex items-center gap-6">
                      <div className="text-right">
                        <div className="text-sm text-muted-foreground">Entradas</div>
                        <div className="text-sm font-medium text-[var(--success)]">
                          {formatCurrency(totalIncome)}
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-sm text-muted-foreground">Saídas</div>
                        <div className="text-sm font-medium text-destructive">
                          {formatCurrency(totalBankExpense + totalCreditCardExpense)}
                        </div>
                      </div>
                      <ChevronDown
                        className={`h-5 w-5 text-muted-foreground transition-transform ${
                          isOpen ? 'transform rotate-180' : ''
                        }`}
                      />
                    </div>
                  </div>
                </CardHeader>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <CardContent className="pt-0">
                  {/* Dashboard do Mês */}
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                    <Card className="p-4">
                      <div className="text-sm text-muted-foreground mb-1">Receitas</div>
                      <div className="text-2xl font-bold text-[var(--success)]">
                        {formatCurrency(totalIncome)}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">Banco</div>
                    </Card>
                    <Card className="p-4">
                      <div className="text-sm text-muted-foreground mb-1">Despesas Conta</div>
                      <div className="text-2xl font-bold text-destructive">
                        {formatCurrency(totalBankExpense)}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">Banco</div>
                    </Card>
                    <Card className="p-4">
                      <div className="text-sm text-muted-foreground mb-1">Fatura Cartão</div>
                      <div className="text-2xl font-bold text-destructive">
                        {formatCurrency(totalCreditCardExpense)}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">Cartão</div>
                    </Card>
                    <Card className="p-4">
                      <div className="text-sm text-muted-foreground mb-1">Resumo</div>
                      <div
                        className={`text-2xl font-bold ${
                          balance >= 0 ? 'text-[var(--success)]' : 'text-destructive'
                        }`}
                      >
                        {formatCurrency(balance)}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        Receitas - Despesas Banco
                      </div>
                    </Card>
                  </div>

                  {/* Tabs para filtrar visualização */}
                  <div className="w-full">
                    <div className="flex border-b mb-4">
                      <button
                        onClick={() => setActiveTab(monthKey, 'overview')}
                        className={`px-4 py-2 border-b-2 transition-colors ${
                          getActiveTab(monthKey) === 'overview'
                            ? 'border-primary text-primary font-medium'
                            : 'border-transparent text-muted-foreground hover:text-foreground'
                        }`}
                      >
                        Visão Geral
                      </button>
                      <button
                        onClick={() => setActiveTab(monthKey, 'bank')}
                        className={`px-4 py-2 border-b-2 transition-colors flex items-center gap-2 ${
                          getActiveTab(monthKey) === 'bank'
                            ? 'border-primary text-primary font-medium'
                            : 'border-transparent text-muted-foreground hover:text-foreground'
                        }`}
                      >
                        <Wallet className="h-4 w-4" />
                        Extrato Bancário
                      </button>
                      <button
                        onClick={() => setActiveTab(monthKey, 'credit')}
                        className={`px-4 py-2 border-b-2 transition-colors flex items-center gap-2 ${
                          getActiveTab(monthKey) === 'credit'
                            ? 'border-primary text-primary font-medium'
                            : 'border-transparent text-muted-foreground hover:text-foreground'
                        }`}
                      >
                        <CreditCard className="h-4 w-4" />
                        Fatura Cartão
                      </button>
                    </div>
                    {getActiveTab(monthKey) === 'overview' && (
                      <TransactionTable
                        transactions={allTransactions}
                        onEdit={onEdit}
                        onDelete={onDelete}
                      />
                    )}
                    {getActiveTab(monthKey) === 'bank' && (
                      <>
                        {bankTransactions.length > 0 ? (
                          <TransactionTable
                            transactions={bankTransactions}
                            onEdit={onEdit}
                            onDelete={onDelete}
                          />
                        ) : (
                          <div className="text-center p-8 text-muted-foreground">
                            <p>Nenhuma transação bancária neste período.</p>
                          </div>
                        )}
                      </>
                    )}
                    {getActiveTab(monthKey) === 'credit' && (
                      <>
                        {creditCardTransactions.length > 0 ? (
                          <TransactionTable
                            transactions={creditCardTransactions}
                            onEdit={onEdit}
                            onDelete={onDelete}
                          />
                        ) : (
                          <div className="text-center p-8 text-muted-foreground">
                            <p>Nenhuma transação de cartão neste período.</p>
                          </div>
                        )}
                      </>
                    )}
                  </div>
                </CardContent>
              </CollapsibleContent>
            </Card>
          </Collapsible>
        )
      })}
    </div>
  )
}
