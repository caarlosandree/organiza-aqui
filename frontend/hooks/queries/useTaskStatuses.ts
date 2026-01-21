import { useQuery } from '@tanstack/react-query'
import { getTaskStatuses } from '@/services/taskService'
import type { TaskStatus } from '@/types/task'

export const useTaskStatuses = () => {
  return useQuery<TaskStatus[]>({
    queryKey: ['task-statuses'],
    queryFn: getTaskStatuses,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}

export const useTaskStatus = (id: string | null) => {
  return useQuery<TaskStatus>({
    queryKey: ['task-statuses', id],
    queryFn: () => getTaskStatuses().then((statuses) => statuses.find((s) => s.id === id)!),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}
