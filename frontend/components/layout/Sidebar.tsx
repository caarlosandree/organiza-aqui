'use client'

import * as React from 'react'
import {
  LayoutDashboard,
  Wallet,
  CheckSquare,
  Calendar,
  FileText,
  Target,
} from 'lucide-react'
import { NavMain } from '@/components/layout/NavMain'
import { NavUser } from '@/components/layout/NavUser'
import { BrandSwitcher } from '@/components/layout/BrandSwitcher'
import { useAuthStore } from '@/stores/authStore'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  SidebarSeparator,
} from '@/components/ui/sidebar'

const navMainItems = [
  {
    title: 'Dashboard',
    url: '/',
    icon: LayoutDashboard,
    items: [
      {
        title: 'Financeiro',
        url: '/financial',
      },
    ],
  },
  {
    title: 'Financeiro',
    url: '/financial',
    icon: Wallet,
    items: [
      {
        title: 'Contas',
        url: '/financial/accounts',
      },
      {
        title: 'Cartões de Crédito',
        url: '/financial/credit-cards',
      },
      {
        title: 'Importar',
        url: '/financial/import',
      },
      {
        title: 'Transações',
        url: '/financial/transactions',
      },
    ],
  },
  {
    title: 'Tarefas',
    url: '/tasks',
    icon: CheckSquare,
  },
  {
    title: 'Agenda',
    url: '/calendar',
    icon: Calendar,
  },
  {
    title: 'Anotações',
    url: '/notes',
    icon: FileText,
  },
  {
    title: 'Hábitos',
    url: '/habits',
    icon: Target,
  },
]

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuthStore()

  const userData = {
    name: user?.name || 'Usuário',
    email: user?.email || 'email@example.com',
    username: user?.username || 'username',
    avatar: undefined,
  }

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="group-data-[collapsible=icon]:flex items-center justify-center">
        <BrandSwitcher />
      </SidebarHeader>
      <div className="px-2 group-data-[collapsible=icon]:px-2">
        <SidebarSeparator className="group-data-[collapsible=icon]:mx-auto group-data-[collapsible=icon]:w-8" />
      </div>
      <SidebarContent>
        <NavMain items={navMainItems} />
      </SidebarContent>
      <div className="px-2 group-data-[collapsible=icon]:px-2">
        <SidebarSeparator className="group-data-[collapsible=icon]:mx-auto group-data-[collapsible=icon]:w-8" />
      </div>
      <SidebarFooter className="group-data-[collapsible=icon]:flex items-center justify-center">
        <NavUser user={userData} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
