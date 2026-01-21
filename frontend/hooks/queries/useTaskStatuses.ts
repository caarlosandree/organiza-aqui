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
    queryFn: async () => {
      const statuses = await getTaskStatuses()
      const status = statuses.find((s) => s.id === id)
      if (!status) {
        throw new Error(`Task status with id ${id} not found`)
      }
      return status
    },
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}
