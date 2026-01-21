import { useQuery } from '@tanstack/react-query'
import { getCalendarEvents, getCalendarEvent } from '@/services/knowledgeService'
import type { CalendarEvent, CalendarEventFilters } from '@/types/knowledge'

export const useCalendarEvents = (filters?: CalendarEventFilters) => {
  return useQuery<CalendarEvent[]>({
    queryKey: ['calendar-events', filters],
    queryFn: () => getCalendarEvents(filters),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export const useCalendarEvent = (id: string) => {
  return useQuery<CalendarEvent>({
    queryKey: ['calendar-events', id],
    queryFn: () => getCalendarEvent(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}
