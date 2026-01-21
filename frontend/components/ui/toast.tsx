'use client'

import * as React from 'react'
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react'
import { cn } from '@/lib/utils'

export type ToastType = 'success' | 'error' | 'info'

export interface Toast {
  id: string
  message: string
  type: ToastType
  duration?: number
}

interface ToastProps {
  toast: Toast
  onClose: (id: string) => void
}

const toastIcons = {
  success: CheckCircle,
  error: AlertCircle,
  info: Info,
}

const toastStyles = {
  success: 'bg-[#1F4E5F] text-white border-[#3FA7A0]',
  error: 'bg-destructive text-destructive-foreground border-destructive',
  info: 'bg-[#1F4E5F] text-white border-[#3FA7A0]',
}

export function ToastComponent({ toast, onClose }: ToastProps) {
  const Icon = toastIcons[toast.type]

  React.useEffect(() => {
    if (toast.duration && toast.duration > 0) {
      const timer = setTimeout(() => {
        onClose(toast.id)
      }, toast.duration)

      return () => clearTimeout(timer)
    }
  }, [toast.id, toast.duration, onClose])

  const [isVisible, setIsVisible] = React.useState(false)

  React.useEffect(() => {
    // Trigger animation on mount
    setIsVisible(true)
  }, [])

  return (
    <div
      className={cn(
        'flex items-center gap-3 px-6 py-4 rounded-xl shadow-2xl border-2 min-w-[300px] max-w-md transform transition-all duration-300',
        isVisible
          ? 'translate-y-0 opacity-100'
          : 'translate-y-24 opacity-0',
        toastStyles[toast.type]
      )}
    >
      <Icon className="w-5 h-5 shrink-0" />
      <span className="flex-1 text-sm font-medium">{toast.message}</span>
      <button
        onClick={() => onClose(toast.id)}
        className="shrink-0 hover:opacity-70 transition-opacity"
        aria-label="Fechar notificação"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  )
}

interface ToastContainerProps {
  toasts: Toast[]
  onClose: (id: string) => void
}

export function ToastContainer({ toasts, onClose }: ToastContainerProps) {
  if (toasts.length === 0) return null

  return (
    <div className="fixed bottom-8 right-8 z-50 flex flex-col gap-3 items-end">
      {toasts.map((toast) => (
        <ToastComponent key={toast.id} toast={toast} onClose={onClose} />
      ))}
    </div>
  )
}
