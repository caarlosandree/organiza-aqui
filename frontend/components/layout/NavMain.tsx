'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ChevronDown } from 'lucide-react'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

interface NavItem {
  title: string
  url: string
  icon: React.ComponentType<{ className?: string }>
  items?: {
    title: string
    url: string
  }[]
  isActive?: boolean
}

interface NavMainProps {
  items: NavItem[]
}

function NavItemWithCollapsible({ item, pathname }: { item: NavItem; pathname: string }) {
  const { state } = useSidebar()
  const Icon = item.icon
  const isActive = item.isActive ?? (pathname === item.url || pathname?.startsWith(item.url + '/'))
  const isOpenByDefault = pathname.startsWith(item.url + '/')
  const isCollapsed = state === 'collapsed'
  const [open, setOpen] = useState(() => isOpenByDefault)
  const [dropdownOpen, setDropdownOpen] = useState(false)

  // Quando colapsado, usar DropdownMenu ao invés de Collapsible
  if (isCollapsed) {
    return (
      <SidebarMenuItem className="flex justify-center w-full">
        <DropdownMenu open={dropdownOpen} onOpenChange={setDropdownOpen}>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              isActive={isActive}
              size="lg"
              className="w-auto justify-center px-0 [&>svg]:w-6 [&>svg]:h-6"
            >
              <Icon className="w-6 h-6 shrink-0" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            side="right"
            align="start"
            sideOffset={0}
            className="w-48 translate-x-[calc(var(--sidebar-width-icon,4rem)/2)]"
          >
            <DropdownMenuItem disabled className="font-semibold">
              {item.title}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            {item.items?.map((subItem) => {
              const isSubActive = pathname === subItem.url
              return (
                <DropdownMenuItem
                  key={subItem.url}
                  asChild
                  className={isSubActive ? 'bg-accent' : ''}
                >
                  <Link href={subItem.url} onClick={() => setDropdownOpen(false)}>
                    {subItem.title}
                  </Link>
                </DropdownMenuItem>
              )
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    )
  }

  // Quando expandido, usar Collapsible normal
  return (
    <Collapsible asChild defaultOpen={isOpenByDefault} open={open} onOpenChange={setOpen}>
      <SidebarMenuItem className="w-full">
        <CollapsibleTrigger asChild>
          <SidebarMenuButton isActive={isActive} size="lg" className="w-full justify-start pl-4 pr-4">
            <div className="flex items-center gap-3 flex-1">
              <Icon className="w-5 h-5 shrink-0" />
              <span className="text-base font-medium flex-1 text-left">{item.title}</span>
            </div>
            <ChevronDown
              className={`w-4 h-4 shrink-0 transition-transform duration-200 ${open ? 'rotate-180' : ''}`}
            />
          </SidebarMenuButton>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <SidebarMenu>
            {item.items?.map((subItem) => {
              const isSubActive = pathname === subItem.url
              return (
                <SidebarMenuItem key={subItem.url}>
                  <SidebarMenuButton
                    asChild
                    isActive={isSubActive}
                    size="lg"
                    className="w-full justify-start px-8"
                  >
                    <Link href={subItem.url} className="w-full">
                      <span className="text-base font-medium">{subItem.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )
            })}
          </SidebarMenu>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  )
}

export function NavMain({ items }: NavMainProps) {
  const pathname = usePathname()

  return (
    <SidebarGroup className="group-data-[collapsible=icon]:flex group-data-[collapsible=icon]:items-center group-data-[collapsible=icon]:justify-center">
      <SidebarGroupLabel className="group-data-[collapsible=icon]:hidden">Menu</SidebarGroupLabel>
      <SidebarGroupContent className="w-full group-data-[collapsible=icon]:flex group-data-[collapsible=icon]:items-center group-data-[collapsible=icon]:justify-center">
        <SidebarMenu className="w-full group-data-[collapsible=icon]:flex group-data-[collapsible=icon]:items-center">
          {items.map((item) => {
            const Icon = item.icon
            const isActive = item.isActive ?? (pathname === item.url || pathname?.startsWith(item.url + '/'))
            const hasItems = item.items && item.items.length > 0

            if (hasItems) {
              return <NavItemWithCollapsible key={item.url} item={item} pathname={pathname} />
            }

            return (
              <SidebarMenuItem key={item.url} className="group-data-[collapsible=icon]:flex justify-center w-full">
                <SidebarMenuButton asChild tooltip={item.title} isActive={isActive} size="lg" className="w-full justify-start pl-4 pr-4 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:w-auto">
                  <Link href={item.url} className="flex items-center gap-3 group-data-[collapsible=icon]:gap-0">
                    <Icon className="w-5 h-5 group-data-[collapsible=icon]:w-6 group-data-[collapsible=icon]:h-6 shrink-0" />
                    <span className="text-base font-medium group-data-[collapsible=icon]:hidden">{item.title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
