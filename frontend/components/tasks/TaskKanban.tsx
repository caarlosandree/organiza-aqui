'use client'

import { useState } from 'react'
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { TaskCard } from './TaskCard'
import { TaskForm } from './TaskForm'
import { useTasks } from '@/hooks/queries/useTasks'
import { useTaskStatuses } from '@/hooks/queries/useTaskStatuses'
import { useReorderTask, useCreateTask, useUpdateTask } from '@/hooks/mutations/useTaskMutations'
import type { Task, TaskStatus } from '@/types/task'
import type { CreateTaskFormData, UpdateTaskFormData } from '@/schemas/taskSchema'

interface TaskColumnProps {
  status: TaskStatus
  tasks: Task[]
  onTaskClick: (task: Task) => void
}

const TaskColumn = ({ status, tasks, onTaskClick }: TaskColumnProps) => {
  const { setNodeRef, transform, transition, isDragging } =
    useSortable({ id: status.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div ref={setNodeRef} style={style} className="shrink-0 w-80">
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div
                className="h-3 w-3 rounded-full"
                style={{ backgroundColor: status.color }}
              />
              <CardTitle className="text-base">{status.name}</CardTitle>
              <span className="text-sm text-muted-foreground">({tasks.length})</span>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-3">
          <SortableContext items={tasks.map((t) => t.id)} strategy={verticalListSortingStrategy}>
            {tasks.map((task) => (
              <SortableTaskCard
                key={task.id}
                task={task}
                onClick={() => onTaskClick(task)}
              />
            ))}
          </SortableContext>
        </CardContent>
      </Card>
    </div>
  )
}

interface SortableTaskCardProps {
  task: Task
  onClick: () => void
}

const SortableTaskCard = ({ task, onClick }: SortableTaskCardProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: task.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <TaskCard task={task} onClick={onClick} />
    </div>
  )
}

export const TaskKanban = () => {
  const { data: statuses = [] } = useTaskStatuses()
  const { data: allTasks = [] } = useTasks()
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [activeTask, setActiveTask] = useState<Task | null>(null)

  const reorderMutation = useReorderTask()
  const createMutation = useCreateTask()
  const updateMutation = useUpdateTask()

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  )

  const handleDragStart = (event: DragStartEvent) => {
    const taskId = event.active.id as string
    const task = allTasks.find((t) => t.id === taskId)
    setActiveTask(task || null)
  }

  const handleDragEnd = (event: DragEndEvent) => {
    setActiveTask(null)

    const { active, over } = event
    if (!over) return

    const taskId = active.id as string
    const newStatusId = over.id as string

    // Encontrar a tarefa atual
    const currentTask = allTasks.find((t) => t.id === taskId)
    if (!currentTask) return

    const currentStatusId = currentTask.status_id

    // Se mudou de coluna
    if (currentStatusId !== newStatusId) {
      reorderMutation.mutate({
        task_id: taskId,
        status_id: newStatusId,
      })
    }
  }

  const handleCreateTask = (data: CreateTaskFormData | UpdateTaskFormData) => {
    if (selectedTask) {
      updateMutation.mutate(
        { id: selectedTask.id, data: data as UpdateTaskFormData },
        {
          onSuccess: () => {
            setIsFormOpen(false)
            setSelectedTask(null)
          },
        }
      )
    } else {
      createMutation.mutate(data as CreateTaskFormData, {
        onSuccess: () => {
          setIsFormOpen(false)
        },
      })
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Tarefas</h1>
        <Button onClick={() => setIsFormOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Nova Tarefa
        </Button>
      </div>

      <DndContext sensors={sensors} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
        <div className="flex gap-4 overflow-x-auto pb-4">
          {statuses.map((status: TaskStatus) => {
            const tasks = allTasks.filter((t) => t.status_id === status.id)
            return (
              <TaskColumn
                key={status.id}
                status={status}
                tasks={tasks}
                onTaskClick={(task) => {
                  setSelectedTask(task)
                  setIsFormOpen(true)
                }}
              />
            )
          })}
        </div>
        <DragOverlay>
          {activeTask ? <TaskCard task={activeTask} /> : null}
        </DragOverlay>
      </DndContext>

      <TaskForm
        open={isFormOpen}
        onOpenChange={(open) => {
          setIsFormOpen(open)
          if (!open) setSelectedTask(null)
        }}
        task={selectedTask || undefined}
        onSubmit={handleCreateTask}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  )
}
