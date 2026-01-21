'use client'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { formatCurrency } from '@/utils/currency'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import type { ImportPreviewResponse } from '@/types/financial'

interface ImportPreviewProps {
  preview: ImportPreviewResponse
  onImport: () => void
  isLoading?: boolean
}

export function ImportPreview({ preview, onImport, isLoading }: ImportPreviewProps) {
  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Resumo da Importação</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold">{preview.total_transactions}</div>
              <div className="text-sm text-muted-foreground">Total</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">{preview.new_transactions}</div>
              <div className="text-sm text-muted-foreground">Novas</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-yellow-600">{preview.duplicates}</div>
              <div className="text-sm text-muted-foreground">Duplicadas</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {preview.transactions.length > 0 && (
        <div className="space-y-2 max-h-64 overflow-y-auto">
          <h4 className="text-sm font-medium">Transações a serem importadas:</h4>
          {preview.transactions.map((transaction, index) => (
            <div
              key={index}
              className="flex items-center justify-between p-2 border rounded text-sm"
            >
              <div className="flex-1">
                <div className="font-medium">{transaction.description}</div>
                <div className="text-muted-foreground">
                  {format(new Date(transaction.date), "dd 'de' MMMM 'de' yyyy", { locale: ptBR })}
                </div>
              </div>
              <div className="text-right">
                <div className="font-semibold">{formatCurrency(transaction.amount, 'BRL')}</div>
                <Badge variant="outline" className="text-xs">
                  {transaction.type === 'income' ? 'Receita' : transaction.type === 'expense' ? 'Despesa' : 'Transferência'}
                </Badge>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="flex gap-2">
        <Button
          onClick={onImport}
          disabled={isLoading || preview.new_transactions === 0}
          className="flex-1"
        >
          {isLoading ? 'Importando...' : `Importar ${preview.new_transactions} transações`}
        </Button>
      </div>
    </div>
  )
}
