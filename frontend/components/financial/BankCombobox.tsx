'use client'

import { useState, useMemo, useRef, useCallback, useEffect } from 'react'
import { Controller, Control, FieldValues, Path } from 'react-hook-form'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Loader2, ChevronDown, X, Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Bank } from '@/types/bank'

interface BankComboboxProps<T extends FieldValues> {
  control: Control<T>
  name: Path<T>
  banks?: Bank[]
  isLoading?: boolean
  onValueChange?: (value: string) => void
}

const ITEMS_PER_PAGE = 20
const INITIAL_ITEMS = 5

export function BankCombobox<T extends FieldValues>({
  control,
  name,
  banks = [],
  isLoading = false,
  onValueChange,
}: BankComboboxProps<T>) {
  const [searchTerm, setSearchTerm] = useState('')
  const [displayCount, setDisplayCount] = useState(INITIAL_ITEMS)
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLDivElement>(null)
  const loadMoreRef = useRef<HTMLDivElement>(null)

  // Filtrar bancos baseado no termo de busca (código ou nome)
  const filteredBanks = useMemo(() => {
    if (!searchTerm.trim()) {
      return banks
    }

    const term = searchTerm.toLowerCase().trim()
    return banks.filter(
      (bank) =>
        bank.code.toString().includes(term) ||
        bank.name.toLowerCase().includes(term) ||
        bank.full_name.toLowerCase().includes(term)
    )
  }, [banks, searchTerm])

  // Bancos visíveis (com paginação/lazy loading)
  const visibleBanks = useMemo(() => {
    return filteredBanks.slice(0, displayCount)
  }, [filteredBanks, displayCount])

  const hasMore = displayCount < filteredBanks.length

  // Fechar ao clicar fora
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false)
        setSearchTerm('')
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

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

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        // Encontrar o banco selecionado para exibir no input
        const selectedBank = field.value
          ? banks.find((bank) => bank.id === field.value)
          : null

        // Obter o valor de exibição do input
        const displayValue = isOpen && searchTerm !== ''
          ? searchTerm
          : selectedBank
            ? `${selectedBank.code} - ${selectedBank.full_name}`
            : searchTerm

        const handleSelect = (bank: Bank) => {
          field.onChange(bank.id)
          onValueChange?.(bank.id)
          setSearchTerm('')
          setIsOpen(false)
          inputRef.current?.blur()
        }

        const handleClear = (e: React.MouseEvent) => {
          e.stopPropagation()
          field.onChange(null)
          onValueChange?.('')
          setSearchTerm('')
          inputRef.current?.focus()
        }

        return (
          <div className="space-y-2">
            <Label htmlFor={name}>Banco</Label>
            {isLoading ? (
              <div className="flex items-center justify-center p-4 border rounded-md">
                <Loader2 className="h-4 w-4 animate-spin" />
              </div>
            ) : (
              <div ref={containerRef} className="relative w-full">
                <div className="relative">
                  <Input
                    ref={inputRef}
                    id={name}
                    className={cn(
                      fieldState.error ? 'border-destructive' : '',
                      'w-full pr-16'
                    )}
                    placeholder="Buscar por código ou nome do banco..."
                    value={displayValue}
                    onChange={(e) => {
                      const newValue = e.target.value
                      setSearchTerm(newValue)
                      setDisplayCount(INITIAL_ITEMS)
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
                    {field.value && (
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

                {isOpen && (
                  <div className="absolute z-50 w-full mt-1 bg-popover border rounded-md shadow-md">
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
                              onClick={() => handleSelect(bank)}
                              className={cn(
                                "w-full flex items-center gap-2 px-2 py-1.5 text-sm rounded-sm hover:bg-accent hover:text-accent-foreground cursor-pointer text-left",
                                field.value === bank.id && "bg-accent"
                              )}
                            >
                              <span className="font-mono text-xs text-muted-foreground min-w-[40px]">
                                {bank.code || '---'}
                              </span>
                              <span className="truncate flex-1">{bank.full_name}</span>
                              {field.value === bank.id && (
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
                  </div>
                )}
              </div>
            )}
            {fieldState.error && (
              <p className="text-sm text-destructive">{fieldState.error.message}</p>
            )}
          </div>
        )
      }}
    />
  )
}
