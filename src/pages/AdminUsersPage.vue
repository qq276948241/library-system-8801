<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { Search, Plus, Edit, Trash2, Users, UserRound, Mail, Shield, Save, X } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseModal from '@/components/BaseModal.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import { usersApi, type User, type UserInput, type Role } from '@/lib/api'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { formatDate, initials } from '@/lib/format'
import { cn } from '@/lib/utils'

useTheme()
const { success, error } = useToast()

const roleOptions = [
  { value: 'admin', label: '管理员' },
  { value: 'student', label: '学生' },
]

interface FormErrors {
  username?: string
  password?: string
  name?: string
  role?: string
  email?: string
}

const users = ref<User[]>([])
const loading = ref(false)
const formLoading = ref(false)
const deleteLoading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref('')

const showModal = ref(false)
const isEdit = ref(false)
const editingUser = ref<User | null>(null)

const form = reactive<UserInput>({
  username: '',
  password: '',
  name: '',
  role: 'student',
  email: '',
})

const formErrors = reactive<FormErrors>({})

const showDeleteConfirm = ref(false)
const deletingUser = ref<User | null>(null)

function roleLabel(role: Role): string {
  return role === 'admin' ? '管理员' : '学生'
}

function avatarBgClass(role: Role): string {
  return role === 'admin' ? 'bg-ink-500' : 'bg-brass-500'
}

