'use client'

import { useTheme } from 'next-themes'
import { Moon, Sun, User, LogOut } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { SidebarTrigger } from '@/components/ui/sidebar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar'
import { useAuthStore } from '@/stores/authStore'
import { useLogout } from '@/hooks/mutations/useAuth'
import { PrivacyToggle } from './PrivacyToggle'

export function Header() {
  const { theme, setTheme } = useTheme()
  const { user } = useAuthStore()
  const logoutMutation = useLogout()
  const router = useRouter()

  const handleLogout = () => {
    logoutMutation.mutate()
  }

  const handleProfile = () => {
    router.push('/profile')
  }

  // Função para gerar iniciais do nome do usuário
  const getInitials = (name: string | undefined) => {
    if (!name) return 'U'
    const words = name.trim().split(' ')
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase()
    }
    return name.substring(0, 2).toUpperCase()
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center justify-between px-4 w-full">
        <div className="flex items-center gap-2">
          <SidebarTrigger />
          <p className="text-sm font-medium text-muted-foreground">
            Bem-vindo de volta, <span className="font-semibold text-foreground">{user?.name || 'Usuário'}</span>
          </p>
        </div>
        <div className="flex items-center gap-2 ml-auto">
          <PrivacyToggle />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          >
            <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span className="sr-only">Alternar tema</span>
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="rounded-full p-0 h-auto w-auto hover:bg-transparent">
                <Avatar className="h-8 w-8 cursor-pointer ring-2 ring-transparent hover:ring-foreground/20 transition-all">
                  <AvatarImage src="https://github.com/shadcn.png" alt={user?.name || 'Usuário'} />
                  <AvatarFallback>{getInitials(user?.name)}</AvatarFallback>
                </Avatar>
                <span className="sr-only">Menu do usuário</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem disabled className="font-semibold py-1">
                {user?.username || 'username'}
              </DropdownMenuItem>
              <DropdownMenuItem disabled className="text-muted-foreground py-1 -mt-1">
                {user?.email || 'email@example.com'}
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={handleProfile}>
                <User className="h-4 w-4" />
                Meu Perfil
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleLogout} variant="destructive">
                <LogOut className="h-4 w-4" />
                Sair
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}
