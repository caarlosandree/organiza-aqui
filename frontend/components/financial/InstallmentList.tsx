'use client'

import { Calendar, CheckCircle2, Clock, XCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { formatCurrency } from '@/utils/currency'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { useInstallments } from '@/hooks/queries/useInstallments'
import { usePrivacyStore } from '@/stores/privacyStore'
import type { Transaction } from '@/types/financial'

interface InstallmentListProps {
  parentId: string
  onEditClick?: (transactionId: string) => void
}

const statusLabels: Record<Transaction['status'], string> = {
  pending: 'Pendente',
  paid: 'Paga',
  cancelled: 'Cancelada',
}

const statusIcons: Record<Transaction['status'], typeof Calendar> = {
  pending: Clock,
  paid: CheckCircle2,
  cancelled: XCircle,
}

const statusColors: Record<Transaction['status'], string> = {
  pending: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  paid: 'bg-green-500/10 text-green-600 border-green-500/20',
  cancelled: 'bg-red-500/10 text-red-600 border-red-500/20',
}

export function InstallmentList({ parentId, onEditClick }: InstallmentListProps) {
  const { data: installments, isLoading } = useInstallments(parentId)
  const { isPrivacyMode } = usePrivacyStore()

  if (isLoading) {
    return <div className="text-center p-4">Carregando parcelas...</div>
  }

  if (!installments || installments.length === 0) {
    return <div className="text-center p-4 text-muted-foreground">Nenhuma parcela encontrada</div>
  }

  return (
    <div className="space-y-3">
      {installments.map((installment) => {
        const StatusIcon = statusIcons[installment.status]
        const displayAmount = isPrivacyMode ? '••••' : formatCurrency(installment.amount, 'BRL')

        return (
          <Card key={installment.id}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Parcela {installment.installment_number}/{installment.total_installments}
              </CardTitle>
              <StatusIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="flex flex-col gap-2">
                <div className="flex items-center justify-between">
                  <Badge className={statusColors[installment.status]} variant="outline">
                    {statusLabels[installment.status]}
                  </Badge>
                  <span className="text-lg font-semibold">{displayAmount}</span>
                </div>
                <p className="text-sm text-muted-foreground">{installment.description}</p>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Calendar className="h-4 w-4" />
                  <span>
                    {format(new Date(installment.date), "dd 'de' MMMM 'de' yyyy", { locale: ptBR })}
                  </span>
                </div>
                {onEditClick && installment.status === 'pending' && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onEditClick(installment.id)}
                    className="w-full mt-2"
                  >
                    Editar Parcela
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
