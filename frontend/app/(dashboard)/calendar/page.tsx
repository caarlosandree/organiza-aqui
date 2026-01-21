'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Plus, Edit, Trash2 } from 'lucide-react'
import { CalendarEventForm } from '@/components/knowledge/CalendarEventForm'
import { useCalendarEvents } from '@/hooks/queries/useCalendarEvents'
import {
  useCreateCalendarEvent,
  useUpdateCalendarEvent,
  useDeleteCalendarEvent,
} from '@/hooks/mutations/useCalendarEventMutations'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import type {
  CreateCalendarEventFormData,
  UpdateCalendarEventFormData,
} from '@/schemas/knowledgeSchema'
import type { CalendarEvent } from '@/types/knowledge'

export default function CalendarPage() {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | undefined>()
  const { data: events = [], isLoading } = useCalendarEvents()
  const createMutation = useCreateCalendarEvent()
  const updateMutation = useUpdateCalendarEvent()
  const deleteMutation = useDeleteCalendarEvent()

  const handleCreate = (data: CreateCalendarEventFormData | UpdateCalendarEventFormData) => {
    if (selectedEvent) {
      updateMutation.mutate(
        { id: selectedEvent.id, data: data as UpdateCalendarEventFormData },
        {
          onSuccess: () => {
            setIsFormOpen(false)
            setSelectedEvent(undefined)
          },
        }
      )
    } else {
      createMutation.mutate(data as CreateCalendarEventFormData, {
        onSuccess: () => {
          setIsFormOpen(false)
        },
      })
    }
  }

  const handleEdit = (event: CalendarEvent) => {
    setSelectedEvent(event)
    setIsFormOpen(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('Tem certeza que deseja deletar este evento?')) {
      deleteMutation.mutate(id)
    }
  }

  const handleNew = () => {
    setSelectedEvent(undefined)
    setIsFormOpen(true)
  }

  if (isLoading) {
    return (
      <div className="container mx-auto py-6">
        <p>Carregando eventos...</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Agenda</h1>
          <p className="text-muted-foreground mt-1">
            Gerencie seus eventos e compromissos
          </p>
        </div>
        <Button onClick={handleNew}>
          <Plus className="h-4 w-4 mr-2" />
          Novo Evento
        </Button>
      </div>

      <div className="grid gap-4">
        {events.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            <p>Nenhum evento cadastrado</p>
            <p className="text-sm mt-2">
              Clique em &quot;Novo Evento&quot; para começar
            </p>
          </div>
        ) : (
          events.map((event) => (
            <div
              key={event.id}
              className="border rounded-lg p-4 hover:bg-accent transition-colors"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <div
                      className="h-4 w-4 rounded-full"
                      style={{ backgroundColor: event.color }}
                    />
                    <h3 className="font-semibold">{event.title}</h3>
                  </div>
                  {event.description && (
                    <p className="text-sm text-muted-foreground mt-1">
                      {event.description}
                    </p>
                  )}
                  <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                    <span>
                      {format(new Date(event.start_date), "dd 'de' MMMM 'de' yyyy 'às' HH:mm", {
                        locale: ptBR,
                      })}
                    </span>
                    {event.location && <span>📍 {event.location}</span>}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleEdit(event)}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleDelete(event.id)}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      <CalendarEventForm
        open={isFormOpen}
        onOpenChange={(open) => {
          setIsFormOpen(open)
          if (!open) setSelectedEvent(undefined)
        }}
        event={selectedEvent}
        onSubmit={handleCreate}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  )
}
