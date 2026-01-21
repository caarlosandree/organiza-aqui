'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Plus, Pin, Tag, Edit, Trash2 } from 'lucide-react'
import { NoteForm } from '@/components/knowledge/NoteForm'
import { useNotes } from '@/hooks/queries/useNotes'
import {
  useCreateNote,
  useUpdateNote,
  useDeleteNote,
} from '@/hooks/mutations/useNoteMutations'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { Badge } from '@/components/ui/badge'
import type {
  CreateNoteFormData,
  UpdateNoteFormData,
} from '@/schemas/knowledgeSchema'
import type { Note } from '@/types/knowledge'

export default function NotesPage() {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedNote, setSelectedNote] = useState<Note | undefined>()
  const { data: notes = [], isLoading } = useNotes()
  const createMutation = useCreateNote()
  const updateMutation = useUpdateNote()
  const deleteMutation = useDeleteNote()

  const handleCreate = (data: CreateNoteFormData | UpdateNoteFormData) => {
    if (selectedNote) {
      updateMutation.mutate(
        { id: selectedNote.id, data: data as UpdateNoteFormData },
        {
          onSuccess: () => {
            setIsFormOpen(false)
            setSelectedNote(undefined)
          },
        }
      )
    } else {
      createMutation.mutate(data as CreateNoteFormData, {
        onSuccess: () => {
          setIsFormOpen(false)
        },
      })
    }
  }

  const handleEdit = (note: Note) => {
    setSelectedNote(note)
    setIsFormOpen(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('Tem certeza que deseja deletar esta anotação?')) {
      deleteMutation.mutate(id)
    }
  }

  const handleNew = () => {
    setSelectedNote(undefined)
    setIsFormOpen(true)
  }

  // Separar notas fixadas e não fixadas
  const pinnedNotes = notes.filter((note) => note.is_pinned)
  const unpinnedNotes = notes.filter((note) => !note.is_pinned)

  if (isLoading) {
    return (
      <div className="container mx-auto py-6">
        <p>Carregando anotações...</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Anotações</h1>
          <p className="text-muted-foreground mt-1">
            Organize suas ideias e pensamentos
          </p>
        </div>
        <Button onClick={handleNew}>
          <Plus className="h-4 w-4 mr-2" />
          Nova Anotação
        </Button>
      </div>

      {notes.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          <p>Nenhuma anotação cadastrada</p>
          <p className="text-sm mt-2">
            Clique em &quot;Nova Anotação&quot; para começar
          </p>
        </div>
      ) : (
        <div className="space-y-6">
          {pinnedNotes.length > 0 && (
            <div>
              <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                <Pin className="h-4 w-4" />
                Fixadas
              </h2>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {pinnedNotes.map((note) => (
                  <NoteCard
                    key={note.id}
                    note={note}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                ))}
              </div>
            </div>
          )}

          {unpinnedNotes.length > 0 && (
            <div>
              {pinnedNotes.length > 0 && (
                <h2 className="text-lg font-semibold mb-4">Outras</h2>
              )}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {unpinnedNotes.map((note) => (
                  <NoteCard
                    key={note.id}
                    note={note}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      <NoteForm
        open={isFormOpen}
        onOpenChange={(open) => {
          setIsFormOpen(open)
          if (!open) setSelectedNote(undefined)
        }}
        note={selectedNote}
        onSubmit={handleCreate}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  )
}

function NoteCard({
  note,
  onEdit,
  onDelete,
}: {
  note: {
    id: string
    title: string
    content: string
    tags: string[]
    is_pinned: boolean
    created_at: string
  }
  onEdit: (note: Note) => void
  onDelete: (id: string) => void
}) {
  return (
    <div className="border rounded-lg p-4 hover:bg-accent transition-colors flex flex-col h-full">
      <div className="flex items-start justify-between mb-2">
        <h3 className="font-semibold flex-1">{note.title}</h3>
        {note.is_pinned && (
          <Pin className="h-4 w-4 text-yellow-500 flex-shrink-0" />
        )}
      </div>
      <p className="text-sm text-muted-foreground flex-1 line-clamp-4 mb-3">
        {note.content}
      </p>
      {note.tags.length > 0 && (
        <div className="flex flex-wrap gap-1 mb-3">
          {note.tags.map((tag) => (
            <Badge key={tag} variant="secondary" className="text-xs">
              <Tag className="h-3 w-3 mr-1" />
              {tag}
            </Badge>
          ))}
        </div>
      )}
      <div className="flex items-center justify-between mt-auto pt-3 border-t">
        <span className="text-xs text-muted-foreground">
          {format(new Date(note.created_at), "dd 'de' MMM 'de' yyyy", {
            locale: ptBR,
          })}
        </span>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onEdit(note as Note)}
          >
            <Edit className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDelete(note.id)}
            className="text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
