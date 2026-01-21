import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createNote, updateNote, deleteNote } from '@/services/knowledgeService'
import type { CreateNoteRequest, UpdateNoteRequest } from '@/types/knowledge'

export const useCreateNote = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateNoteRequest) => createNote(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes'] })
    },
  })
}

export const useUpdateNote = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateNoteRequest }) =>
      updateNote(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['notes'] })
      queryClient.invalidateQueries({ queryKey: ['notes', variables.id] })
    },
  })
}

export const useDeleteNote = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteNote(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes'] })
    },
  })
}
