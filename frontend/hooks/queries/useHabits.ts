import { useQuery } from '@tanstack/react-query'
import { getHabits, getHabit, getHabitStats, getHabitTracking } from '@/services/knowledgeService'
import type { Habit, HabitStats, HabitTracking } from '@/types/knowledge'

export const useHabits = () => {
  return useQuery<Habit[]>({
    queryKey: ['habits'],
    queryFn: getHabits,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export const useHabit = (id: string) => {
  return useQuery<Habit>({
    queryKey: ['habits', id],
    queryFn: () => getHabit(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export const useHabitStats = (
  id: string,
  startDate?: string,
  endDate?: string
) => {
  return useQuery<HabitStats>({
    queryKey: ['habits', id, 'stats', startDate, endDate],
    queryFn: () => getHabitStats(id, startDate, endDate),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export const useHabitTracking = (
  habitId: string,
  startDate?: string,
  endDate?: string
) => {
  return useQuery<HabitTracking[]>({
    queryKey: ['habits', habitId, 'tracking', startDate, endDate],
    queryFn: () => getHabitTracking(habitId, startDate, endDate),
    enabled: !!habitId,
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}
