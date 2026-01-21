'use client'

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar'
import { User, LogOut } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useLogout } from '@/hooks/mutations/useAuth'
import { cn } from '@/lib/utils'

interface User {
  name?: string
  email?: string
  username?: string
  avatar?: string
}

interface NavUserProps {
  user: User
}

export function NavUser({ user }: NavUserProps) {
  const { isMobile, toggleSidebar, state } = useSidebar()
  const router = useRouter()
  const logoutMutation = useLogout()
  const isCollapsed = state === 'collapsed'

  const handleLogout = () => {
    logoutMutation.mutate()
  }

  const handleProfile = () => {
    router.push('/profile')
    if (isMobile) {
      toggleSidebar()
    }
  }

  const getInitials = (name?: string) => {
    if (!name) return 'U'
    const words = name.trim().split(' ')
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase()
    }
    return name.substring(0, 2).toUpperCase()
  }

  return (
    <SidebarMenu className="group-data-[collapsible=icon]:flex items-center justify-center w-full">
      <SidebarMenuItem className="group-data-[collapsible=icon]:flex justify-center w-full">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              tooltip={user.name || 'Usuário'}
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground w-full justify-start px-4 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:w-auto"
            >
              <Avatar className="h-8 w-8 group-data-[collapsible=icon]:h-10 group-data-[collapsible=icon]:w-10 rounded-full ring-2 ring-transparent group-hover:ring-foreground/20 transition-all">
                <AvatarImage src={user.avatar} alt={user.name || 'Usuário'} />
                <AvatarFallback className="rounded-full">
                  {getInitials(user.name)}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                <span className="truncate font-semibold">{user.name || 'Usuário'}</span>
                <span className="truncate text-xs text-sidebar-foreground/60">
                  {user.email || 'email@example.com'}
                </span>
              </div>
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className={cn(
              'w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg',
              isCollapsed && !isMobile && 'translate-x-[calc(var(--sidebar-width-icon,4rem)/2)]'
            )}
            side={isMobile ? 'bottom' : 'right'}
            align="end"
            sideOffset={0}
          >
            <DropdownMenuItem disabled>
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium leading-none">{user.username || 'username'}</p>
                <p className="text-xs leading-none text-muted-foreground">{user.email || 'email@example.com'}</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleProfile}>
              <User className="h-4 w-4" />
              <span>Meu Perfil</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} variant="destructive">
              <LogOut className="h-4 w-4" />
              <span>Sair</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
