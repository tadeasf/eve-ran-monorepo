import { useState, useCallback } from 'react'

export interface Toast {
  id: string
  title?: string
  description?: string
  action?: React.ReactNode
  variant?: 'default' | 'destructive'
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([])

  const toast = useCallback(
    ({ title, description, action, variant = 'default' }: Omit<Toast, 'id'>) => {
      const id = Math.random().toString(36).substr(2, 9)
      const newToast: Toast = { id, title, description, action, variant }

      setToasts((prev) => [...prev, newToast])

      // Auto-dismiss after 5 seconds
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id))
      }, 5000)

      return {
        id,
        dismiss: () => setToasts((prev) => prev.filter((t) => t.id !== id)),
        update: (props: Partial<Omit<Toast, 'id'>>) =>
          setToasts((prev) =>
            prev.map((t) => (t.id === id ? { ...t, ...props } : t))
          ),
      }
    },
    []
  )

  return {
    toast,
    toasts,
    dismiss: (id: string) => setToasts((prev) => prev.filter((t) => t.id !== id)),
  }
}
