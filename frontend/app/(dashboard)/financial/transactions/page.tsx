'use client'

import { useState, useMemo } from 'react'
import { useQueryStates, parseAsString, parseAsInteger } from 'nuqs'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import {
  Plus,
  CreditCard,
  Landmark,
  ChevronDown,
  Calendar,
  Trash2,
  Edit,
  Filter,
  PieChart,
  Target,
  Loader2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { TransactionForm } from '@/components/financial/TransactionForm'
import { useTransactions } from '@/hooks/queries/useTransactions'
import { useTransactionPeriods } from '@/hooks/queries/useTransactionPeriods'
import {
  useCreateTransaction,
  useUpdateTransaction,
  useDeleteTransaction,
} from '@/hooks/mutations/useTransactionMutations'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useCategories } from '@/hooks/queries/useCategories'
import { formatCurrency } from '@/utils/currency'
import { cn } from '@/lib/utils'
import type { TransactionFormData } from '@/schemas/financialSchema'
import type {
  TransactionFilters,
  TransactionPeriodFilters,
  Transaction,
  PeriodWithTransactions,
  Account,
  Category,
} from '@/types/financial'

// --- Helpers ---

// Parse uma string de data no formato "YYYY-MM-DD" como data local (sem conversão de fuso horário)
const parseLocalDate = (dateString: string): Date => {
  const [year, month, day] = dateString.split('-').map(Number)
  return new Date(year, month - 1, day)
}

// Extrai o mês no formato "yyyy-MM" de uma string de data "YYYY-MM-DD"
const getMonthFromDate = (dateString: string): string => {
  const [year, month] = dateString.split('-')
  return `${year}-${month}`
}

// --- Componentes ---

interface KpiCardProps {
  title: string
  value: number
  subtext?: string
  icon: React.ElementType
  trend?: 'up' | 'down'
}

const KpiCard = ({ title, value, subtext, icon: Icon, trend }: KpiCardProps) => {
  const isPositive = trend === 'up'
  const isNegative = trend === 'down'

  // Determinar cores baseadas na tendência
  let bgColor = 'bg-blue-50/50 dark:bg-blue-950/30'
  let iconColor = 'text-blue-600 dark:text-blue-400'
  let valueColor = 'text-blue-600 dark:text-blue-400'

  if (isPositive) {
    bgColor = 'bg-emerald-50/50 dark:bg-emerald-950/30'
    iconColor = 'text-emerald-600 dark:text-emerald-400'
    valueColor = 'text-emerald-600 dark:text-emerald-400'
  } else if (isNegative) {
    bgColor = 'bg-red-50/50 dark:bg-red-950/30'
    iconColor = 'text-red-600 dark:text-red-400'
    valueColor = 'text-red-600 dark:text-red-400'
  }

  return (
    <Card className={cn('flex flex-col justify-between h-full border-0', bgColor)}>
      <CardContent className="p-5">
        <div className="flex justify-between items-start mb-2">
          <p className="text-xs text-muted-foreground font-bold uppercase tracking-wider">
            {title}
          </p>
          <div className={cn('p-2 rounded-lg', iconColor)}>
            <Icon size={18} />
          </div>
        </div>
        <div>
          <h3 className={cn('text-2xl font-bold', valueColor)}>{formatCurrency(value)}</h3>
          {subtext && <p className="text-xs text-muted-foreground mt-1">{subtext}</p>}
        </div>
      </CardContent>
    </Card>
  )
}

interface AccountBadgeProps {
  account: Account
}

const AccountBadge = ({ account }: AccountBadgeProps) => {
  const colorMap: Record<string, string> = {
    checking: 'bg-purple-100 text-purple-700 border-purple-200',
    savings: 'bg-blue-100 text-blue-700 border-blue-200',
    credit: 'bg-orange-100 text-orange-700 border-orange-200',
    investment: 'bg-emerald-100 text-emerald-700 border-emerald-200',
  }

  const colors = colorMap[account.type] || 'bg-gray-100 text-gray-700 border-gray-200'

  return (
    <span className={`text-[10px] font-bold px-2 py-0.5 rounded border ${colors}`}>
      {account.name}
    </span>
  )
}

