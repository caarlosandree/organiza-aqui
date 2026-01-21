import { useQuery } from '@tanstack/react-query'
import { getTasks, getTask } from '@/services/taskService'
import type { Task, TaskFilters } from '@/types/task'

export const useTasks = (filters?: TaskFilters) => {
  return useQuery<Task[]>({
    queryKey: ['tasks', filters],
    queryFn: () => getTasks(filters),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export const useTask = (id: string | null) => {
  return useQuery<Task>({
    queryKey: ['tasks', id],
    queryFn: () => getTask(id!),
    enabled: !!id,
    staleTime: 1 * 60 * 1000,
  })
}

export const useTasksByStatus = (statusId: string | null) => {
  return useTasks(statusId ? { status_id: statusId } : undefined)
}
