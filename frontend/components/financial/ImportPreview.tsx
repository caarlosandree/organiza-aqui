'use client'

import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { formatCurrency } from '@/utils/currency'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { CheckCircle2, XCircle, AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { ImportPreviewResponse } from '@/types/financial'

interface ImportPreviewProps {
  preview: ImportPreviewResponse
  onImport: (externalIds: string[]) => void
  isLoading?: boolean
}

export function ImportPreview({ preview, onImport, isLoading }: ImportPreviewProps) {
  // Filtrar apenas transações novas
  const newTransactions = preview.transactions.filter(t => t.status === 'new')
  const existingTransactions = preview.transactions.filter(t => t.status === 'existing')
  
  // IDs válidos (apenas de transações novas)
  const validExternalIds = new Set(newTransactions.map(t => t.external_id))
  
  // Estado inicial: apenas transações novas selecionadas
  const [selectedExternalIds, setSelectedExternalIds] = useState<Set<string>>(() => 
    new Set(newTransactions.map(t => t.external_id))
  )
  
  // Filtrar selectedExternalIds para garantir apenas IDs válidos (de transações novas)
  // Isso previne que transações existentes sejam importadas
  const filteredSelectedIds = new Set(
    Array.from(selectedExternalIds).filter(id => validExternalIds.has(id))
  )

  const handleToggleTransaction = (externalId: string, isNew: boolean) => {
    // Apenas permitir selecionar transações novas
    if (!isNew) {
      return
    }
    
    const newSet = new Set(selectedExternalIds)
    if (newSet.has(externalId)) {
      newSet.delete(externalId)
    } else {
      newSet.add(externalId)
    }
    setSelectedExternalIds(newSet)
  }

  const handleSelectAll = () => {
    if (filteredSelectedIds.size === newTransactions.length) {
      setSelectedExternalIds(new Set())
    } else {
      setSelectedExternalIds(new Set(newTransactions.map(t => t.external_id)))
    }
  }

  const handleImport = () => {
    // Filtrar apenas external_ids de transações novas para garantir que não importamos duplicatas
    // Usar filteredSelectedIds que já garante apenas transações novas
    const newExternalIds = Array.from(filteredSelectedIds)
    
    if (newExternalIds.length === 0) {
      return // Não há transações para importar
    }
    
    onImport(newExternalIds)
  }

  return (
    <div className="space-y-4">
      {preview.file_hash && (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Hash do arquivo: <code className="text-xs">{preview.file_hash}</code>
          </AlertDescription>
        </Alert>
      )}

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
              <div className="text-sm text-muted-foreground">Já Existem</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {newTransactions.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg">Transações Novas ({newTransactions.length})</CardTitle>
              <Button
                variant="outline"
                size="sm"
                onClick={handleSelectAll}
              >
                {filteredSelectedIds.size === newTransactions.length ? 'Desmarcar Todas' : 'Selecionar Todas'}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {newTransactions.map((transaction) => {
                const isSelected = filteredSelectedIds.has(transaction.external_id)
                return (
                  <div
                    key={transaction.external_id}
                    className={cn(
                      'flex items-center gap-3 p-2 border rounded text-sm cursor-pointer hover:bg-muted/50',
                      isSelected && 'bg-muted'
                    )}
                    onClick={() => handleToggleTransaction(transaction.external_id, true)}
                  >
                    <Checkbox
                      checked={isSelected}
                      onCheckedChange={() => handleToggleTransaction(transaction.external_id, true)}
                      onClick={(e) => e.stopPropagation()}
                    />
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
                    <CheckCircle2 className="h-4 w-4 text-green-600" />
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {existingTransactions.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Transações Já Existentes ({existingTransactions.length})</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {existingTransactions.map((transaction) => (
                <div
                  key={transaction.external_id}
                  className="flex items-center gap-3 p-2 border rounded text-sm opacity-60"
                >
                  <XCircle className="h-4 w-4 text-yellow-600" />
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
          </CardContent>
        </Card>
      )}

      <div className="flex gap-2">
        <Button
          onClick={handleImport}
          disabled={isLoading || filteredSelectedIds.size === 0}
          className="flex-1"
        >
          {isLoading ? 'Importando...' : `Importar ${filteredSelectedIds.size} transação(ões) selecionada(s)`}
        </Button>
      </div>
    </div>
  )
}