interface TransactionTableProps {
  transactions: Transaction[]
  accounts: Account[]
  categories: Category[]
  onEdit: (transaction: Transaction) => void
  onDelete: (id: string) => void
  emptyMessage: string
}

const TransactionTable = ({
  transactions,
  accounts,
  categories,
  onEdit,
  onDelete,
  emptyMessage,
}: TransactionTableProps) => {
  if (transactions.length === 0) {
    return (
      <div className="p-8 text-center text-muted-foreground border border-dashed rounded-lg bg-muted/50 text-sm">
        {emptyMessage}
      </div>
    )
  }

  const getAccountName = (accountId: string) => {
    return accounts.find((a) => a.id === accountId)?.name || 'Conta'
  }

  const getCategoryName = (categoryId?: string) => {
    if (!categoryId) return 'Outros'
    return categories.find((c) => c.id === categoryId)?.name || 'Outros'
  }

  const formatDateShort = (dateString: string) => {
    const date = parseLocalDate(dateString)
    return format(date, 'dd/MM/yyyy', { locale: ptBR })
  }

  return (
    <div className="overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Dia</TableHead>
            <TableHead>Conta</TableHead>
            <TableHead>Descrição</TableHead>
            <TableHead>Categoria</TableHead>
            <TableHead className="text-right">Valor</TableHead>
            <TableHead className="text-right">Ações</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {transactions.map((t) => (
            <TableRow key={t.id}>
              <TableCell className="text-muted-foreground font-mono text-xs whitespace-nowrap">
                {formatDateShort(t.date)}
              </TableCell>
              <TableCell>
                {(() => {
                  const account = accounts.find((a) => a.id === t.account_id)
                  return account ? <AccountBadge account={account} /> : <span>N/A</span>
                })()}
              </TableCell>
              <TableCell className="font-medium">{t.description || '-'}</TableCell>
              <TableCell>
                <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-muted text-muted-foreground">
                  {getCategoryName(t.category_id || undefined)}
                </span>
              </TableCell>
              <TableCell
                className={`text-right font-bold whitespace-nowrap ${
                  t.type === 'income' ? 'text-emerald-600' : 'text-rose-600'
                }`}
              >
                {t.type === 'expense' ? '-' : '+'} {formatCurrency(Math.abs(t.amount))}
              </TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => onEdit(t)}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:text-destructive"
                    onClick={() => onDelete(t.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

interface MonthAccordionProps {
  monthKey: string
  transactions: Transaction[]
  accounts: Account[]
  categories: Category[]
  onEdit: (transaction: Transaction) => void
  onDelete: (id: string) => void
  isOpenDefault?: boolean
  showAccountTags?: boolean
}

const MonthAccordion = ({
  monthKey,
  transactions,
  accounts,
  categories,
  onEdit,
  onDelete,
  isOpenDefault = false,
  showAccountTags = false,
}: MonthAccordionProps) => {
  const [isOpen, setIsOpen] = useState(isOpenDefault)
  const [activeTab, setActiveTab] = useState<'all' | 'bank' | 'card'>('all')

  // Obter contas únicas envolvidas neste mês
  const uniqueAccounts = useMemo(() => {
    if (!showAccountTags) return []
    const accountIds = new Set(transactions.map((t) => t.account_id))
    return accounts.filter((acc) => accountIds.has(acc.id))
  }, [transactions, accounts, showAccountTags])

  const stats = useMemo(() => {
    const bankTx = transactions.filter(
      (t) => accounts.find((a) => a.id === t.account_id)?.type !== 'credit'
    )
    const cardTx = transactions.filter(
      (t) => accounts.find((a) => a.id === t.account_id)?.type === 'credit'
    )

    const bankIncome = bankTx.filter((t) => t.type === 'income').reduce((acc, t) => acc + t.amount, 0)
    const bankExpense = bankTx.filter((t) => t.type === 'expense').reduce((acc, t) => acc + t.amount, 0)
    const cardBill = cardTx.reduce((acc, t) => acc + t.amount, 0)

    return { bankIncome, bankExpense, cardBill, balance: bankIncome - bankExpense }
  }, [transactions, accounts])

  const filteredList = useMemo(() => {
    let filtered = transactions
    if (activeTab === 'bank') {
      filtered = transactions.filter(
        (t) => accounts.find((a) => a.id === t.account_id)?.type !== 'credit'
      )
    } else if (activeTab === 'card') {
      filtered = transactions.filter(
        (t) => accounts.find((a) => a.id === t.account_id)?.type === 'credit'
      )
    }
    return filtered.sort(
      (a, b) => parseLocalDate(b.date).getTime() - parseLocalDate(a.date).getTime()
    )
  }, [transactions, activeTab, accounts])

  const getMonthName = (monthKey: string) => {
    if (!monthKey) return ''
    const [year, month] = monthKey.split('-')
    const date = new Date(parseInt(year), parseInt(month) - 1, 1)
    return format(date, "MMMM 'de' yyyy", { locale: ptBR })
  }

  return (
    <Card className="mb-4 overflow-hidden py-0 gap-0">
      <Collapsible open={isOpen} onOpenChange={setIsOpen}>
        <CollapsibleTrigger className="w-full">
          <CardContent className="p-4 cursor-pointer hover:bg-muted/50 transition-colors flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div
                className={`w-10 h-10 rounded-full flex items-center justify-center ${
                  isOpen ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'
                }`}
              >
                <Calendar size={18} />
              </div>
              <div className="flex flex-col items-start">
                <h3 className="text-lg font-bold leading-tight capitalize">
                  {getMonthName(monthKey)}
                </h3>
                <div className="flex items-center gap-2 mt-1">
                  <p className="text-xs text-muted-foreground leading-none">
                    {transactions.length} lançamentos
                  </p>
                  {showAccountTags && uniqueAccounts.length > 0 && (
                    <div className="flex items-center gap-1 flex-wrap">
                      {uniqueAccounts.slice(0, 3).map((account) => (
                        <AccountBadge key={account.id} account={account} />
                      ))}
                      {uniqueAccounts.length > 3 && (
                        <span className="text-[10px] text-muted-foreground font-medium">
                          +{uniqueAccounts.length - 3}
                        </span>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-center gap-6">
              <div className="hidden md:block text-right">
                <p className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">
                  Entradas
                </p>
                <p className="text-sm font-bold text-emerald-600">
                  {formatCurrency(stats.bankIncome)}
                </p>
              </div>
              <div className="hidden md:block text-right">
                <p className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">
                  Saídas Totais
                </p>
                <p className="text-sm font-bold text-rose-600">
                  {formatCurrency(stats.bankExpense + stats.cardBill)}
                </p>
              </div>
              <div
                className={`p-1 rounded-full transition-transform ${
                  isOpen ? 'rotate-180 bg-muted' : ''
                }`}
              >
                <ChevronDown size={20} className="text-muted-foreground" />
              </div>
            </div>
          </CardContent>
        </CollapsibleTrigger>

        <CollapsibleContent>
          <div className="border-t bg-background">
            {/* Dashboard Mês */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-px bg-muted border-b">
              <div className="bg-background p-4">
                <span className="text-xs text-muted-foreground font-semibold block mb-1">
                  Receitas
                </span>
                <span className="text-lg font-bold text-emerald-600">
                  {formatCurrency(stats.bankIncome)}
                </span>
              </div>
              <div className="bg-background p-4">
                <span className="text-xs text-muted-foreground font-semibold block mb-1">
                  Despesas Conta
                </span>
                <span className="text-lg font-bold text-rose-600">
                  {formatCurrency(stats.bankExpense)}
                </span>
              </div>
              <div className="bg-background p-4">
                <span className="text-xs text-muted-foreground font-semibold block mb-1">
                  Fatura Cartão
                </span>
                <span className="text-lg font-bold text-rose-600">
                  {formatCurrency(stats.cardBill)}
                </span>
              </div>
              <div className="bg-background p-4">
                <span className="text-xs text-muted-foreground font-semibold block mb-1">
                  Resumo (S/ Fatura)
                </span>
                <span
                  className={`text-lg font-bold ${
                    stats.balance >= 0 ? 'text-emerald-600' : 'text-rose-600'
                  }`}
                >
                  {formatCurrency(stats.balance)}
                </span>
              </div>
            </div>

            {/* Tabs */}
            <div className="px-4 pt-4">
              <div className="flex gap-2 border-b">
                {(['all', 'bank', 'card'] as const).map((tab) => (
                  <button
                    key={tab}
                    onClick={() => setActiveTab(tab)}
                    className={`pb-3 px-4 text-sm font-medium transition-all relative ${
                      activeTab === tab
                        ? 'text-primary'
                        : 'text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    {tab === 'all' && 'Visão Geral'}
                    {tab === 'bank' && 'Extrato Bancário'}
                    {tab === 'card' && 'Fatura Cartão'}
                    {activeTab === tab && (
                      <div className="absolute bottom-0 left-0 w-full h-0.5 bg-primary rounded-t-full"></div>
                    )}
                  </button>
                ))}
              </div>
            </div>

            <div className="min-h-[100px] p-4 pb-4">
              <TransactionTable
                transactions={filteredList}
                accounts={accounts}
                categories={categories}
                onEdit={onEdit}
                onDelete={onDelete}
                emptyMessage="Nenhuma transação encontrada neste filtro."
              />
            </div>
          </div>
        </CollapsibleContent>
      </Collapsible>
    </Card>
  )
}

// --- Página Principal ---

export default function TransactionsPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null)

  // Filtros na URL usando Nuqs
  const [filters, setFilters] = useQueryStates({
    account_id: parseAsString,
    reference_month: parseAsString,
  })

  // Queries
  const { data: accounts, isLoading: isLoadingAccounts } = useAccounts()
  const { data: categories, isLoading: isLoadingCategories } = useCategories()

  const transactionFilters: TransactionFilters = {
    account_id: filters.account_id || undefined,
    limit: 1000, // Buscar muitas transações para agrupar
    offset: 0,
  }

  const periodFilters: TransactionPeriodFilters = {
    account_id: filters.account_id || undefined,
    reference_month: filters.reference_month || undefined,
  }

  const { data: transactionsData, isLoading: isLoadingTransactions } =
    useTransactions(transactionFilters)
  const { data: periods, isLoading: isLoadingPeriods } =
    useTransactionPeriods(periodFilters)

  // Mutations
  const createMutation = useCreateTransaction()
  const updateMutation = useUpdateTransaction()
  const deleteMutation = useDeleteTransaction()

  // Agrupar transações por mês de referência
  const groupedTransactions = useMemo(() => {
    if (!transactionsData?.data) return []

    const grouped = transactionsData.data.reduce(
      (acc, t) => {
        const monthKey = t.reference_month || getMonthFromDate(t.date)
        if (!acc[monthKey]) acc[monthKey] = []
        acc[monthKey].push(t)
        return acc
      },
      {} as Record<string, Transaction[]>
    )

    return Object.keys(grouped)
      .sort()
      .reverse()
      .map((key) => ({
        monthKey: key,
        transactions: grouped[key],
      }))
  }, [transactionsData])

  // KPIs Globais
  const globalStats = useMemo(() => {
    if (!transactionsData?.data || !accounts) {
      return { currentBalance: 0, creditCardDebt: 0, filteredIncome: 0, filteredExpense: 0 }
    }

    // Filtrar transações por conta selecionada
    const accountTxs =
      filters.account_id === null || filters.account_id === undefined
        ? transactionsData.data
        : transactionsData.data.filter((t) => t.account_id === filters.account_id)

    // Calcular saldo atual (receitas - despesas bancárias) de todas as transações
    const currentBalance = accountTxs
      .filter((t) => {
        const account = accounts.find((a) => a.id === t.account_id)
        return account?.type !== 'credit'
      })
      .reduce((acc, t) => {
        if (t.type === 'income') return acc + t.amount
        if (t.type === 'expense') return acc - t.amount
        return acc
      }, 0)

    // Faturas de cartão (transações em contas tipo credit)
    const creditCardDebt = accountTxs
      .filter((t) => {
        const account = accounts.find((a) => a.id === t.account_id)
        return account?.type === 'credit'
      })
      .reduce((acc, t) => acc + t.amount, 0)

    // Receitas e despesas do período filtrado
    const filteredTxs = filters.reference_month
      ? accountTxs.filter((t) => t.reference_month === filters.reference_month)
      : accountTxs

    const filteredIncome = filteredTxs
      .filter((t) => t.type === 'income')
      .reduce((acc, t) => acc + t.amount, 0)

    const filteredExpense = filteredTxs
      .filter((t) => {
        const account = accounts.find((a) => a.id === t.account_id)
        return account?.type !== 'credit' && t.type === 'expense'
      })
      .reduce((acc, t) => acc + t.amount, 0)

    return { currentBalance, creditCardDebt, filteredIncome, filteredExpense }
  }, [transactionsData, accounts, filters])

  const isLoading =
    isLoadingAccounts || isLoadingCategories || isLoadingTransactions || isLoadingPeriods

  const handleSubmit = async (data: TransactionFormData) => {
    try {
      if (editingTransaction) {
        await updateMutation.mutateAsync({ id: editingTransaction.id, data })
      } else {
        await createMutation.mutateAsync(data)
      }
      setIsDialogOpen(false)
      setEditingTransaction(null)
    } catch (error) {
      console.error('Erro ao salvar transação:', error)
    }
  }

  const handleEdit = (transaction: Transaction) => {
    setEditingTransaction(transaction)
    setIsDialogOpen(true)
  }

  const handleDelete = async (id: string) => {
    if (confirm('Tem certeza que deseja deletar esta transação?')) {
      try {
        await deleteMutation.mutateAsync(id)
      } catch (error) {
        console.error('Erro ao deletar transação:', error)
      }
    }
  }

  // Obter meses disponíveis para o filtro
  const availableMonths = useMemo(() => {
    if (!transactionsData?.data) return []
    const months = new Set<string>()
    transactionsData.data.forEach((t) => {
      const month = t.reference_month || getMonthFromDate(t.date)
      months.add(month)
    })
    return Array.from(months).sort().reverse()
  }, [transactionsData])

  return (
    <div className="flex flex-col gap-6">
      {/* Header com Botão */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Transações</h1>
          <p className="text-muted-foreground">Gerencie suas receitas e despesas</p>
        </div>
        <Dialog
          open={isDialogOpen}
          onOpenChange={(open) => {
            setIsDialogOpen(open)
            if (!open) {
              setEditingTransaction(null)
            }
          }}
        >
          <DialogTrigger asChild>
            <Button
              onClick={() => setEditingTransaction(null)}
              className="flex items-center gap-2"
            >
              <Plus size={18} />
              Nova Transação
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>
                {editingTransaction ? 'Editar Transação' : 'Nova Transação'}
              </DialogTitle>
              <DialogDescription>
                {editingTransaction
                  ? 'Atualize os dados da transação'
                  : 'Registre uma nova transação financeira'}
              </DialogDescription>
            </DialogHeader>
            <TransactionForm
              transaction={editingTransaction || undefined}
              onSubmit={handleSubmit}
              onCancel={() => {
                setIsDialogOpen(false)
                setEditingTransaction(null)
              }}
              isLoading={createMutation.isPending || updateMutation.isPending}
            />
          </DialogContent>
        </Dialog>
      </div>

      {/* Barra de Filtros */}
      <div className="flex flex-col sm:flex-row gap-4 items-end sm:items-center justify-between">
        <div className="flex items-center gap-2 w-full sm:w-auto overflow-x-auto no-scrollbar">
          <div className="flex items-center gap-2 bg-muted p-1 rounded-lg border">
            <button
              onClick={() => setFilters({ account_id: null })}
              className={`px-3 py-1.5 rounded-md text-xs font-semibold transition-all whitespace-nowrap ${
                !filters.account_id
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground hover:bg-background/50'
              }`}
            >
              Todas as Contas
            </button>
            {accounts?.map((acc) => (
              <button
                key={acc.id}
                onClick={() => setFilters({ account_id: acc.id })}
                className={`px-3 py-1.5 rounded-md text-xs font-semibold transition-all flex items-center gap-2 whitespace-nowrap ${
                  filters.account_id === acc.id
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:bg-background/50'
                }`}
              >
                <div className="w-2 h-2 rounded-full bg-primary"></div>
                {acc.name}
              </button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2 bg-muted p-1 rounded-lg border w-full sm:w-auto">
          <Filter size={14} className="text-muted-foreground ml-2" />
          <Select
            value={filters.reference_month || 'all'}
            onValueChange={(value) =>
              setFilters({ reference_month: value === 'all' ? null : value })
            }
          >
            <SelectTrigger className="bg-transparent border-none outline-none focus:ring-0 w-full sm:w-auto">
              <SelectValue placeholder="Todos os Meses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Todos os Meses</SelectItem>
              {availableMonths.map((month) => {
                const [year, monthNum] = month.split('-')
                const date = new Date(parseInt(year), parseInt(monthNum) - 1, 1)
                return (
                  <SelectItem key={month} value={month}>
                    {format(date, "MMMM 'de' yyyy", { locale: ptBR })}
                  </SelectItem>
                )
              })}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <KpiCard
          title="Saldo Disponível"
          value={globalStats.currentBalance}
          icon={Landmark}
          subtext={
            filters.account_id
              ? accounts?.find((a) => a.id === filters.account_id)?.name || 'Nesta conta'
              : 'Em todas as contas'
          }
          trend={globalStats.currentBalance >= 0 ? 'up' : 'down'}
        />
        <KpiCard
          title="Previsão Faturas"
          value={globalStats.creditCardDebt}
          icon={CreditCard}
          subtext="Total de compras no crédito (Geral)"
          trend="down"
        />
        <KpiCard
          title="Resultado do Período"
          value={globalStats.filteredIncome - globalStats.filteredExpense}
          icon={Target}
          subtext={filters.reference_month ? 'Receitas - Despesas (Mês)' : 'Receitas - Despesas (Visível)'}
          trend={
            globalStats.filteredIncome - globalStats.filteredExpense >= 0 ? 'up' : 'down'
          }
        />
      </div>

      {/* Lista de Transações */}
      {isLoading ? (
        <div className="flex justify-center items-center p-8">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : groupedTransactions.length > 0 && accounts && categories ? (
        <div className="space-y-4">
          {groupedTransactions.map((group, index) => (
            <MonthAccordion
              key={group.monthKey}
              monthKey={group.monthKey}
              transactions={group.transactions}
              accounts={accounts}
              categories={categories}
              onEdit={handleEdit}
              onDelete={handleDelete}
              isOpenDefault={index === 0 || filters.reference_month !== null}
              showAccountTags={!filters.account_id}
            />
          ))}
        </div>
      ) : (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-24">
            <div className="bg-muted w-20 h-20 rounded-full flex items-center justify-center mb-6">
              <PieChart className="text-muted-foreground" size={40} />
            </div>
            <h3 className="text-xl font-bold mb-2">Nada encontrado</h3>
            <p className="text-muted-foreground max-w-md text-center mb-6">
              Não há movimentações para os filtros selecionados de conta ou data.
            </p>
            <Button
              variant="outline"
              onClick={() => setFilters({ account_id: null, reference_month: null })}
            >
              Limpar filtros
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
