import { toast } from 'sonner'
import type { ToastType } from '@/components/ui/toast'

export function useToast() {
  return {
    showToast: (message: string, type?: ToastType, duration?: number) => {
      const toastId = toast[type || 'info'](message, {
        duration: duration || 3000,
      })
      return toastId.toString()
    },
    hideToast: (id: string) => {
      toast.dismiss(id)
    },
    success: (message: string, duration?: number) => {
      const toastId = toast.success(message, {
        duration: duration || 3000,
      })
      return toastId.toString()
    },
    error: (message: string, duration?: number) => {
      const toastId = toast.error(message, {
        duration: duration || 3000,
      })
      return toastId.toString()
    },
    info: (message: string, duration?: number) => {
      const toastId = toast.info(message, {
        duration: duration || 3000,
      })
      return toastId.toString()
    },
  }
}
