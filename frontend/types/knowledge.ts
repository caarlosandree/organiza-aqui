export interface CalendarEvent {
  id: string
  user_id: string
  title: string
  description: string
  start_date: string
  end_date?: string
  all_day: boolean
  location?: string
  color: string
  created_at: string
  updated_at: string
}

export interface CreateCalendarEventRequest {
  title: string
  description?: string
  start_date: string
  end_date?: string
  all_day?: boolean
  location?: string
  color?: string
}

export interface UpdateCalendarEventRequest {
  title: string
  description?: string
  start_date: string
  end_date?: string
  all_day?: boolean
  location?: string
  color?: string
}

export interface CalendarEventFilters {
  start_date?: string
  end_date?: string
  limit?: number
  offset?: number
}

export interface Note {
  id: string
  user_id: string
  title: string
  content: string
  tags: string[]
  is_pinned: boolean
  created_at: string
  updated_at: string
}

export interface CreateNoteRequest {
  title: string
  content: string
  tags?: string[]
  is_pinned?: boolean
}

export interface UpdateNoteRequest {
  title: string
  content: string
  tags?: string[]
  is_pinned?: boolean
}

export interface NoteFilters {
  tag?: string
  is_pinned?: boolean
  limit?: number
  offset?: number
}

export interface Habit {
  id: string
  user_id: string
  name: string
  description?: string
  color: string
  frequency: 'daily' | 'weekly' | 'monthly'
  target_days: number
  created_at: string
  updated_at: string
}

export interface CreateHabitRequest {
  name: string
  description?: string
  color?: string
  frequency: 'daily' | 'weekly' | 'monthly'
  target_days: number
}

export interface UpdateHabitRequest {
  name: string
  description?: string
  color?: string
  frequency: 'daily' | 'weekly' | 'monthly'
  target_days: number
}

export interface HabitTracking {
  id: string
  habit_id: string
  date: string
  completed: boolean
  notes?: string
  created_at: string
}

export interface CreateHabitTrackingRequest {
  habit_id: string
  date: string
  completed: boolean
  notes?: string
}

export interface UpdateHabitTrackingRequest {
  completed: boolean
  notes?: string
}

export interface HabitStats {
  habit_id: string
  total_days: number
  completed_days: number
  completion_rate: number
  current_streak: number
  longest_streak: number
}
