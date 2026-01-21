'use client'

import { useState, useMemo, useRef, useEffect, useCallback } from 'react'
import { createPortal } from 'react-dom'
import { Plus, Trash2, Edit, Wallet, Building2, TrendingUp, ChevronDown, ChevronUp, X, Check, Loader2, Filter, ArrowUpDown, ArrowUp, ArrowDown, ChevronLeft, ChevronRight, DollarSign, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { AccountForm } from '@/components/financial/AccountForm'
import { UpdateBalanceDialog } from '@/components/financial/UpdateBalanceDialog'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useBanks } from '@/hooks/queries/useBanks'
import { useCreateAccount, useUpdateAccount, useDeleteAccount, useRecalculateBalance } from '@/hooks/mutations/useAccountMutations'
import { type AccountFormData } from '@/schemas/financialSchema'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { formatCurrency } from '@/utils/currency'
import { cn } from '@/lib/utils'
import type { Account } from '@/types/financial'
import type { Bank } from '@/types/bank'

const accountTypeLabels: Record<Account['type'], string> = {
  checking: 'Conta Corrente',
  savings: 'Poupança',
  credit: 'Cartão de Crédito',
  investment: 'Investimento',
}

const ITEMS_PER_PAGE = 20
const INITIAL_ITEMS = 5

interface BankFilterComboboxProps {
  label: string
  banks: Bank[]
  isLoading?: boolean
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

function BankFilterCombobox({
  label,
  banks,
  isLoading = false,
  value,
  onChange,
  placeholder = 'Selecione um banco',
}: BankFilterComboboxProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [displayCount, setDisplayCount] = useState(INITIAL_ITEMS)
  const [isOpen, setIsOpen] = useState(false)
  const [position, setPosition] = useState({ top: 0, left: 0, width: 0 })
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLDivElement>(null)
  const loadMoreRef = useRef<HTMLDivElement>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

  // Adicionar opção "Todos os bancos"
  const banksWithAll = useMemo(() => {
    return [{ id: 'all', code: 0, name: 'Todos os bancos', full_name: 'Todos os bancos', ispb: '' } as Bank, ...banks]
  }, [banks])

  // Filtrar bancos baseado no termo de busca (código ou nome)
  const filteredBanks = useMemo(() => {
    if (!searchTerm.trim()) {
      return banksWithAll
    }

    const term = searchTerm.toLowerCase().trim()
    return banksWithAll.filter(
      (bank) =>
        bank.code.toString().includes(term) ||
        bank.name.toLowerCase().includes(term) ||
        bank.full_name.toLowerCase().includes(term)
    )
  }, [banksWithAll, searchTerm])

  // Bancos visíveis (com paginação/lazy loading)
  const visibleBanks = useMemo(() => {
    return filteredBanks.slice(0, displayCount)
  }, [filteredBanks, displayCount])

  const hasMore = displayCount < filteredBanks.length

  // Calcular posição do dropdown quando abrir ou quando a janela mudar
  const updatePosition = useCallback(() => {
    if (containerRef.current) {
      const rect = containerRef.current.getBoundingClientRect()
      setPosition({
        top: rect.bottom + window.scrollY + 4,
        left: rect.left + window.scrollX,
        width: rect.width,
      })
    }
  }, [])

  useEffect(() => {
    if (isOpen) {
      updatePosition()
      
      const handleResize = () => updatePosition()
      const handleScroll = () => updatePosition()
      
      window.addEventListener('resize', handleResize)
      window.addEventListener('scroll', handleScroll, true)
      
      return () => {
        window.removeEventListener('resize', handleResize)
        window.removeEventListener('scroll', handleScroll, true)
      }
    }
  }, [isOpen, updatePosition])

  // Fechar ao clicar fora
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node) &&
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
        setSearchTerm('')
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
      return () => document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  // Carregar mais itens quando chegar no fim da lista
  const loadMore = useCallback(() => {
    if (hasMore) {
      setDisplayCount((prev) => Math.min(prev + ITEMS_PER_PAGE, filteredBanks.length))
    }
  }, [hasMore, filteredBanks.length])

