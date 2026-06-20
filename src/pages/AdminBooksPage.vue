<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { Search, Plus, Edit, Trash2, BookOpen, Filter, Save, X } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTextarea from '@/components/BaseTextarea.vue'
import BookCover from '@/components/BookCover.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import { booksApi, type Book, type CategoryCount, type BookInput } from '@/lib/api'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { formatDate } from '@/lib/format'
import { cn } from '@/lib/utils'

useTheme()
const { success, error } = useToast()

const coverColorOptions = [
  { value: '#1F3D2B', label: '墨绿' },
  { value: '#2D3A4A', label: '藏蓝' },
  { value: '#5C2E2E', label: '暗红' },
  { value: '#4A3F2D', label: '棕褐' },
  { value: '#3D2B4A', label: '暗紫' },
  { value: '#2B4A4A', label: '青色' },
  { value: '#4A4A2B', label: '橄榄' },
  { value: '#3B3B3B', label: '炭黑' },
]

interface FormErrors {
  title?: string
  author?: string
  isbn?: string
  category?: string
  cover_color?: string
  publisher?: string
  published_year?: string
  total_copies?: string
}

const books = ref<Book[]>([])
const categories = ref<CategoryCount[]>([])
const loading = ref(false)
const formLoading = ref(false)
const deleteLoading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref('')
const categoryFilter = ref('')

const showModal = ref(false)
const isEdit = ref(false)
const editingBook = ref<Book | null>(null)

const form = reactive<BookInput>({
  title: '',
  author: '',
  isbn: '',
  category: '',
  description: '',
  cover_color: '#1F3D2B',
  publisher: '',
  published_year: undefined,
  total_copies: 1,
})

const formErrors = reactive<FormErrors>({})

const showDeleteConfirm = ref(false)
const deletingBook = ref<Book | null>(null)

const categoryOptions = computed(() => {
  return categories.value.map(c => ({
    value: c.category,
    label: `${c.category} (${c.count})`,
  }))
})

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
      category: categoryFilter.value || undefined,
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

function resetForm() {
  form.title = ''
  form.author = ''
  form.isbn = ''
  form.category = ''
  form.description = ''
  form.cover_color = '#1F3D2B'
  form.publisher = ''
  form.published_year = undefined
  form.total_copies = 1
  Object.keys(formErrors).forEach(key => {
    delete formErrors[key as keyof FormErrors]
  })
}

function validateForm(): boolean {
  Object.keys(formErrors).forEach(key => {
    delete formErrors[key as keyof FormErrors]
  })

  let valid = true

  if (!form.title.trim()) {
    formErrors.title = '请输入书名'
    valid = false
  }

  if (!form.author.trim()) {
    formErrors.author = '请输入作者'
    valid = false
  }

  if (form.published_year !== undefined) {
    const year = Number(form.published_year)
    if (isNaN(year) || year < 1000 || year > new Date().getFullYear()) {
      formErrors.published_year = '请输入有效的出版年份'
      valid = false
    }
  }

  if (form.total_copies !== undefined) {
    const copies = Number(form.total_copies)
    if (isNaN(copies) || copies < 1) {
      formErrors.total_copies = '总册数必须大于0'
      valid = false
    }
  }

  return valid
}

function openCreateModal() {
  isEdit.value = false
  editingBook.value = null
  resetForm()
  showModal.value = true
}

function openEditModal(book: Book) {
  isEdit.value = true
  editingBook.value = book
  form.title = book.title
  form.author = book.author
  form.isbn = book.isbn
  form.category = book.category
  form.description = book.description
  form.cover_color = book.cover_color
  form.publisher = book.publisher
  form.published_year = book.published_year
  form.total_copies = book.total_copies
  showModal.value = true
}

async function onSubmit() {
  if (!validateForm()) return

  formLoading.value = true
  try {
    const submitData: BookInput = {
      title: form.title.trim(),
      author: form.author.trim(),
      isbn: form.isbn?.trim() || undefined,
      category: form.category?.trim() || undefined,
      description: form.description?.trim() || undefined,
      cover_color: form.cover_color || undefined,
      publisher: form.publisher?.trim() || undefined,
      published_year: form.published_year ? Number(form.published_year) : undefined,
      total_copies: form.total_copies ? Number(form.total_copies) : undefined,
    }

    if (isEdit.value && editingBook.value) {
      await booksApi.update(editingBook.value.id, submitData)
      success(`成功更新《${form.title}》`)
    } else {
      await booksApi.create(submitData)
      success(`成功添加《${form.title}》`)
    }

    showModal.value = false
    loadBooks()
    loadCategories()
  } catch (e) {
    error(e instanceof Error ? e.message : '保存失败')
  } finally {
    formLoading.value = false
  }
}

function openDeleteConfirm(book: Book) {
  deletingBook.value = book
  showDeleteConfirm.value = true
}

