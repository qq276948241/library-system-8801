import type { BorrowStatus } from '@/lib/api'

const pad = (n: number) => String(n).padStart(2, '0')

export function parseDate(s: string | null): Date | null {
  if (!s) return null
  const d = new Date(s)
  return isNaN(d.getTime()) ? null : d
}

export function formatDate(s: string | null): string {
  const d = parseDate(s)
  if (!d) return '—'
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

export function formatDateTime(s: string | null): string {
  const d = parseDate(s)
  if (!d) return '—'
  return `${formatDate(s)} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/** 距到期日的相对描述 */
export function dueCountdown(dueDate: string | null, returnDate: string | null): string {
  if (returnDate) return '已归还'
  const d = parseDate(dueDate)
  if (!d) return ''
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const end = new Date(d.getFullYear(), d.getMonth(), d.getDate())
  const diff = Math.round((end.getTime() - start.getTime()) / 86400000)
  if (diff < 0) return `逾期 ${Math.abs(diff)} 天`
  if (diff === 0) return '今日到期'
  if (diff === 1) return '明日到期'
  return `剩余 ${diff} 天`
}

export interface StatusMeta {
  label: string
  dot: string
  text: string
  bg: string
}

export function statusMeta(status: BorrowStatus): StatusMeta {
  switch (status) {
    case 'borrowed':
      return {
        label: '借阅中',
        dot: 'bg-ink-500',
        text: 'text-ink-700',
        bg: 'bg-ink-50',
      }
    case 'overdue':
      return {
        label: '已逾期',
        dot: 'bg-red-500',
        text: 'text-red-700',
        bg: 'bg-red-50',
      }
    case 'returned':
      return {
        label: '已归还',
        dot: 'bg-brass-500',
        text: 'text-brass-700',
        bg: 'bg-brass-300/30',
      }
  }
}

export function initials(name: string): string {
  if (!name) return '?'
  return name.charAt(0)
}
