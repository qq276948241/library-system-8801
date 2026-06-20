<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { BookOpen, RotateCcw } from 'lucide-vue-next'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseButton from '@/components/BaseButton.vue'
import BookCover from '@/components/BookCover.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import { borrowsApi, type Borrow, type BorrowStatus } from '@/lib/api'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { formatDate, dueCountdown } from '@/lib/format'

useTheme()
const { success, error } = useToast()

const borrows = ref<Borrow[]>([])
const loading = ref(false)
const returnLoading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const status = ref<BorrowStatus | ''>('')

const selectedBorrow = ref<Borrow | null>(null)
const confirmReturn = ref(false)

const statusOptions = [
  { value: '', label: '全部' },
  { value: 'borrowed', label: '借阅中' },
  { value: 'overdue', label: '已逾期' },
  { value: 'returned', label: '已归还' },
]

async function loadBorrows() {
  loading.value = true
  try {
    const res = await borrowsApi.list({
      page: page.value,
      pageSize: pageSize.value,
      status: status.value || undefined,
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

function openReturnDialog(borrow: Borrow) {
  selectedBorrow.value = borrow
  confirmReturn.value = true
}

async function onReturn() {
  if (!selectedBorrow.value) return
  returnLoading.value = true
  try {
    await borrowsApi.return(selectedBorrow.value.id)
    success(`成功归还《${selectedBorrow.value.book?.title}》`)
    confirmReturn.value = false
    loadBorrows()
  } catch (e) {
    error(e instanceof Error ? e.message : '还书失败')
  } finally {
    returnLoading.value = false
  }
}

onMounted(() => {
  loadBorrows()
})
</script>

<template>
  <div class="animate-fade-up">
    <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="font-serif text-2xl font-bold text-ink-700">我的借阅</h1>
        <p class="mt-1 text-sm text-charcoal-muted">
          共 <span class="font-mono text-charcoal-soft">{{ total }}</span> 条借阅记录
        </p>
      </div>
      <div class="flex flex-1 flex-col gap-3 sm:flex-row sm:max-w-xs">
        <BaseSelect
          v-model="status"
          :options="statusOptions"
          placeholder="选择状态"
          @update:model-value="onStatusChange"
        />
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner size="lg" />
    </div>

    <div v-else-if="borrows.length === 0">
      <EmptyState
        :icon="BookOpen"
        title="暂无借阅记录"
        :description="status ? '当前筛选条件下没有借阅记录' : '您还没有借阅过任何图书'"
      >
        <BaseButton v-if="!status" @click="$router.push('/books')">
          <template #icon>
            <BookOpen class="h-4 w-4" />
          </template>
          去借阅
        </BaseButton>
      </EmptyState>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="borrow in borrows"
        :key="borrow.id"
        class="overflow-hidden rounded-2xl border border-paper-300 bg-paper-50 shadow-card transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lift"
      >
        <div class="flex flex-col sm:flex-row">
          <div class="p-3 sm:w-32">
            <BookCover
              :color="borrow.book?.cover_color"
              :title="borrow.book?.title"
              :author="borrow.book?.author"
              :category="borrow.book?.category"
              size="sm"
            />
          </div>
          <div class="flex flex-1 flex-col justify-between border-t border-paper-300/60 p-4 sm:border-l sm:border-t-0">
            <div>
              <div class="flex flex-wrap items-start justify-between gap-2">
                <div class="flex-1 min-w-0">
                  <h3 class="line-clamp-1 font-serif text-base font-semibold text-ink-700">
                    {{ borrow.book?.title }}
                  </h3>
                  <p class="mt-0.5 line-clamp-1 text-sm text-charcoal-muted">
                    {{ borrow.book?.author }}
                  </p>
                </div>
                <StatusBadge :status="borrow.status" />
              </div>
              <div class="mt-3 grid grid-cols-2 gap-2 text-sm">
                <div>
                  <span class="text-charcoal-muted">借阅日期：</span>
                  <span class="font-mono text-charcoal-soft">{{ formatDate(borrow.borrow_date) }}</span>
                </div>
                <div>
                  <span class="text-charcoal-muted">应还日期：</span>
                  <span class="font-mono text-charcoal-soft">{{ formatDate(borrow.due_date) }}</span>
                </div>
                <div class="col-span-2">
                  <span class="text-charcoal-muted">状态：</span>
                  <span
                    :class="[
                      'font-medium',
                      borrow.status === 'overdue'
                        ? 'text-red-600'
                        : borrow.status === 'returned'
                          ? 'text-brass-600'
                          : 'text-ink-600',
                    ]"
                  >
                    {{ dueCountdown(borrow.due_date, borrow.return_date) }}
                  </span>
                </div>
              </div>
            </div>
            <div v-if="borrow.status !== 'returned'" class="mt-3 flex justify-end">
              <BaseButton size="sm" @click="openReturnDialog(borrow)">
                <template #icon>
                  <RotateCcw class="h-3.5 w-3.5" />
                </template>
                还书
              </BaseButton>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="borrows.length > 0" class="mt-8">
      <Pagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="loadBorrows"
      />
    </div>

    <ConfirmDialog
      v-model="confirmReturn"
      title="确认还书"
      :message="`您确定要归还《${selectedBorrow?.book?.title}》吗？`"
      confirm-text="确认还书"
      :loading="returnLoading"
      @confirm="onReturn"
    />

    <ToastContainer />
  </div>
</template>
