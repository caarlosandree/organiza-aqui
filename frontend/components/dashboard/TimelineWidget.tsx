'use client'

import { Calendar, Clock } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useTimelineSummary, useTimelineEvents } from '@/hooks/queries/useTimeline'
import { Skeleton } from '@/components/ui/skeleton'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

export const TimelineWidget = () => {
  const { data: summary, isLoading: isLoadingSummary } = useTimelineSummary()
  const { data: recentEvents = [], isLoading: isLoadingEvents } = useTimelineEvents({
    limit: 5,
  })

  if (isLoadingSummary || isLoadingEvents) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-24" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Timeline</CardTitle>
        <Calendar className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{summary?.total_events || 0}</div>
        <p className="text-xs text-muted-foreground mt-1">Total de eventos</p>
        <div className="flex flex-col gap-2 mt-4">
          <div className="flex items-center justify-between text-sm">
            <span>Hoje</span>
            <span className="font-medium">{summary?.today_events || 0}</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span>Próximos</span>
            <span className="font-medium">{summary?.upcoming_events || 0}</span>
          </div>
        </div>
        {recentEvents.length > 0 && (
          <div className="mt-4 pt-4 border-t">
            <p className="text-xs font-medium mb-2">Eventos Recentes</p>
            <div className="space-y-2">
              {recentEvents.slice(0, 3).map((event) => (
                <div key={event.id} className="flex items-start gap-2 text-xs">
                  <Clock className="h-3 w-3 mt-0.5 text-muted-foreground" />
                  <div className="flex-1 min-w-0">
                    <p className="truncate font-medium">{event.title}</p>
                    <p className="text-muted-foreground">
                      {format(new Date(event.event_date), "dd 'de' MMM", {
                        locale: ptBR,
                      })}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
