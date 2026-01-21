import api from '@/lib/axios'
import type { TimelineEvent, TimelineSummary, TimelineFilters } from '@/types/timeline'

export const getTimelineEvents = async (
  filters?: TimelineFilters
): Promise<TimelineEvent[]> => {
  const params = new URLSearchParams()
  if (filters?.entity_type) params.append('entity_type', filters.entity_type)
  if (filters?.start_date) params.append('start_date', filters.start_date)
  if (filters?.end_date) params.append('end_date', filters.end_date)
  if (filters?.limit) params.append('limit', String(filters.limit))
  if (filters?.offset) params.append('offset', String(filters.offset))

  const response = await api.get<TimelineEvent[]>(`/timeline?${params.toString()}`)
  return response.data
}

export const getTimelineSummary = async (): Promise<TimelineSummary> => {
  const response = await api.get<TimelineSummary>('/timeline/summary')
  return response.data
}
