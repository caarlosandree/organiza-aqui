import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createCalendarEvent,
  updateCalendarEvent,
  deleteCalendarEvent,
} from '@/services/knowledgeService'
import type {
  CreateCalendarEventRequest,
  UpdateCalendarEventRequest,
} from '@/types/knowledge'

export const useCreateCalendarEvent = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateCalendarEventRequest) =>
      createCalendarEvent(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
    },
  })
}

export const useUpdateCalendarEvent = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data: UpdateCalendarEventRequest
    }) => updateCalendarEvent(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
      queryClient.invalidateQueries({
        queryKey: ['calendar-events', variables.id],
      })
    },
  })
}

export const useDeleteCalendarEvent = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteCalendarEvent(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
    },
  })
}
