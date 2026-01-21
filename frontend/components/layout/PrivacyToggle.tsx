'use client'

import { Eye, EyeOff } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { usePrivacyStore } from '@/stores/privacyStore'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

export function PrivacyToggle() {
  const { isPrivacyMode, togglePrivacy } = usePrivacyStore()

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            onClick={togglePrivacy}
            className="h-9 w-9"
          >
            {isPrivacyMode ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
            <span className="sr-only">
              {isPrivacyMode ? 'Desativar modo privacidade' : 'Ativar modo privacidade'}
            </span>
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{isPrivacyMode ? 'Desativar modo privacidade' : 'Ativar modo privacidade'}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
