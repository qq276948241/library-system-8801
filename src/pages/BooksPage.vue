<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Search, BookOpen, Library, Filter, BookPlus } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BookCover from '@/components/BookCover.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import { booksApi, borrowsApi, type Book, type CategoryCount } from '@/lib/api'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { cn } from '@/lib/utils'

useTheme()
const { success, error } = useToast()

const books = ref<Book[]>([])
const categories = ref<CategoryCount[]>([])
const loading = ref(false)
const borrowLoading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const keyword = ref('')
const category = ref('')

const selectedBook = ref<Book | null>(null)
const confirmBorrow = ref(false)

const totalBooks = computed(() => categories.value.reduce((sum, c) => sum + c.count, 0))

async function loadCategories() {
  try {
    categories.value = await booksApi.categories()
  } catch (e) {
    // ignore
  }
}

async function loadBooks() {
  loading.value = true
  try {
    const res = await booksApi.list({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value || undefined,
      category: category.value || undefined,
    })
    books.value = res.list
    total.value = res.total
  } catch (e) {
    error(e instanceof Error ? e.message : '加载图书列表失败')
  } finally {
    loading.value = false
  }
}

function onSearch() {
  page.value = 1
  loadBooks()
}

function onCategoryChange() {
  page.value = 1
  loadBooks()
}

function openBorrowDialog(book: Book) {
  if (book.available_copies <= 0) {
    error('该书暂无可用副本')
    return
  }
  selectedBook.value = book
  confirmBorrow.value = true
}

async function onBorrow() {
  if (!selectedBook.value) return
  borrowLoading.value = true
  try {
    await borrowsApi.borrow(selectedBook.value.id)
    success(`成功借阅《${selectedBook.value.title}》`)
    confirmBorrow.value = false
    loadBooks()
  } catch (e) {
    error(e instanceof Error ? e.message : '借阅失败')
  } finally {
    borrowLoading.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadBooks()
})
</script>

<template>
  <div class="animate-fade-up">
    <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="font-serif text-2xl font-bold text-ink-700">图书列表</h1>
        <p class="mt-1 text-sm text-charcoal-muted">
          馆藏共 <span class="font-mono text-charcoal-soft">{{ totalBooks }}</span> 种图书
        </p>
      </div>
      <div class="flex flex-1 flex-col gap-3 sm:flex-row sm:max-w-md">
        <BaseInput
          v-model="keyword"
          placeholder="搜索书名、作者、ISBN..."
          @keyup.enter="onSearch"
        >
          <template #prefix>
            <Search class="h-4 w-4" />
          </template>
        </BaseInput>
        <BaseButton @click="onSearch">
          <template #icon>
            <Search class="h-4 w-4" />
          </template>
          搜索
        </BaseButton>
      </div>
    </div>

    <div class="mb-5 flex flex-wrap gap-2">
      <button
        @click="category = ''; onCategoryChange()"
        :class="[
          'inline-flex items-center gap-1.5 rounded-full px-3.5 py-1.5 text-sm font-medium transition',
          category === ''
            ? 'bg-ink-600 text-paper-50 shadow-card'
            : 'bg-paper-50 text-charcoal-soft hover:bg-ink-50 border border-paper-300',
        ]"
      >
        <Library class="h-3.5 w-3.5" />
        全部
      </button>
      <button
        v-for="c in categories"
        :key="c.category"
        @click="category = c.category; onCategoryChange()"
        :class="[
          'inline-flex items-center gap-1.5 rounded-full px-3.5 py-1.5 text-sm font-medium transition',
          category === c.category
            ? 'bg-ink-600 text-paper-50 shadow-card'
            : 'bg-paper-50 text-charcoal-soft hover:bg-ink-50 border border-paper-300',
        ]"
      >
        <Filter class="h-3.5 w-3.5" />
        {{ c.category }}
        <span class="font-mono text-xs opacity-80">{{ c.count }}</span>
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner size="lg" />
    </div>

    <div v-else-if="books.length === 0">
      <EmptyState
        :icon="BookOpen"
        title="暂无符合条件的图书"
        description="试试其他搜索关键词或分类吧"
      />
    </div>

    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <div
        v-for="book in books"
        :key="book.id"
        class="group overflow-hidden rounded-2xl border border-paper-300 bg-paper-50 shadow-card transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lift"
      >
        <div class="p-3">
          <BookCover
            :color="book.cover_color"
            :title="book.title"
            :author="book.author"
            :category="book.category"
            size="md"
          />
        </div>
        <div class="border-t border-paper-300/60 p-4">
          <h3 class="line-clamp-1 font-serif text-base font-semibold text-ink-700">
            {{ book.title }}
          </h3>
          <p class="mt-0.5 line-clamp-1 text-sm text-charcoal-muted">{{ book.author }}</p>
          <div class="mt-3 flex items-center justify-between">
            <span
              :class="[
                'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium',
                book.available_copies > 0
                  ? 'bg-ink-50 text-ink-700'
                  : 'bg-red-50 text-red-700',
              ]"
            >
              <span
                :class="[
                  'h-1.5 w-1.5 rounded-full',
                  book.available_copies > 0 ? 'bg-ink-500' : 'bg-red-500',
                ]"
              />
              {{ book.available_copies > 0 ? `可借 ${book.available_copies}` : '已借完' }}
            </span>
            <BaseButton
              size="sm"
              :disabled="book.available_copies <= 0"
              @click="openBorrowDialog(book)"
            >
              <template #icon>
                <BookPlus class="h-3.5 w-3.5" />
              </template>
              借阅
            </BaseButton>
          </div>
        </div>
      </div>
    </div>

    <div v-if="books.length > 0" class="mt-8">
      <Pagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="loadBooks"
      />
    </div>

    <ConfirmDialog
      v-model="confirmBorrow"
      title="确认借阅"
      :message="`您确定要借阅《${selectedBook?.title}》吗？借阅期限为 30 天，请按时归还。`"
      confirm-text="确认借阅"
      :loading="borrowLoading"
      @confirm="onBorrow"
    />

    <ToastContainer />
  </div>
</template>
