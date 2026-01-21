'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Plus, Check, X, Edit } from 'lucide-react'
import { HabitForm } from '@/components/knowledge/HabitForm'
import { useHabits, useHabitTracking } from '@/hooks/queries/useHabits'
import {
  useCreateHabit,
  useUpdateHabit,
  useDeleteHabit,
  useCreateHabitTracking,
} from '@/hooks/mutations/useHabitMutations'
import { format } from 'date-fns'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type {
  CreateHabitFormData,
  UpdateHabitFormData,
} from '@/schemas/knowledgeSchema'
import type { Habit } from '@/types/knowledge'

export default function HabitsPage() {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedHabit, setSelectedHabit] = useState<Habit | undefined>()
  const { data: habits = [], isLoading } = useHabits()
  const createMutation = useCreateHabit()
  const updateMutation = useUpdateHabit()
  const deleteMutation = useDeleteHabit()

  const handleCreate = (data: CreateHabitFormData | UpdateHabitFormData) => {
    if (selectedHabit) {
      updateMutation.mutate(
        { id: selectedHabit.id, data: data as UpdateHabitFormData },
        {
          onSuccess: () => {
            setIsFormOpen(false)
            setSelectedHabit(undefined)
          },
        }
      )
    } else {
      createMutation.mutate(data as CreateHabitFormData, {
        onSuccess: () => {
          setIsFormOpen(false)
        },
      })
    }
  }

  const handleEdit = (habit: Habit) => {
    setSelectedHabit(habit)
    setIsFormOpen(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('Tem certeza que deseja deletar este hábito?')) {
      deleteMutation.mutate(id)
    }
  }

  const handleNew = () => {
    setSelectedHabit(undefined)
    setIsFormOpen(true)
  }

  if (isLoading) {
    return (
      <div className="container mx-auto py-6">
        <p>Carregando hábitos...</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Hábitos</h1>
          <p className="text-muted-foreground mt-1">
            Acompanhe seus hábitos e mantenha a consistência
          </p>
        </div>
        <Button onClick={handleNew}>
          <Plus className="h-4 w-4 mr-2" />
          Novo Hábito
        </Button>
      </div>

      {habits.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          <p>Nenhum hábito cadastrado</p>
          <p className="text-sm mt-2">
            Clique em &quot;Novo Hábito&quot; para começar
          </p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {habits.map((habit) => (
            <HabitCard
              key={habit.id}
              habit={habit}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      <HabitForm
        open={isFormOpen}
        onOpenChange={(open) => {
          setIsFormOpen(open)
          if (!open) setSelectedHabit(undefined)
        }}
        habit={selectedHabit}
        onSubmit={handleCreate}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  )
}

function HabitCard({
  habit,
  onEdit,
  onDelete,
}: {
  habit: Habit
  onEdit: (habit: Habit) => void
  onDelete: (id: string) => void
}) {
  const today = format(new Date(), 'yyyy-MM-dd')
  const { data: tracking = [] } = useHabitTracking(habit.id)
  const createTrackingMutation = useCreateHabitTracking()

  const todayTracking = tracking.find((t) => t.date === today)

  const handleToggle = () => {
    if (todayTracking) {
      // Se já existe, atualizar
      // Por enquanto, vamos apenas criar/atualizar
      createTrackingMutation.mutate({
        habit_id: habit.id,
        date: today,
        completed: !todayTracking.completed,
      })
    } else {
      // Criar novo
      createTrackingMutation.mutate({
        habit_id: habit.id,
        date: today,
        completed: true,
      })
    }
  }

  const frequencyLabels = {
    daily: 'Diário',
    weekly: 'Semanal',
    monthly: 'Mensal',
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-2 flex-1">
            <div
              className="h-4 w-4 rounded-full"
              style={{ backgroundColor: habit.color }}
            />
            <CardTitle className="text-lg">{habit.name}</CardTitle>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onEdit(habit)}
            >
              <Edit className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onDelete(habit.id)}
              className="text-destructive hover:text-destructive"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>
        {habit.description && (
          <p className="text-sm text-muted-foreground mt-1">
            {habit.description}
          </p>
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-2">
          <Badge variant="secondary">
            {frequencyLabels[habit.frequency as keyof typeof frequencyLabels]}
          </Badge>
          <span className="text-sm text-muted-foreground">
            {habit.target_days} vez{habit.target_days !== 1 ? 'es' : ''} por
            período
          </span>
        </div>

        <div className="flex items-center justify-between pt-4 border-t">
          <span className="text-sm font-medium">Hoje</span>
          <Button
            variant={todayTracking?.completed ? 'default' : 'outline'}
            size="sm"
            onClick={handleToggle}
            disabled={createTrackingMutation.isPending}
          >
            {todayTracking?.completed ? (
              <>
                <Check className="h-4 w-4 mr-1" />
                Concluído
              </>
            ) : (
              'Marcar'
            )}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
