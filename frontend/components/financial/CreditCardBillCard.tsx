'use client'

import { Calendar, CheckCircle2, Clock } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { formatCurrency } from '@/utils/currency'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { usePrivacyStore } from '@/stores/privacyStore'
import type { CreditCardBill } from '@/types/financial'

interface CreditCardBillCardProps {
  bill: CreditCardBill
  onPayClick?: () => void
  onCloseClick?: () => void
}

const statusLabels: Record<CreditCardBill['status'], string> = {
  open: 'Aberta',
  closed: 'Fechada',
  paid: 'Paga',
}

const statusIcons: Record<CreditCardBill['status'], typeof Calendar> = {
  open: Clock,
  closed: Calendar,
  paid: CheckCircle2,
}

const statusColors: Record<CreditCardBill['status'], string> = {
  open: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  closed: 'bg-blue-500/10 text-blue-600 border-blue-500/20',
  paid: 'bg-green-500/10 text-green-600 border-green-500/20',
}

export function CreditCardBillCard({ bill, onPayClick, onCloseClick }: CreditCardBillCardProps) {
  const { isPrivacyMode } = usePrivacyStore()
  const StatusIcon = statusIcons[bill.status]
  const displayAmount = isPrivacyMode ? '••••' : formatCurrency(bill.total_amount, 'BRL')

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">
          Fatura {String(bill.month).padStart(2, '0')}/{bill.year}
        </CardTitle>
        <StatusIcon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <Badge className={statusColors[bill.status]} variant="outline">
              {statusLabels[bill.status]}
            </Badge>
            <span className="text-2xl font-bold">{displayAmount}</span>
          </div>
          <div className="space-y-1 text-sm text-muted-foreground">
            <div className="flex items-center justify-between">
              <span>Fechamento:</span>
              <span className="font-medium">
                {format(new Date(bill.closing_date), "dd 'de' MMMM", { locale: ptBR })}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span>Vencimento:</span>
              <span className="font-medium">
                {format(new Date(bill.due_date), "dd 'de' MMMM", { locale: ptBR })}
              </span>
            </div>
          </div>
          {bill.status === 'open' && onCloseClick && (
            <Button variant="outline" size="sm" onClick={onCloseClick} className="w-full">
              Fechar Fatura
            </Button>
          )}
          {bill.status === 'closed' && onPayClick && (
            <Button size="sm" onClick={onPayClick} className="w-full">
              Pagar Fatura
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
