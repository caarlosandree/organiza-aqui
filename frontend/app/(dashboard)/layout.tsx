'use client'

import { AppSidebar } from '@/components/layout/Sidebar'
import { Header } from '@/components/layout/Header'
import { AuthGuard } from '@/components/layout/AuthGuard'
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'
import { Separator } from '@/components/ui/separator'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <AuthGuard>
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset>
          <Header />
          <Separator />
          <div className="flex flex-1 flex-col gap-4 p-4">
            {children}
          </div>
        </SidebarInset>
      </SidebarProvider>
    </AuthGuard>
  )
}
