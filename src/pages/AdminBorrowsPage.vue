<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { BookCopy, Users, Clock, AlertTriangle, RotateCcw, TrendingUp, BookOpen, Filter } from 'lucide-vue-next'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseButton from '@/components/BaseButton.vue'
import BookCover from '@/components/BookCover.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import { statsApi, borrowsApi, usersApi, type Borrow, type Stats, type User, type BorrowStatus } from '@/lib/api'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { formatDate, dueCountdown, initials } from '@/lib/format'

useTheme()
const { success, error } = useToast()

const stats = ref<Stats | null>(null)
const borrows = ref<Borrow[]>([])
const users = ref<User[]>([])
const loading = ref(false)
const statsLoading = ref(false)
const usersLoading = ref(false)
const returnLoading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref<BorrowStatus | ''>('')
const userFilter = ref('')

const showReturnConfirm = ref(false)
const returningBorrow = ref<Borrow | null>(null)

const statusOptions: { value: BorrowStatus | ''; label: string }[] = [
  { value: '', label: '全部状态' },
  { value: 'borrowed', label: '借阅中' },
  { value: 'overdue', label: '已逾期' },
  { value: 'returned', label: '已归还' },
]

const userOptions = computed(() => {
  return [
    { value: '', label: '全部用户' },
    ...users.value.map(u => ({
      value: String(u.id),
      label: `${u.name} (${u.username})`,
    })),
  ]
})

const statCards = computed(() => {
  if (!stats.value) return []
  return [
    {
      title: '总藏书量',
      value: stats.value.total_copies,
      icon: BookCopy,
      gradient: 'from-ink-500 to-ink-700',
      iconBg: 'bg-ink-400/20',
    },
    {
      title: '当前借阅',
      value: stats.value.active_borrows,
      icon: TrendingUp,
      gradient: 'from-blue-500 to-blue-700',
      iconBg: 'bg-blue-400/20',
    },
    {
      title: '逾期数量',
      value: stats.value.overdue_borrows,
      icon: AlertTriangle,
      gradient: 'from-red-500 to-red-700',
      iconBg: 'bg-red-400/20',
    },
    {
      title: '总用户数',
      value: stats.value.total_users,
      icon: Users,
      gradient: 'from-brass-500 to-brass-700',
      iconBg: 'bg-brass-400/20',
    },
  ]
})

async function loadStats() {
  statsLoading.value = true
  try {
    stats.value = await statsApi.overview()
  } catch (e) {
    error(e instanceof Error ? e.message : '加载统计数据失败')
  } finally {
    statsLoading.value = false
  }
}

async function loadUsers() {
  usersLoading.value = true
  try {
    const res = await usersApi.list({ pageSize: 1000 })
    users.value = res.list
  } catch (e) {
    error(e instanceof Error ? e.message : '加载用户列表失败')
  } finally {
    usersLoading.value = false
  }
}

async function loadBorrows() {
  loading.value = true
  try {
    const res = await borrowsApi.list({
      page: page.value,
      pageSize: pageSize.value,
      status: statusFilter.value || undefined,
      userId: userFilter.value ? Number(userFilter.value) : undefined,
    })
    borrows.value = res.list
    total.value = res.total
  } catch (e) {
    error(e instanceof Error ? e.message : '加载借阅记录失败')
  } finally {
    loading.value = false
  }
}

function onStatusChange() {
  page.value = 1
  loadBorrows()
}

function onUserChange() {
  page.value = 1
  loadBorrows()
}

function openReturnConfirm(borrow: Borrow) {
  returningBorrow.value = borrow
  showReturnConfirm.value = true
}

async function onReturn() {
  if (!returningBorrow.value) return

  returnLoading.value = true
  try {
    await borrowsApi.return(returningBorrow.value.id)
    success(`《${returningBorrow.value.book?.title}》已成功归还`)
    showReturnConfirm.value = false
    returningBorrow.value = null
    loadBorrows()
    loadStats()
  } catch (e) {
    error(e instanceof Error ? e.message : '还书失败')
  } finally {
    returnLoading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadUsers()
  loadBorrows()
})
</script>