  // IntersectionObserver para lazy loading
  useEffect(() => {
    const loadMoreElement = loadMoreRef.current
    const listElement = listRef.current
    if (!loadMoreElement || !listElement || !isOpen) return

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore) {
          loadMore()
        }
      },
      {
        root: listElement,
        rootMargin: '50px',
        threshold: 0.1,
      }
    )

    observer.observe(loadMoreElement)

    return () => {
      observer.disconnect()
    }
  }, [isOpen, hasMore, loadMore])

  // Encontrar o banco selecionado para exibir no input
  const selectedBank = value && value !== 'all' ? banksWithAll.find((bank) => bank.id === value) : null

  // Obter o valor de exibição do input
  const displayValue = isOpen && searchTerm !== ''
    ? searchTerm
    : value === 'all'
      ? 'Todos os bancos'
      : selectedBank
        ? `${selectedBank.code || '---'} - ${selectedBank.full_name}`
        : placeholder

  const handleSelect = (bankId: string) => {
    onChange(bankId)
    setSearchTerm('')
    setIsOpen(false)
    inputRef.current?.blur()
  }

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    onChange('all')
    setSearchTerm('')
    inputRef.current?.focus()
  }

  return (
    <div className="space-y-2">
      <Label htmlFor="filter-bank">{label}</Label>
      {isLoading ? (
        <div className="flex items-center justify-center p-4 border rounded-md">
          <Loader2 className="h-4 w-4 animate-spin" />
        </div>
      ) : (
        <div ref={containerRef} className="relative w-full">
          <div className="relative">
            <Input
              ref={inputRef}
              id="filter-bank"
              className="w-full pr-16"
              placeholder={placeholder}
              value={displayValue}
              onChange={(e) => {
                const newValue = e.target.value
                setSearchTerm(newValue)
                setDisplayCount(INITIAL_ITEMS) // Reset display count when search changes
                if (!isOpen) {
                  setIsOpen(true)
                }
              }}
              onFocus={() => {
                setIsOpen(true)
                setDisplayCount(INITIAL_ITEMS)
              }}
              autoComplete="off"
            />
            <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
              {value && value !== 'all' && (
                <button
                  type="button"
                  onClick={handleClear}
                  className="p-1 hover:bg-muted rounded-sm"
                >
                  <X className="h-4 w-4 text-muted-foreground" />
                </button>
              )}
              <button
                type="button"
                onClick={() => {
                  setIsOpen(!isOpen)
                  if (!isOpen) {
                    inputRef.current?.focus()
                  }
                }}
                className="p-1 hover:bg-muted rounded-sm"
              >
                <ChevronDown className={cn(
                  "h-4 w-4 text-muted-foreground transition-transform",
                  isOpen && "rotate-180"
                )} />
              </button>
            </div>
          </div>

          {isOpen &&
            typeof window !== 'undefined' &&
            createPortal(
              <div
                ref={dropdownRef}
                className="fixed z-[100] bg-popover border rounded-md shadow-md"
                style={{
                  top: `${position.top}px`,
                  left: `${position.left}px`,
                  width: `${position.width}px`,
                }}
              >
                <div
                  ref={listRef}
                  className="max-h-[200px] overflow-y-auto p-1"
                >
                  {visibleBanks.length > 0 ? (
                    <>
                      {visibleBanks.map((bank) => (
                        <button
                          key={bank.id}
                          type="button"
                          onClick={() => handleSelect(bank.id)}
                          className={cn(
                            "w-full flex items-center gap-2 px-2 py-1.5 text-sm rounded-sm hover:bg-accent hover:text-accent-foreground cursor-pointer text-left",
                            value === bank.id && "bg-accent"
                          )}
                        >
                          {bank.id !== 'all' && (
                            <span className="font-mono text-xs text-muted-foreground min-w-[40px]">
                              {bank.code || '---'}
                            </span>
                          )}
                          <span className="truncate flex-1">{bank.full_name}</span>
                          {value === bank.id && (
                            <Check className="h-4 w-4 shrink-0" />
                          )}
                        </button>
                      ))}
                      {/* Sentinela para carregar mais */}
                      {hasMore && (
                        <div
                          ref={loadMoreRef}
                          className="flex items-center justify-center py-2 text-xs text-muted-foreground"
                        >
                          <Loader2 className="h-3 w-3 animate-spin mr-2" />
                          Carregando mais...
                        </div>
                      )}
                    </>
                  ) : (
                    <div className="py-6 text-center text-sm text-muted-foreground">
                      Nenhum banco encontrado
                    </div>
                  )}
                </div>
              </div>,
              document.body
            )}
        </div>
      )}
    </div>
  )
}

