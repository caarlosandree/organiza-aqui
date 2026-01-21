'use client'

import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { Edit, Trash2 } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { formatCurrency } from '@/utils/currency'
import type { Transaction } from '@/types/financial'

interface TransactionTableProps {
  transactions: Transaction[]
  onEdit?: (transaction: Transaction) => void
  onDelete?: (id: string) => void
  getPeriodType?: (transaction: Transaction) => 'bank' | 'credit_card'
}

const typeLabels: Record<Transaction['type'], string> = {
  income: 'Receita',
  expense: 'Despesa',
  transfer: 'Transferência',
}

const typeColors: Record<Transaction['type'], string> = {
  income: 'bg-[var(--success)]/20 text-[var(--success)]',
  expense: 'bg-destructive/20 text-destructive',
  transfer: 'bg-accent/20 text-accent-foreground',
}

export function TransactionTable({ transactions, onEdit, onDelete, getPeriodType }: TransactionTableProps) {
  if (transactions.length === 0) {
    return (
      <div className="text-center p-8 text-muted-foreground">
        <p>Nenhuma transação encontrada.</p>
      </div>
    )
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Data</TableHead>
            <TableHead>Descrição</TableHead>
            <TableHead>Categoria</TableHead>
            <TableHead>Tipo Transação</TableHead>
            {getPeriodType && <TableHead>Tipo</TableHead>}
            <TableHead className="text-right">Valor</TableHead>
            <TableHead className="text-right">Ações</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {transactions.map((transaction) => {
            const amountColor =
              transaction.type === 'income' ? 'text-[var(--success)]' : 'text-destructive'
            const date = new Date(transaction.date)

            return (
              <TableRow key={transaction.id}>
                <TableCell>
                  {format(date, "dd 'de' MMM 'de' yyyy", { locale: ptBR })}
                </TableCell>
                <TableCell className="font-medium">
                  {transaction.description || '-'}
                </TableCell>
                <TableCell>
                  {transaction.category_id ? (
                    <span className="text-sm text-muted-foreground">Categoria</span>
                  ) : (
                    '-'
                  )}
                </TableCell>
                <TableCell>
                  <Badge className={typeColors[transaction.type]}>
                    {typeLabels[transaction.type]}
                  </Badge>
                </TableCell>
                {getPeriodType && (
                  <TableCell>
                    {(() => {
                      const periodType = getPeriodType(transaction)
                      return (
                        <Badge
                          variant="outline"
                          className={
                            periodType === 'bank'
                              ? 'border-blue-500/20 bg-blue-500/10 text-blue-700'
                              : 'border-orange-500/20 bg-orange-500/10 text-orange-700'
                          }
                        >
                          {periodType === 'bank' ? 'Extrato' : 'Fatura'}
                        </Badge>
                      )
                    })()}
                  </TableCell>
                )}
                <TableCell className={`text-right font-medium ${amountColor}`}>
                  {transaction.type === 'income' ? '+' : '-'}
                  {formatCurrency(Math.abs(transaction.amount))}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    {onEdit && (
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onEdit(transaction)}
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                    )}
                    {onDelete && (
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onDelete(transaction.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </div>
  )
}
