import { useQuery } from '@tanstack/react-query'
import { getNotes, getNote } from '@/services/knowledgeService'
import type { Note, NoteFilters } from '@/types/knowledge'

export const useNotes = (filters?: NoteFilters) => {
  return useQuery<Note[]>({
    queryKey: ['notes', filters],
    queryFn: () => getNotes(filters),
    staleTime: 1 * 60 * 1000, // 1 minuto
  })
}

export const useNote = (id: string) => {
  return useQuery<Note>({
    queryKey: ['notes', id],
    queryFn: () => getNote(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })
}
