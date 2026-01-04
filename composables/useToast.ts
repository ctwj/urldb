/**
 * Toast 通知系统
 * 提供全局通知功能
 */

interface ToastOptions {
  message: string
  type?: 'success' | 'error' | 'warning' | 'info'
  duration?: number
  action?: {
    label: string
    handler: () => void
  }
}

interface Toast extends ToastOptions {
  id: string
  timestamp: number
}

export const useToast = () => {
  const toasts = ref<Toast[]>([])

  const addToast = (options: ToastOptions) => {
    const toast: Toast = {
      id: Math.random().toString(36).substr(2, 9),
      timestamp: Date.now(),
      type: options.type || 'info',
      duration: options.duration || 3000,
      ...options
    }

    toasts.value.push(toast)

    // 自动移除
    if (toast.duration > 0) {
      setTimeout(() => {
        removeToast(toast.id)
      }, toast.duration)
    }

    return toast.id
  }

  const removeToast = (id: string) => {
    const index = toasts.value.findIndex(toast => toast.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  const clearToasts = () => {
    toasts.value = []
  }

  // 便捷方法
  const success = (message: string, options?: Omit<ToastOptions, 'message' | 'type'>) => {
    return addToast({ message, type: 'success', ...options })
  }

  const error = (message: string, options?: Omit<ToastOptions, 'message' | 'type'>) => {
    return addToast({ message, type: 'error', duration: 5000, ...options })
  }

  const warning = (message: string, options?: Omit<ToastOptions, 'message' | 'type'>) => {
    return addToast({ message, type: 'warning', duration: 4000, ...options })
  }

  const info = (message: string, options?: Omit<ToastOptions, 'message' | 'type'>) => {
    return addToast({ message, type: 'info', ...options })
  }

  return {
    toasts: readonly(toasts),
    addToast,
    removeToast,
    clearToasts,
    success,
    error,
    warning,
    info
  }
}