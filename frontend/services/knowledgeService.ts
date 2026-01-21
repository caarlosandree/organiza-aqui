import api from '@/lib/axios'
import type {
  CalendarEvent,
  CreateCalendarEventRequest,
  UpdateCalendarEventRequest,
  CalendarEventFilters,
  Note,
  CreateNoteRequest,
  UpdateNoteRequest,
  NoteFilters,
  Habit,
  CreateHabitRequest,
  UpdateHabitRequest,
  HabitTracking,
  CreateHabitTrackingRequest,
  UpdateHabitTrackingRequest,
  HabitStats,
} from '@/types/knowledge'

// Calendar Events
export const getCalendarEvents = async (
  filters?: CalendarEventFilters
): Promise<CalendarEvent[]> => {
  const params = new URLSearchParams()
  if (filters?.start_date) params.append('start_date', filters.start_date)
  if (filters?.end_date) params.append('end_date', filters.end_date)
  if (filters?.limit) params.append('limit', String(filters.limit))
  if (filters?.offset) params.append('offset', String(filters.offset))

  const response = await api.get<CalendarEvent[]>(
    `/calendar-events?${params.toString()}`
  )
  return response.data
}

export const getCalendarEvent = async (id: string): Promise<CalendarEvent> => {
  const response = await api.get<CalendarEvent>(`/calendar-events/${id}`)
  return response.data
}

export const createCalendarEvent = async (
  data: CreateCalendarEventRequest
): Promise<CalendarEvent> => {
  const response = await api.post<CalendarEvent>('/calendar-events', data)
  return response.data
}

export const updateCalendarEvent = async (
  id: string,
  data: UpdateCalendarEventRequest
): Promise<CalendarEvent> => {
  const response = await api.put<CalendarEvent>(`/calendar-events/${id}`, data)
  return response.data
}

export const deleteCalendarEvent = async (id: string): Promise<void> => {
  await api.delete(`/calendar-events/${id}`)
}

// Notes
export const getNotes = async (filters?: NoteFilters): Promise<Note[]> => {
  const params = new URLSearchParams()
  if (filters?.tag) params.append('tag', filters.tag)
  if (filters?.is_pinned !== undefined)
    params.append('is_pinned', String(filters.is_pinned))
  if (filters?.limit) params.append('limit', String(filters.limit))
  if (filters?.offset) params.append('offset', String(filters.offset))

  const response = await api.get<Note[]>(`/notes?${params.toString()}`)
  return response.data
}

export const getNote = async (id: string): Promise<Note> => {
  const response = await api.get<Note>(`/notes/${id}`)
  return response.data
}

export const createNote = async (data: CreateNoteRequest): Promise<Note> => {
  const response = await api.post<Note>('/notes', data)
  return response.data
}

export const updateNote = async (
  id: string,
  data: UpdateNoteRequest
): Promise<Note> => {
  const response = await api.put<Note>(`/notes/${id}`, data)
  return response.data
}

export const deleteNote = async (id: string): Promise<void> => {
  await api.delete(`/notes/${id}`)
}

// Habits
export const getHabits = async (): Promise<Habit[]> => {
  const response = await api.get<Habit[]>('/habits')
  return response.data
}

export const getHabit = async (id: string): Promise<Habit> => {
  const response = await api.get<Habit>(`/habits/${id}`)
  return response.data
}

export const createHabit = async (data: CreateHabitRequest): Promise<Habit> => {
  const response = await api.post<Habit>('/habits', data)
  return response.data
}

export const updateHabit = async (
  id: string,
  data: UpdateHabitRequest
): Promise<Habit> => {
  const response = await api.put<Habit>(`/habits/${id}`, data)
  return response.data
}

export const deleteHabit = async (id: string): Promise<void> => {
  await api.delete(`/habits/${id}`)
}

export const getHabitStats = async (
  id: string,
  startDate?: string,
  endDate?: string
): Promise<HabitStats> => {
  const params = new URLSearchParams()
  if (startDate) params.append('start_date', startDate)
  if (endDate) params.append('end_date', endDate)

  const response = await api.get<HabitStats>(
    `/habits/${id}/stats?${params.toString()}`
  )
  return response.data
}

// Habit Tracking
export const getHabitTracking = async (
  habitId: string,
  startDate?: string,
  endDate?: string
): Promise<HabitTracking[]> => {
  const params = new URLSearchParams()
  if (startDate) params.append('start_date', startDate)
  if (endDate) params.append('end_date', endDate)

  const response = await api.get<HabitTracking[]>(
    `/habits/${habitId}/tracking?${params.toString()}`
  )
  return response.data
}

export const createHabitTracking = async (
  data: CreateHabitTrackingRequest
): Promise<HabitTracking> => {
  const response = await api.post<HabitTracking>('/habit-tracking', data)
  return response.data
}

export const updateHabitTracking = async (
  id: string,
  data: UpdateHabitTrackingRequest
): Promise<HabitTracking> => {
  const response = await api.put<HabitTracking>(`/habit-tracking/${id}`, data)
  return response.data
}

export const deleteHabitTracking = async (id: string): Promise<void> => {
  await api.delete(`/habit-tracking/${id}`)
}
