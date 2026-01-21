import { useQuery } from '@tanstack/react-query'
import { getTimelineEvents, getTimelineSummary } from '@/services/timelineService'
import type { TimelineEvent, TimelineSummary, TimelineFilters } from '@/types/timeline'

export const useTimelineEvents = (filters?: TimelineFilters) => {
  return useQuery<TimelineEvent[]>({
    queryKey: ['timeline', filters],
    queryFn: () => getTimelineEvents(filters),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export const useTimelineSummary = () => {
  return useQuery<TimelineSummary>({
    queryKey: ['timeline', 'summary'],
    queryFn: getTimelineSummary,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}