<template>
  <div class="animate-fade-up">
    <div class="mb-6">
      <h1 class="font-serif text-2xl font-bold text-ink-700">借阅管理</h1>
      <p class="mt-1 text-sm text-charcoal-muted">
        管理所有用户的借阅记录，支持还书操作和状态筛选
      </p>
    </div>

    <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="(card, index) in statCards"
        :key="index"
        :class="[
          'relative overflow-hidden rounded-2xl p-5 text-white shadow-card',
          `bg-gradient-to-br ${card.gradient}`,
        ]"
      >
        <div class="relative z-10 flex items-start justify-between">
          <div>
            <p class="text-sm font-medium text-white/80">{{ card.title }}</p>
            <p class="mt-1 font-serif text-3xl font-bold">
              {{ statsLoading ? '—' : card.value.toLocaleString() }}
            </p>
          </div>
          <div :class="['flex h-11 w-11 items-center justify-center rounded-xl', card.iconBg]">
            <component :is="card.icon" class="h-5 w-5" />
          </div>
        </div>
        <div class="absolute -right-6 -bottom-6 h-24 w-24 rounded-full bg-white/10" />
        <div class="absolute -right-2 -bottom-2 h-16 w-16 rounded-full bg-white/10" />
      </div>
    </div>

    <div class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end">
      <BaseSelect
        v-model="statusFilter"
        :options="statusOptions"
        placeholder="全部状态"
        class="sm:max-w-xs"
        @change="onStatusChange"
      >
        <template #prefix>
          <Clock class="h-4 w-4" />
        </template>
      </BaseSelect>
      <BaseSelect
        v-model="userFilter"
        :options="userOptions"
        placeholder="全部用户"
        class="sm:max-w-xs"
        :loading="usersLoading"
        @change="onUserChange"
      >
        <template #prefix>
          <Filter class="h-4 w-4" />
        </template>
      </BaseSelect>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner :size="40" />
    </div>

    <div v-else-if="borrows.length === 0">
      <EmptyState
        :icon="BookOpen"
        title="暂无借阅记录"
        description="没有符合当前筛选条件的借阅记录，试试调整筛选条件"
      />
    </div>

    <div v-else class="overflow-hidden rounded-2xl border border-paper-300 bg-paper-50 shadow-card">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-paper-300 bg-ink-50/50">
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">图书</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">借阅人</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">借阅日期</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">应还日期</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">归还日期</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">状态</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">剩余天数</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="borrow in borrows"
              :key="borrow.id"
              class="border-b border-paper-300/60 transition-colors hover:bg-ink-50/30 last:border-b-0"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-12 shrink-0">
                    <BookCover
                      v-if="borrow.book"
                      :color="borrow.book.cover_color"
                      :title="borrow.book.title"
                      :author="borrow.book.author"
                      size="sm"
                    />
                  </div>
                  <div class="min-w-0">
                    <h3 class="line-clamp-1 font-serif text-sm font-semibold text-ink-700">
                      {{ borrow.book?.title || '—' }}
                    </h3>
                    <p class="mt-0.5 line-clamp-1 text-xs text-charcoal-muted">
                      {{ borrow.book?.author || '—' }}
                    </p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-2.5">
                  <div class="flex h-8 w-8 items-center justify-center rounded-full bg-ink-100 font-serif text-sm font-semibold text-ink-600">
                    {{ initials(borrow.user?.name || '') }}
                  </div>
                  <div class="min-w-0">
                    <p class="line-clamp-1 text-sm font-medium text-ink-700">
                      {{ borrow.user?.name || '—' }}
                    </p>
                    <p class="line-clamp-1 text-xs text-charcoal-muted">
                      {{ borrow.user?.username || '' }}
                    </p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-xs text-charcoal-soft">{{ formatDate(borrow.borrow_date) }}</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-xs text-charcoal-soft">{{ formatDate(borrow.due_date) }}</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-xs text-charcoal-soft">{{ formatDate(borrow.return_date) }}</span>
              </td>
              <td class="px-4 py-3 text-center">
                <StatusBadge :status="borrow.status" />
              </td>
              <td class="px-4 py-3 text-center">
                <span
                  :class="[
                    'font-mono text-xs',
                    borrow.status === 'overdue' ? 'text-red-600' :
                    borrow.status === 'returned' ? 'text-charcoal-muted' :
                    'text-ink-600',
                  ]"
                >
                  {{ dueCountdown(borrow.due_date, borrow.return_date) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-center">
                  <BaseButton
                    v-if="borrow.status === 'borrowed' || borrow.status === 'overdue'"
                    size="sm"
                    variant="ghost"
                    @click="openReturnConfirm(borrow)"
                  >
                    <template #icon>
                      <RotateCcw class="h-3.5 w-3.5" />
                    </template>
                    还书
                  </BaseButton>
                  <span v-else class="text-xs text-charcoal-muted">—</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="borrows.length > 0" class="mt-6">
      <Pagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="loadBorrows"
      />
    </div>

    <ConfirmDialog
      v-model="showReturnConfirm"
      title="确认还书"
      :message="`您确定要将《${returningBorrow?.book?.title}》标记为已归还吗？`"
      confirm-text="确认归还"
      :loading="returnLoading"
      @confirm="onReturn"
    />

    <ToastContainer />
  </div>
</template>
