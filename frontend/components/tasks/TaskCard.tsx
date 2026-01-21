'use client'

import { Check, Clock, Flag } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { getPriorityClasses } from '@/lib/theme'
import type { Task } from '@/types/task'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

interface TaskCardProps {
  task: Task
  onClick?: () => void
}

const priorityColors = {
  low: cn(getPriorityClasses('low').background, getPriorityClasses('low').text),
  medium: cn(getPriorityClasses('medium').background, getPriorityClasses('medium').text),
  high: cn(getPriorityClasses('high').background, getPriorityClasses('high').text),
  urgent: cn(getPriorityClasses('urgent').background, getPriorityClasses('urgent').text),
}

const priorityIcons = {
  low: <Flag className="h-3 w-3" />,
  medium: <Flag className="h-3 w-3" />,
  high: <Flag className="h-3 w-3" />,
  urgent: <Flag className="h-3 w-3" />,
}

export const TaskCard = ({ task, onClick }: TaskCardProps) => {
  const isCompleted = !!task.completed_at
  const isOverdue = task.due_date && !isCompleted && new Date(task.due_date) < new Date()

  return (
    <Card
      className={cn(
        'cursor-pointer transition-all hover:shadow-md',
        isCompleted && 'opacity-60',
        isOverdue && 'border-destructive'
      )}
      onClick={onClick}
    >
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-2">
          <CardTitle
            className={cn(
              'text-base font-medium line-clamp-2',
              isCompleted && 'line-through'
            )}
          >
            {task.title}
          </CardTitle>
          {isCompleted && (
            <Check className="h-5 w-5 text-success shrink-0" />
          )}
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {task.description && (
          <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
            {task.description}
          </p>
        )}
        <div className="flex items-center gap-2 flex-wrap">
          <Badge
            variant="outline"
            className={cn('text-xs', priorityColors[task.priority])}
          >
            {priorityIcons[task.priority]}
            <span className="ml-1 capitalize">{task.priority}</span>
          </Badge>
          {task.due_date && (
            <Badge
              variant="outline"
              className={cn(
                'text-xs',
                isOverdue && 'border-destructive text-destructive'
              )}
            >
              <Clock className="h-3 w-3 mr-1" />
              {format(new Date(task.due_date), "dd 'de' MMM", { locale: ptBR })}
            </Badge>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
