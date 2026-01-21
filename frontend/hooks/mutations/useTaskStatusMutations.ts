import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createTaskStatus,
  updateTaskStatus,
  deleteTaskStatus,
  reorderTaskStatuses,
} from '@/services/taskService'
import type {
  CreateTaskStatusRequest,
  UpdateTaskStatusRequest,
} from '@/types/task'

export const useCreateTaskStatus = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTaskStatusRequest) => createTaskStatus(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['task-statuses'] })
    },
  })
}

export const useUpdateTaskStatus = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTaskStatusRequest }) =>
      updateTaskStatus(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['task-statuses'] })
      queryClient.invalidateQueries({ queryKey: ['task-statuses', variables.id] })
    },
  })
}

export const useDeleteTaskStatus = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteTaskStatus(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['task-statuses'] })
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export const useReorderTaskStatuses = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (statusIds: string[]) => reorderTaskStatuses(statusIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['task-statuses'] })
    },
  })
}
