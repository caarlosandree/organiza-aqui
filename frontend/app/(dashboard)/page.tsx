'use client'

import { FinancialWidget } from '@/components/dashboard/FinancialWidget'
import { TasksWidget } from '@/components/dashboard/TasksWidget'
import { TimelineWidget } from '@/components/dashboard/TimelineWidget'

export default function DashboardPage() {
  return (
    <div className="container mx-auto py-6 space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-1">
          Visão geral da sua organização pessoal
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <FinancialWidget />
        <TasksWidget />
        <TimelineWidget />
      </div>
    </div>
  )
}