async function onDelete() {
  if (!deletingBook.value) return

  deleteLoading.value = true
  try {
    await booksApi.remove(deletingBook.value.id)
    success(`成功删除《${deletingBook.value.title}》`)
    showDeleteConfirm.value = false
    deletingBook.value = null
    loadBooks()
    loadCategories()
  } catch (e) {
    error(e instanceof Error ? e.message : '删除失败')
  } finally {
    deleteLoading.value = false
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
        <h1 class="font-serif text-2xl font-bold text-ink-700">书籍管理</h1>
        <p class="mt-1 text-sm text-charcoal-muted">
          管理馆藏图书信息，支持新增、编辑和删除操作
        </p>
      </div>
      <BaseButton @click="openCreateModal">
        <template #icon>
          <Plus class="h-4 w-4" />
        </template>
        新增图书
      </BaseButton>
    </div>

    <div class="mb-5 flex flex-col gap-3 sm:flex-row">
      <div class="flex flex-1 flex-col gap-3 sm:flex-row sm:max-w-xl">
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
      <BaseSelect
        v-model="categoryFilter"
        :options="categoryOptions"
        placeholder="全部分类"
        class="sm:max-w-xs"
        @change="onCategoryChange"
      >
        <template #prefix>
          <Filter class="h-4 w-4" />
        </template>
      </BaseSelect>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner :size="40" />
    </div>

    <div v-else-if="books.length === 0">
      <EmptyState
        :icon="BookOpen"
        title="暂无符合条件的图书"
        description="试试其他搜索关键词或分类，或者点击新增图书添加"
      />
    </div>

    <div v-else class="overflow-hidden rounded-2xl border border-paper-300 bg-paper-50 shadow-card">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-paper-300 bg-ink-50/50">
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">图书</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">ISBN</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">分类</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">出版社</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">出版年</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">总册数</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">可借</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">创建时间</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="book in books"
              :key="book.id"
              class="border-b border-paper-300/60 transition-colors hover:bg-ink-50/30 last:border-b-0"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-12 shrink-0">
                    <BookCover
                      :color="book.cover_color"
                      :title="book.title"
                      :author="book.author"
                      size="sm"
                    />
                  </div>
                  <div class="min-w-0">
                    <h3 class="line-clamp-1 font-serif text-sm font-semibold text-ink-700">
                      {{ book.title }}
                    </h3>
                    <p class="mt-0.5 line-clamp-1 text-xs text-charcoal-muted">{{ book.author }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="font-mono text-xs text-charcoal-soft">{{ book.isbn || '—' }}</span>
              </td>
              <td class="px-4 py-3">
                <span class="inline-flex items-center rounded-full bg-ink-50 px-2 py-0.5 text-xs font-medium text-ink-700">
                  {{ book.category }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span class="text-sm text-charcoal-soft">{{ book.publisher || '—' }}</span>
              </td>
              <td class="px-4 py-3">
                <span class="text-sm text-charcoal-soft">{{ book.published_year || '—' }}</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-sm text-charcoal">{{ book.total_copies }}</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span
                  :class="[
                    'font-mono text-sm',
                    book.available_copies > 0 ? 'text-ink-600' : 'text-red-600',
                  ]"
                >
                  {{ book.available_copies }}
                </span>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-xs text-charcoal-muted">{{ formatDate(book.created_at) }}</span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-center gap-1.5">
                  <BaseButton size="sm" variant="ghost" @click="openEditModal(book)">
                    <template #icon>
                      <Edit class="h-3.5 w-3.5" />
                    </template>
                    编辑
                  </BaseButton>
                  <BaseButton size="sm" variant="ghost" class="text-red-500 hover:text-red-600" @click="openDeleteConfirm(book)">
                    <template #icon>
                      <Trash2 class="h-3.5 w-3.5" />
                    </template>
                    删除
                  </BaseButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="books.length > 0" class="mt-6">
      <Pagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="loadBooks"
      />
    </div>

    <BaseModal
      v-model="showModal"
      :title="isEdit ? '编辑图书' : '新增图书'"
      width="max-w-2xl"
    >
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="md:col-span-2">
          <BaseInput
            v-model="form.title"
            label="书名"
            placeholder="请输入书名"
            required
            :error="formErrors.title"
          />
        </div>
        <div class="md:col-span-2">
          <BaseInput
            v-model="form.author"
            label="作者"
            placeholder="请输入作者"
            required
            :error="formErrors.author"
          />
        </div>
        <BaseInput
          v-model="form.isbn"
          label="ISBN"
          placeholder="请输入ISBN"
          :error="formErrors.isbn"
        />
        <BaseSelect
          v-model="form.category"
          label="分类"
          :options="categories.map(c => ({ value: c.category, label: c.category }))"
          placeholder="请选择分类"
          :error="formErrors.category"
        />
        <BaseSelect
          v-model="form.cover_color"
          label="封面颜色"
          :options="coverColorOptions"
          placeholder="请选择封面颜色"
          :error="formErrors.cover_color"
        />
        <BaseInput
          v-model="form.publisher"
          label="出版社"
          placeholder="请输入出版社"
          :error="formErrors.publisher"
        />
        <BaseInput
          v-model.number="form.published_year"
          label="出版年"
          type="number"
          placeholder="请输入出版年份"
          :error="formErrors.published_year"
        />
        <BaseInput
          v-model.number="form.total_copies"
          label="总册数"
          type="number"
          placeholder="请输入总册数"
          required
          :error="formErrors.total_copies"
        />
        <div class="md:col-span-2">
          <BaseTextarea
            v-model="form.description"
            label="描述"
            placeholder="请输入图书描述"
            :rows="3"
          />
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2.5">
          <BaseButton variant="ghost" size="md" @click="showModal = false">
            <template #icon>
              <X class="h-4 w-4" />
            </template>
            取消
          </BaseButton>
          <BaseButton size="md" :loading="formLoading" @click="onSubmit">
            <template #icon>
              <Save class="h-4 w-4" />
            </template>
            {{ isEdit ? '保存修改' : '添加图书' }}
          </BaseButton>
        </div>
      </template>
    </BaseModal>

    <ConfirmDialog
      v-model="showDeleteConfirm"
      title="确认删除"
      :message="`您确定要删除《${deletingBook?.title}》吗？此操作不可撤销，将永久删除该图书信息。`"
      confirm-text="确认删除"
      danger
      :loading="deleteLoading"
      @confirm="onDelete"
    />

    <ToastContainer />
  </div>
</template>
