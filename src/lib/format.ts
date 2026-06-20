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

/** 计算距到期日的天数（负数表示逾期） */
export function daysUntilDue(dueDate: string | null, returnDate: string | null): number | null {
  if (returnDate) return null
  const d = parseDate(dueDate)
  if (!d) return null
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const end = new Date(d.getFullYear(), d.getMonth(), d.getDate())
  return Math.round((end.getTime() - start.getTime()) / 86400000)
}

/** 到期警告级别：normal 正常 | warning 3天内到期 | danger 已逾期 */
export type DueLevel = 'normal' | 'warning' | 'danger' | 'returned'

export function getDueLevel(dueDate: string | null, returnDate: string | null): DueLevel {
  if (returnDate) return 'returned'
  const diff = daysUntilDue(dueDate, returnDate)
  if (diff === null) return 'normal'
  if (diff < 0) return 'danger'
  if (diff <= 3) return 'warning'
  return 'normal'
}

/** 距到期日的相对描述 */
export function dueCountdown(dueDate: string | null, returnDate: string | null): string {
  if (returnDate) return '已归还'
  const diff = daysUntilDue(dueDate, returnDate)
  if (diff === null) return ''
  if (diff < 0) return `逾期 ${Math.abs(diff)} 天`
  if (diff === 0) return '今日到期'
  if (diff === 1) return '明日到期'
  return `剩余 ${diff} 天`
}

/** 获取到期级别对应的样式类 */
export interface DueStyle {
  borderClass: string
  bgClass: string
  textClass: string
  badgeClass: string
}

export function getDueStyles(level: DueLevel): DueStyle {
  switch (level) {
    case 'danger':
      return {
        borderClass: 'border-red-300',
        bgClass: 'bg-red-50/60',
        textClass: 'text-red-700',
        badgeClass: 'bg-red-100 text-red-700 border-red-200',
      }
    case 'warning':
      return {
        borderClass: 'border-yellow-300',
        bgClass: 'bg-yellow-50/60',
        textClass: 'text-yellow-800',
        badgeClass: 'bg-yellow-100 text-yellow-800 border-yellow-200',
      }
    default:
      return {
        borderClass: 'border-paper-300',
        bgClass: 'bg-paper-50',
        textClass: 'text-ink-600',
        badgeClass: 'bg-ink-50 text-ink-700 border-ink-200',
      }
  }
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
