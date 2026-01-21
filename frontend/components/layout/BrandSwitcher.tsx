'use client'

import Link from 'next/link'
import {
  SidebarMenuButton,
  SidebarMenu,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

export function BrandSwitcher() {
  return (
    <SidebarMenu className="group-data-[collapsible=icon]:flex items-center justify-center w-full">
      <SidebarMenuItem className="group-data-[collapsible=icon]:flex justify-center w-full">
        <SidebarMenuButton size="lg" asChild tooltip="Organiza Aqui" className="group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:w-auto">
          <Link href="/" className="flex items-center justify-center gap-3 w-full group group-data-[collapsible=icon]:gap-0">
            <div className="relative flex items-center justify-center w-8 h-8 group-data-[collapsible=icon]:w-10 group-data-[collapsible=icon]:h-10 rounded-lg bg-sidebar-primary/10 group-hover:bg-sidebar-primary/20 transition-colors shrink-0">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className="text-sidebar-primary transition-transform group-hover:scale-110 w-5 h-5 group-data-[collapsible=icon]:w-6 group-data-[collapsible=icon]:h-6"
              >
                <defs>
                  <linearGradient id="logoHexGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stopColor="currentColor" stopOpacity="0.3" />
                    <stop offset="100%" stopColor="currentColor" stopOpacity="0.1" />
                  </linearGradient>
                </defs>
                
                {/* Hexágono externo (Estrutura/Organização) */}
                <path
                  d="M12 2L20.66 7V17L12 22L3.34 17V7L12 2Z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  fill="url(#logoHexGradient)"
                  className="opacity-30"
                />
                
                {/* Conexões internas (Organização/Sistema) */}
                <path
                  d="M12 22V12"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.8"
                />
                <path
                  d="M12 12L20.66 7"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.6"
                />
                <path
                  d="M12 12L3.34 7"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.6"
                />
                
                {/* Conexões adicionais para modernizar */}
                <path
                  d="M12 12L20.66 17"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.5"
                />
                <path
                  d="M12 12L3.34 17"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.5"
                />
                
                {/* Ponto central (Núcleo de organização) */}
                <circle
                  cx="12"
                  cy="12"
                  r="2.5"
                  fill="currentColor"
                  className="drop-shadow-sm"
                />
                
                {/* Pontos nas conexões para modernizar */}
                <circle cx="20.66" cy="7" r="1" fill="currentColor" opacity="0.7" />
                <circle cx="3.34" cy="7" r="1" fill="currentColor" opacity="0.7" />
                <circle cx="20.66" cy="17" r="1" fill="currentColor" opacity="0.7" />
                <circle cx="3.34" cy="17" r="1" fill="currentColor" opacity="0.7" />
              </svg>
            </div>
            
            <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
              <span className="truncate font-semibold text-sidebar-foreground">Organiza Aqui</span>
              <span className="truncate text-xs text-sidebar-foreground/60">Sua vida organizada</span>
            </div>
          </Link>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