export default function AccountsPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingAccount, setEditingAccount] = useState<string | null>(null)
  const [isFiltersOpen, setIsFiltersOpen] = useState(true)
  const [updatingBalanceAccount, setUpdatingBalanceAccount] = useState<Account | null>(null)
  
  // Filtros
  const [filterName, setFilterName] = useState('')
  const [filterBank, setFilterBank] = useState<string>('all')
  const [filterType, setFilterType] = useState<string>('all')
  
  // Ordenação
  const [sortColumn, setSortColumn] = useState<'name' | 'bank' | 'balance' | null>(null)
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')
  
  // Paginação
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage, setItemsPerPage] = useState(10)

  const { data: accounts, isLoading } = useAccounts()
  const { data: banks, isLoading: isLoadingBanks } = useBanks()
  const createMutation = useCreateAccount()
  const updateMutation = useUpdateAccount()
  const deleteMutation = useDeleteAccount()
  const recalculateBalanceMutation = useRecalculateBalance()

  // Mapeamento de bank_id para nome do banco
  const bankMap = useMemo(() => {
    if (!banks) return new Map<string, string>()
    return new Map(banks.map((bank) => [bank.id, bank.full_name || bank.name]))
  }, [banks])

  // Filtrar e ordenar contas
  const filteredAccounts = useMemo(() => {
    if (!accounts) return []

    let result = accounts.filter((account) => {
      // Filtro por nome
      if (filterName && !account.name.toLowerCase().includes(filterName.toLowerCase())) {
        return false
      }

      // Filtro por banco
      if (filterBank !== 'all' && account.bank_id !== filterBank) {
        return false
      }

      // Filtro por tipo
      if (filterType !== 'all' && account.type !== filterType) {
        return false
      }

      return true
    })

    // Ordenação
    if (sortColumn) {
      result = [...result].sort((a, b) => {
        let aValue: string | number
        let bValue: string | number

        if (sortColumn === 'name') {
          aValue = a.name.toLowerCase()
          bValue = b.name.toLowerCase()
        } else if (sortColumn === 'bank') {
          aValue = (bankMap.get(a.bank_id) || 'N/A').toLowerCase()
          bValue = (bankMap.get(b.bank_id) || 'N/A').toLowerCase()
        } else if (sortColumn === 'balance') {
          aValue = a.balance
          bValue = b.balance
        } else {
          return 0
        }

        if (aValue < bValue) {
          return sortDirection === 'asc' ? -1 : 1
        }
        if (aValue > bValue) {
          return sortDirection === 'asc' ? 1 : -1
        }
        return 0
      })
    }

    return result
  }, [accounts, filterName, filterBank, filterType, sortColumn, sortDirection, bankMap])

  // Calcular paginação com ajuste automático de página
  const totalPages = Math.ceil(filteredAccounts.length / itemsPerPage)
  // Usar useMemo para calcular a página válida sem causar re-renders desnecessários
  const validCurrentPage = useMemo(() => {
    if (totalPages === 0) return 1
    return Math.min(currentPage, totalPages)
  }, [currentPage, totalPages])
  const startIndex = (validCurrentPage - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  const paginatedAccounts = filteredAccounts.slice(startIndex, endIndex)

  // KPIs
  const kpis = useMemo(() => {
    if (!filteredAccounts) {
      return {
        totalBalance: 0,
        accountCount: 0,
        positiveBalance: 0,
        negativeBalance: 0,
      }
    }

    const totalBalance = filteredAccounts.reduce((sum, account) => sum + account.balance, 0)
    const positiveBalance = filteredAccounts.filter((a) => a.balance >= 0).length
    const negativeBalance = filteredAccounts.filter((a) => a.balance < 0).length

    return {
      totalBalance,
      accountCount: filteredAccounts.length,
      positiveBalance,
      negativeBalance,
    }
  }, [filteredAccounts])

  const handleSubmit = async (data: AccountFormData) => {
    try {
      if (editingAccount) {
        await updateMutation.mutateAsync({ id: editingAccount, data })
      } else {
        await createMutation.mutateAsync(data)
      }
      setIsDialogOpen(false)
      setEditingAccount(null)
    } catch (error) {
      console.error('Erro ao salvar conta:', error)
    }
  }

  const handleEdit = (accountId: string) => {
    setEditingAccount(accountId)
    setIsDialogOpen(true)
  }

  const handleDelete = async (accountId: string) => {
    if (confirm('Tem certeza que deseja deletar esta conta?')) {
      try {
        await deleteMutation.mutateAsync(accountId)
      } catch (error) {
        console.error('Erro ao deletar conta:', error)
      }
    }
  }

  const handleClearFilters = () => {
    setFilterName('')
    setFilterBank('all')
    setFilterType('all')
  }

  const handleSort = (column: 'name' | 'bank' | 'balance') => {
    if (sortColumn === column) {
      // Se já está ordenando por essa coluna, inverte a direção
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      // Nova coluna, começa com ascendente
      setSortColumn(column)
      setSortDirection('asc')
    }
  }

  const accountToEdit = editingAccount ? accounts?.find((a) => a.id === editingAccount) : undefined

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Contas</h1>
          <p className="text-muted-foreground">
            Gerencie suas contas financeiras
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={async () => {
              if (!accounts || accounts.length === 0) return
              
              if (confirm('Deseja recalcular o saldo de todas as contas? Isso irá atualizar os saldos baseado no saldo inicial e nas transações.')) {
                try {
                  await Promise.all(
                    accounts.map((account) => recalculateBalanceMutation.mutateAsync(account.id))
                  )
                } catch (error) {
                  console.error('Erro ao recalcular saldos:', error)
                }
              }
            }}
            disabled={recalculateBalanceMutation.isPending || !accounts || accounts.length === 0}
          >
            {recalculateBalanceMutation.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Recalculando...
              </>
            ) : (
              <>
                <RefreshCw className="mr-2 h-4 w-4" />
                Recalcular Saldos
              </>
            )}
          </Button>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setEditingAccount(null)}>
                <Plus className="mr-2 h-4 w-4" />
                Nova Conta
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>
                  {editingAccount ? 'Editar Conta' : 'Nova Conta'}
                </DialogTitle>
                <DialogDescription>
                  {editingAccount
                    ? 'Atualize os dados da conta'
                    : 'Crie uma nova conta financeira'}
                </DialogDescription>
              </DialogHeader>
              <AccountForm
                account={accountToEdit}
                onSubmit={handleSubmit}
                onCancel={() => {
                  setIsDialogOpen(false)
                  setEditingAccount(null)
                }}
                isLoading={createMutation.isPending || updateMutation.isPending}
              />
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* KPIs */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className={cn(
          "border-0",
          kpis.totalBalance >= 0 
            ? "bg-emerald-50/50 dark:bg-emerald-950/30" 
            : "bg-red-50/50 dark:bg-red-950/30"
        )}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo Total</CardTitle>
            <Wallet className={cn(
              "h-4 w-4",
              kpis.totalBalance >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-red-600 dark:text-red-400"
            )} />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${kpis.totalBalance >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}`}>
              {formatCurrency(kpis.totalBalance)}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {kpis.accountCount} conta{kpis.accountCount !== 1 ? 's' : ''}
            </p>
          </CardContent>
        </Card>

        <Card className="border-0 bg-blue-50/50 dark:bg-blue-950/30">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Contas</CardTitle>
            <Building2 className="h-4 w-4 text-blue-600 dark:text-blue-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">{kpis.accountCount}</div>
            <p className="text-xs text-muted-foreground mt-1">
              {kpis.positiveBalance} positiva{kpis.positiveBalance !== 1 ? 's' : ''}, {kpis.negativeBalance} negativa{kpis.negativeBalance !== 1 ? 's' : ''}
            </p>
          </CardContent>
        </Card>

        <Card className="border-0 bg-emerald-50/50 dark:bg-emerald-950/30">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Contas Positivas</CardTitle>
            <TrendingUp className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{kpis.positiveBalance}</div>
            <p className="text-xs text-muted-foreground mt-1">
              Com saldo positivo
            </p>
          </CardContent>
        </Card>

        <Card className="border-0 bg-red-50/50 dark:bg-red-950/30">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Contas Negativas</CardTitle>
            <TrendingUp className="h-4 w-4 text-red-600 dark:text-red-400 rotate-180" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600 dark:text-red-400">{kpis.negativeBalance}</div>
            <p className="text-xs text-muted-foreground mt-1">
              Com saldo negativo
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Tabela */}
      {isLoading ? (
        <div className="flex justify-center items-center p-8">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : (
        <Card>
          <CardContent className="p-6">
            {/* Botão de Filtros */}
            <Collapsible open={isFiltersOpen} onOpenChange={setIsFiltersOpen}>
              <div className="flex items-center justify-between mb-4">
                <CollapsibleTrigger asChild>
                  <Button variant="ghost" size="sm" className="gap-2">
                    <Filter className="h-4 w-4" />
                    Filtros
                    {isFiltersOpen ? (
                      <ChevronUp className="h-4 w-4" />
                    ) : (
                      <ChevronDown className="h-4 w-4" />
                    )}
                  </Button>
                </CollapsibleTrigger>
              </div>

              {/* Filtros */}
              <CollapsibleContent>
                <div className="mb-6 pb-6 border-b">
                  <div className="grid gap-4 md:grid-cols-4">
                <div className="space-y-2">
                  <Label htmlFor="filter-name">Nome da Conta</Label>
                  <Input
                    id="filter-name"
                    placeholder="Buscar por nome..."
                    value={filterName}
                    onChange={(e) => setFilterName(e.target.value)}
                  />
                </div>

                <BankFilterCombobox
                  label="Banco"
                  banks={banks || []}
                  isLoading={isLoadingBanks}
                  value={filterBank}
                  onChange={setFilterBank}
                  placeholder="Todos os bancos"
                />

                <div className="space-y-2">
                  <Label htmlFor="filter-type">Tipo da Conta</Label>
                  <Select value={filterType} onValueChange={setFilterType}>
                    <SelectTrigger id="filter-type" className="!w-full">
                      <SelectValue placeholder="Todos os tipos" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">Todos os tipos</SelectItem>
                      <SelectItem value="checking">Conta Corrente</SelectItem>
                      <SelectItem value="savings">Poupança</SelectItem>
                      <SelectItem value="credit">Cartão de Crédito</SelectItem>
                      <SelectItem value="investment">Investimento</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label className="invisible">Ações</Label>
                  <Tooltip delayDuration={0}>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => {
                          if (filterName !== '' || filterBank !== 'all' || filterType !== 'all') {
                            handleClearFilters()
                          }
                        }}
                        className={cn(
                          "h-9 w-9",
                          filterName === '' && filterBank === 'all' && filterType === 'all' && "opacity-50 cursor-not-allowed"
                        )}
                        aria-label="Limpar filtros"
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent side="top">
                      <p>Limpar filtros</p>
                    </TooltipContent>
                  </Tooltip>
                </div>
                  </div>
                </div>
              </CollapsibleContent>
            </Collapsible>

            {/* Tabela */}
            {filteredAccounts && filteredAccounts.length > 0 ? (
              <>
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow className="bg-muted/50 border-b-2">
                        <TableHead className="py-3 px-4 font-semibold">ID</TableHead>
                        <TableHead className="py-3 px-4 font-semibold">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-0 font-semibold hover:bg-transparent"
                            onClick={() => handleSort('name')}
                          >
                            <span className="flex items-center gap-2">
                              Nome da Conta
                              {sortColumn === 'name' ? (
                                sortDirection === 'asc' ? (
                                  <ArrowUp className="h-4 w-4" />
                                ) : (
                                  <ArrowDown className="h-4 w-4" />
                                )
                              ) : (
                                <ArrowUpDown className="h-4 w-4 opacity-50" />
                              )}
                            </span>
                          </Button>
                        </TableHead>
                        <TableHead className="py-3 px-4 font-semibold">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-0 font-semibold hover:bg-transparent"
                            onClick={() => handleSort('bank')}
                          >
                            <span className="flex items-center gap-2">
                              Banco
                              {sortColumn === 'bank' ? (
                                sortDirection === 'asc' ? (
                                  <ArrowUp className="h-4 w-4" />
                                ) : (
                                  <ArrowDown className="h-4 w-4" />
                                )
                              ) : (
                                <ArrowUpDown className="h-4 w-4 opacity-50" />
                              )}
                            </span>
                          </Button>
                        </TableHead>
                        <TableHead className="py-3 px-4 font-semibold">Tipo da Conta</TableHead>
                        <TableHead className="text-right py-3 px-4 font-semibold">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-0 font-semibold hover:bg-transparent ml-auto"
                            onClick={() => handleSort('balance')}
                          >
                            <span className="flex items-center gap-2">
                              Saldo
                              {sortColumn === 'balance' ? (
                                sortDirection === 'asc' ? (
                                  <ArrowUp className="h-4 w-4" />
                                ) : (
                                  <ArrowDown className="h-4 w-4" />
                                )
                              ) : (
                                <ArrowUpDown className="h-4 w-4 opacity-50" />
                              )}
                            </span>
                          </Button>
                        </TableHead>
                        <TableHead className="text-right w-[120px] py-3 px-4 font-semibold">Ações</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {paginatedAccounts.map((account) => {
                      const bankName = bankMap.get(account.bank_id) || 'N/A'
                      const balanceColor = account.balance >= 0 ? 'text-[var(--success)]' : 'text-destructive'

                      return (
                        <TableRow key={account.id}>
                          <TableCell className="font-mono text-xs">{account.id}</TableCell>
                          <TableCell className="font-medium">{account.name}</TableCell>
                          <TableCell>{bankName}</TableCell>
                          <TableCell>
                            <Badge variant="outline">
                              {accountTypeLabels[account.type]}
                            </Badge>
                          </TableCell>
                          <TableCell className={`text-right font-semibold ${balanceColor}`}>
                            {formatCurrency(account.balance, account.currency)}
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex justify-end gap-2">
                              <Tooltip delayDuration={0}>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => setUpdatingBalanceAccount(account)}
                                  >
                                    <DollarSign className="h-4 w-4" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent side="top">
                                  <p>Atualizar Saldo</p>
                                </TooltipContent>
                              </Tooltip>
                              <Tooltip delayDuration={0}>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handleEdit(account.id)}
                                  >
                                    <Edit className="h-4 w-4" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent side="top">
                                  <p>Editar</p>
                                </TooltipContent>
                              </Tooltip>
                              <Tooltip delayDuration={0}>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handleDelete(account.id)}
                                  >
                                    <Trash2 className="h-4 w-4 text-destructive" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent side="top">
                                  <p>Excluir</p>
                                </TooltipContent>
                              </Tooltip>
                            </div>
                          </TableCell>
                        </TableRow>
                      )
                    })}
                  </TableBody>
                </Table>
              </div>
              
              {/* Controles de Paginação */}
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mt-6 pt-6 border-t">
                {/* Seletor de itens por página */}
                <div className="flex items-center gap-2">
                  <Label htmlFor="items-per-page" className="text-sm text-muted-foreground">
                    Itens por página:
                  </Label>
                  <Select
                    value={itemsPerPage.toString()}
                    onValueChange={(value) => {
                      setItemsPerPage(Number(value))
                      setCurrentPage(1)
                    }}
                  >
                    <SelectTrigger id="items-per-page" className="w-[80px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="10">10</SelectItem>
                      <SelectItem value="25">25</SelectItem>
                      <SelectItem value="50">50</SelectItem>
                      <SelectItem value="100">100</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Informações de paginação */}
                <div className="text-sm text-muted-foreground">
                  Mostrando {startIndex + 1} a {Math.min(endIndex, filteredAccounts.length)} de {filteredAccounts.length} conta{filteredAccounts.length !== 1 ? 's' : ''}
                </div>

                {/* Navegação de páginas */}
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                    disabled={currentPage === 1}
                  >
                    <ChevronLeft className="h-4 w-4" />
                    <span className="sr-only">Página anterior</span>
                  </Button>
                  
                  {/* Números de página */}
                  <div className="flex items-center gap-1">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      let pageNumber: number
                      if (totalPages <= 5) {
                        pageNumber = i + 1
                      } else if (currentPage <= 3) {
                        pageNumber = i + 1
                      } else if (currentPage >= totalPages - 2) {
                        pageNumber = totalPages - 4 + i
                      } else {
                        pageNumber = currentPage - 2 + i
                      }
                      
                      return (
                        <Button
                          key={pageNumber}
                          variant={currentPage === pageNumber ? 'default' : 'outline'}
                          size="sm"
                          className="w-9 h-9 p-0"
                          onClick={() => setCurrentPage(pageNumber)}
                        >
                          {pageNumber}
                        </Button>
                      )
                    })}
                  </div>
                  
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setCurrentPage((prev) => Math.min(totalPages, prev + 1))}
                    disabled={currentPage === totalPages}
                  >
                    <ChevronRight className="h-4 w-4" />
                    <span className="sr-only">Próxima página</span>
                  </Button>
                </div>
              </div>
            </>
            ) : (
              <div className="flex flex-col items-center justify-center p-8 text-center">
                <Wallet className="h-12 w-12 text-muted-foreground mb-4" />
                <p className="text-lg font-medium mb-2">Nenhuma conta encontrada</p>
                <p className="text-sm text-muted-foreground">
                  {accounts && accounts.length > 0
                    ? 'Tente ajustar os filtros para encontrar contas.'
                    : 'Clique em "Nova Conta" para começar.'}
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Dialog de Atualização de Saldo */}
      {updatingBalanceAccount && (
        <UpdateBalanceDialog
          account={updatingBalanceAccount}
          open={!!updatingBalanceAccount}
          onOpenChange={(open) => {
            if (!open) {
              setUpdatingBalanceAccount(null)
            }
          }}
        />
      )}
    </div>
  )
}
