import { reactive } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  type: ToastType
  message: string
}

const toasts = reactive<Toast[]>([])
let seq = 0

function push(type: ToastType, message: string, duration = 2800) {
  const id = ++seq
  toasts.push({ id, type, message })
  setTimeout(() => remove(id), duration)
}

function remove(id: number) {
  const idx = toasts.findIndex((t) => t.id === id)
  if (idx >= 0) toasts.splice(idx, 1)
}

export function useToast() {
  return {
    toasts,
    success: (m: string) => push('success', m),
    error: (m: string) => push('error', m),
    info: (m: string) => push('info', m),
    remove,
  }
}
