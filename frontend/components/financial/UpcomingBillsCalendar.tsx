'use client'

import { Calendar, Clock, CheckCircle2, XCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatCurrency } from '@/utils/currency'
import { format, isToday, isPast, differenceInDays } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { useUpcomingBills } from '@/hooks/queries/useAnalytics'
import { usePrivacyStore } from '@/stores/privacyStore'

const statusIcons: Record<string, typeof Clock> = {
  pending: Clock,
  paid: CheckCircle2,
  cancelled: XCircle,
  open: Clock,
  partially_paid: Clock,
}

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  paid: 'bg-green-500/10 text-green-600 border-green-500/20',
  cancelled: 'bg-red-500/10 text-red-600 border-red-500/20',
  open: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  partially_paid: 'bg-blue-500/10 text-blue-600 border-blue-500/20',
}

export function UpcomingBillsCalendar({ days = 30 }: { days?: number }) {
  const { data: upcomingBills, isLoading } = useUpcomingBills(days)
  const { isPrivacyMode } = usePrivacyStore()

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Próximos Vencimentos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center p-4">Carregando...</div>
        </CardContent>
      </Card>
    )
  }

  if (!upcomingBills || !upcomingBills.items || upcomingBills.items.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Próximos Vencimentos
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center p-4 text-muted-foreground">
            Nenhum vencimento nos próximos {days} dias
          </div>
        </CardContent>
      </Card>
    )
  }

  const sortedItems = [...upcomingBills.items].sort((a, b) => {
    const dateA = new Date(a.date)
    const dateB = new Date(b.date)
    return dateA.getTime() - dateB.getTime()
  })

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Calendar className="h-5 w-5" />
          Próximos Vencimentos
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {sortedItems.map((item) => {
            const itemDate = new Date(item.date)
            const daysUntilDue = differenceInDays(itemDate, new Date())
            const isOverdue = isPast(itemDate) && !isToday(itemDate)
            const StatusIcon = statusIcons[item.status] || Clock
            const displayAmount = isPrivacyMode ? '••••' : formatCurrency(item.amount, 'BRL')

            return (
              <div
                key={`${item.type}-${item.id}`}
                className="flex items-center justify-between p-3 border rounded-lg hover:bg-accent transition-colors"
              >
                <div className="flex items-center gap-3 flex-1">
                  <StatusIcon className="h-5 w-5 text-muted-foreground" />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="font-medium">{item.description}</span>
                      <Badge className={statusColors[item.status] || ''} variant="outline">
                        {item.status}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-muted-foreground mt-1">
                      <Calendar className="h-3 w-3" />
                      <span>
                        {format(itemDate, "dd 'de' MMMM 'de' yyyy", { locale: ptBR })}
                      </span>
                      {isOverdue && (
                        <Badge variant="destructive" className="ml-2">
                          Vencido
                        </Badge>
                      )}
                      {!isOverdue && daysUntilDue <= 7 && daysUntilDue >= 0 && (
                        <Badge variant="outline" className="ml-2 text-yellow-600">
                          {daysUntilDue === 0 ? 'Vence hoje' : `${daysUntilDue} dias`}
                        </Badge>
                      )}
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <div className={`font-semibold ${isOverdue ? 'text-red-600' : ''}`}>
                    {displayAmount}
                  </div>
                  <div className="text-xs text-muted-foreground capitalize">{item.type}</div>
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}