async function loadUsers() {
  loading.value = true
  try {
    const res = await usersApi.list({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value || undefined,
    })
    users.value = res.list
    total.value = res.total
  } catch (e) {
    error(e instanceof Error ? e.message : '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

function onSearch() {
  page.value = 1
  loadUsers()
}

function resetForm() {
  form.username = ''
  form.password = ''
  form.name = ''
  form.role = 'student'
  form.email = ''
  Object.keys(formErrors).forEach(key => {
    delete formErrors[key as keyof FormErrors]
  })
}

function validateForm(): boolean {
  Object.keys(formErrors).forEach(key => {
    delete formErrors[key as keyof FormErrors]
  })

  let valid = true

  if (!form.username?.trim()) {
    formErrors.username = '请输入用户名'
    valid = false
  }

  if (!isEdit.value && !form.password?.trim()) {
    formErrors.password = '请输入密码'
    valid = false
  }

  if (form.password && form.password.length < 6) {
    formErrors.password = '密码长度至少6位'
    valid = false
  }

  if (!form.name?.trim()) {
    formErrors.name = '请输入姓名'
    valid = false
  }

  if (!form.role) {
    formErrors.role = '请选择角色'
    valid = false
  }

  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    formErrors.email = '请输入有效的邮箱地址'
    valid = false
  }

  return valid
}

function openCreateModal() {
  isEdit.value = false
  editingUser.value = null
  resetForm()
  showModal.value = true
}

function openEditModal(user: User) {
  isEdit.value = true
  editingUser.value = user
  form.username = user.username
  form.password = ''
  form.name = user.name
  form.role = user.role
  form.email = user.email
  showModal.value = true
}

async function onSubmit() {
  if (!validateForm()) return

  formLoading.value = true
  try {
    const submitData: UserInput = {
      username: form.username?.trim(),
      password: form.password?.trim() || undefined,
      name: form.name?.trim(),
      role: form.role,
      email: form.email?.trim() || undefined,
    }

    if (isEdit.value && editingUser.value) {
      await usersApi.update(editingUser.value.id, submitData)
      success(`成功更新用户 ${form.name}`)
    } else {
      await usersApi.create(submitData)
      success(`成功添加用户 ${form.name}`)
    }

    showModal.value = false
    loadUsers()
  } catch (e) {
    error(e instanceof Error ? e.message : '保存失败')
  } finally {
    formLoading.value = false
  }
}

function openDeleteConfirm(user: User) {
  deletingUser.value = user
  showDeleteConfirm.value = true
}

async function onDelete() {
  if (!deletingUser.value) return

  deleteLoading.value = true
  try {
    await usersApi.remove(deletingUser.value.id)
    success(`成功删除用户 ${deletingUser.value.name}`)
    showDeleteConfirm.value = false
    deletingUser.value = null
    loadUsers()
  } catch (e) {
    error(e instanceof Error ? e.message : '删除失败')
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<template>
  <div class="animate-fade-up">
    <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="font-serif text-2xl font-bold text-ink-700">用户管理</h1>
        <p class="mt-1 text-sm text-charcoal-muted">
          管理系统用户信息，支持新增、编辑和删除操作
        </p>
      </div>
      <BaseButton @click="openCreateModal">
        <template #icon>
          <Plus class="h-4 w-4" />
        </template>
        新增用户
      </BaseButton>
    </div>

    <div class="mb-5 flex flex-col gap-3 sm:flex-row">
      <div class="flex flex-1 flex-col gap-3 sm:flex-row sm:max-w-xl">
        <BaseInput
          v-model="keyword"
          placeholder="搜索用户名、姓名、邮箱..."
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

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner :size="40" />
    </div>

    <div v-else-if="users.length === 0">
      <EmptyState
        :icon="Users"
        title="暂无符合条件的用户"
        description="试试其他搜索关键词，或者点击新增用户添加"
      />
    </div>

    <div v-else class="overflow-hidden rounded-2xl border border-paper-300 bg-paper-50 shadow-card">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-paper-300 bg-ink-50/50">
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">用户</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">用户名</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">角色</th>
              <th class="px-4 py-3 text-left font-serif text-sm font-semibold text-ink-700">邮箱</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">注册时间</th>
              <th class="px-4 py-3 text-center font-serif text-sm font-semibold text-ink-700">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="user in users"
              :key="user.id"
              class="border-b border-paper-300/60 transition-colors hover:bg-ink-50/30 last:border-b-0"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div
                    :class="[
                      'flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-white font-medium',
                      avatarBgClass(user.role),
                    ]"
                  >
                    {{ initials(user.name) }}
                  </div>
                  <div class="min-w-0">
                    <h3 class="line-clamp-1 font-serif text-sm font-semibold text-ink-700">
                      {{ user.name }}
                    </h3>
                    <p class="mt-0.5 flex items-center gap-1 text-xs text-charcoal-muted">
                      <UserRound class="h-3 w-3" />
                      {{ user.username }}
                    </p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="font-mono text-sm text-charcoal-soft">{{ user.username }}</span>
              </td>
              <td class="px-4 py-3">
                <span
                  :class="[
                    'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium',
                    user.role === 'admin' ? 'bg-ink-50 text-ink-700' : 'bg-brass-300/30 text-brass-700',
                  ]"
                >
                  <Shield v-if="user.role === 'admin'" class="h-3 w-3" />
                  {{ roleLabel(user.role) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span v-if="user.email" class="inline-flex items-center gap-1 text-sm text-charcoal-soft">
                  <Mail class="h-3.5 w-3.5" />
                  {{ user.email }}
                </span>
                <span v-else class="text-sm text-charcoal-muted">—</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span class="font-mono text-xs text-charcoal-muted">{{ formatDate(user.created_at) }}</span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-center gap-1.5">
                  <BaseButton size="sm" variant="ghost" @click="openEditModal(user)">
                    <template #icon>
                      <Edit class="h-3.5 w-3.5" />
                    </template>
                    编辑
                  </BaseButton>
                  <BaseButton size="sm" variant="ghost" class="text-red-500 hover:text-red-600" @click="openDeleteConfirm(user)">
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

    <div v-if="users.length > 0" class="mt-6">
      <Pagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="loadUsers"
      />
    </div>

    <BaseModal
      v-model="showModal"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="max-w-xl"
    >
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="md:col-span-2">
          <BaseInput
            v-model="form.username"
            label="用户名"
            placeholder="请输入用户名"
            required
            :error="formErrors.username"
          />
        </div>
        <div class="md:col-span-2">
          <BaseInput
            v-model="form.password"
            :label="isEdit ? '密码（留空不修改）' : '密码'"
            type="password"
            :placeholder="isEdit ? '留空则不修改密码' : '请输入密码，至少6位'"
            :required="!isEdit"
            :error="formErrors.password"
          />
        </div>
        <div class="md:col-span-2">
          <BaseInput
            v-model="form.name"
            label="姓名"
            placeholder="请输入姓名"
            required
            :error="formErrors.name"
          />
        </div>
        <BaseSelect
          v-model="form.role"
          label="角色"
          :options="roleOptions"
          placeholder="请选择角色"
          required
          :error="formErrors.role"
        />
        <BaseInput
          v-model="form.email"
          label="邮箱"
          type="email"
          placeholder="请输入邮箱地址"
          :error="formErrors.email"
        />
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
            {{ isEdit ? '保存修改' : '添加用户' }}
          </BaseButton>
        </div>
      </template>
    </BaseModal>

    <ConfirmDialog
      v-model="showDeleteConfirm"
      title="确认删除"
      :message="`您确定要删除用户 ${deletingUser?.name} 吗？此操作不可撤销，将永久删除该用户信息。`"
      confirm-text="确认删除"
      danger
      :loading="deleteLoading"
      @confirm="onDelete"
    />

    <ToastContainer />
  </div>
</template>
