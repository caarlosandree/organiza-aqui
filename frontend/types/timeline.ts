export interface TimelineEvent {
  id: string
  user_id: string
  entity_type: 'transaction' | 'task' | 'calendar_event' | 'note'
  entity_id: string
  title: string
  description: string
  event_date: string
  metadata?: Record<string, unknown>
  created_at: string
}

export interface TimelineSummary {
  total_events: number
  today_events: number
  upcoming_events: number
  by_type: Record<string, number>
}

export interface TimelineFilters {
  entity_type?: 'transaction' | 'task' | 'calendar_event' | 'note'
  start_date?: string
  end_date?: string
  limit?: number
  offset?: number
}
