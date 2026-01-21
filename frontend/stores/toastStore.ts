import { create } from 'zustand'
import type { Toast } from '@/components/ui/toast'

interface ToastState {
  toasts: Toast[]
  showToast: (message: string, type?: Toast['type'], duration?: number) => void
  hideToast: (id: string) => void
}

export const useToastStore = create<ToastState>((set) => ({
  toasts: [],
  showToast: (message, type = 'info', duration = 3000) => {
    const id = Math.random().toString(36).substring(7)
    const newToast: Toast = { id, message, type, duration }

    set((state) => ({
      toasts: [...state.toasts, newToast],
    }))

    return id
  },
  hideToast: (id) => {
    set((state) => ({
      toasts: state.toasts.filter((toast) => toast.id !== id),
    }))
  },
}))
