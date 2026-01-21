import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createTask,
  updateTask,
  deleteTask,
  reorderTask,
  completeTask,
  uncompleteTask,
} from '@/services/taskService'
import type { CreateTaskRequest, UpdateTaskRequest, ReorderTaskRequest } from '@/types/task'

export const useCreateTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateTaskRequest) => createTask(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export const useUpdateTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTaskRequest }) =>
      updateTask(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
      queryClient.invalidateQueries({ queryKey: ['tasks', variables.id] })
    },
  })
}

export const useDeleteTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteTask(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export const useReorderTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: ReorderTaskRequest) => reorderTask(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export const useCompleteTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => completeTask(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
      queryClient.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })
}

export const useUncompleteTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => uncompleteTask(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
      queryClient.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })
}
