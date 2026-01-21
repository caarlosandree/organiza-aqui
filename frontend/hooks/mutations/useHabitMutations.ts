import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createHabit,
  updateHabit,
  deleteHabit,
  createHabitTracking,
  updateHabitTracking,
  deleteHabitTracking,
} from '@/services/knowledgeService'
import type {
  CreateHabitRequest,
  UpdateHabitRequest,
  CreateHabitTrackingRequest,
  UpdateHabitTrackingRequest,
} from '@/types/knowledge'

export const useCreateHabit = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateHabitRequest) => createHabit(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['habits'] })
    },
  })
}

export const useUpdateHabit = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateHabitRequest }) =>
      updateHabit(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['habits'] })
      queryClient.invalidateQueries({ queryKey: ['habits', variables.id] })
      queryClient.invalidateQueries({
        queryKey: ['habits', variables.id, 'stats'],
      })
    },
  })
}

export const useDeleteHabit = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteHabit(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['habits'] })
    },
  })
}

export const useCreateHabitTracking = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateHabitTrackingRequest) =>
      createHabitTracking(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['habits', variables.habit_id, 'tracking'],
      })
      queryClient.invalidateQueries({
        queryKey: ['habits', variables.habit_id, 'stats'],
      })
    },
  })
}

export const useUpdateHabitTracking = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data: UpdateHabitTrackingRequest
    }) => updateHabitTracking(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['habit-tracking'] })
      // Invalida também as queries relacionadas ao hábito
      queryClient.invalidateQueries({ queryKey: ['habits'] })
    },
  })
}

export const useDeleteHabitTracking = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteHabitTracking(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['habit-tracking'] })
      queryClient.invalidateQueries({ queryKey: ['habits'] })
    },
  })
}
