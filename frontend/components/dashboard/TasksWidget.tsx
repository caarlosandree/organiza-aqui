'use client'

import { CheckSquare, Clock, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useTasks } from '@/hooks/queries/useTasks'
import { Skeleton } from '@/components/ui/skeleton'

export const TasksWidget = () => {
  const { data: tasks = [], isLoading } = useTasks()

  const completedTasks = tasks.filter((t) => t.completed_at).length
  const pendingTasks = tasks.filter((t) => !t.completed_at).length
  const overdueTasks = tasks.filter(
    (t) => !t.completed_at && t.due_date && new Date(t.due_date) < new Date()
  ).length

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-24" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Tarefas</CardTitle>
        <CheckSquare className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{tasks.length}</div>
        <p className="text-xs text-muted-foreground mt-1">Total de tarefas</p>
        <div className="flex flex-col gap-2 mt-4">
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-2">
              <CheckSquare className="h-3 w-3 text-[var(--success)]" />
              <span>Concluídas</span>
            </div>
            <span className="font-medium">{completedTasks}</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-2">
              <Clock className="h-3 w-3 text-[var(--warning)]" />
              <span>Pendentes</span>
            </div>
            <span className="font-medium">{pendingTasks}</span>
          </div>
          {overdueTasks > 0 && (
            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <AlertCircle className="h-3 w-3 text-destructive" />
                <span>Vencidas</span>
              </div>
              <span className="font-medium text-destructive">
                {overdueTasks}
              </span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
