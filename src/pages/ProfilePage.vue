<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { UserRound, Mail, Calendar, Key, Lock, Eye, EyeOff, BookOpen, CheckCircle, AlertTriangle } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput.vue'
import BaseButton from '@/components/BaseButton.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import Spinner from '@/components/Spinner.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { authApi, borrowsApi, type User, type Borrow } from '@/lib/api'
import { useAuth } from '@/composables/useAuth'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import { formatDate } from '@/lib/format'
import { initials } from '@/lib/format'
import { cn } from '@/lib/utils'

useTheme()
const { user, logout, refreshUser } = useAuth()
const { success, error } = useToast()
const router = useRouter()

const loading = ref(false)
const submitLoading = ref(false)
const showConfirmDialog = ref(false)
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const formErrors = ref<Record<string, string>>({})

const currentBorrows = ref(0)
const completedBorrows = ref(0)
const overdueCount = ref(0)

const roleLabel = computed(() => {
  return user.value?.role === 'admin' ? '管理员' : '学生'
})

const roleBadgeClass = computed(() => {
  return user.value?.role === 'admin'
    ? 'bg-red-50 text-red-700 border-red-200'
    : 'bg-ink-50 text-ink-700 border-ink-200'
})

async function loadBorrowStats() {
  try {
    const res = await borrowsApi.list({ pageSize: 1000 })
    const borrows: Borrow[] = res.list
    currentBorrows.value = borrows.filter((b) => b.status === 'borrowed').length
    completedBorrows.value = borrows.filter((b) => b.status === 'returned').length
    overdueCount.value = borrows.filter((b) => b.status === 'overdue').length
  } catch (e) {
    // ignore
  }
}

function validateForm(): boolean {
  formErrors.value = {}

  if (!oldPassword.value) {
    formErrors.value.oldPassword = '请输入原密码'
  }

  if (!newPassword.value) {
    formErrors.value.newPassword = '请输入新密码'
  } else if (newPassword.value.length < 6) {
    formErrors.value.newPassword = '新密码长度不能少于6位'
  }

  if (!confirmPassword.value) {
    formErrors.value.confirmPassword = '请确认新密码'
  } else if (newPassword.value !== confirmPassword.value) {
    formErrors.value.confirmPassword = '两次输入的密码不一致'
  }

  return Object.keys(formErrors.value).length === 0
}

function openConfirmDialog() {
  if (!validateForm()) return
  showConfirmDialog.value = true
}

async function onChangePassword() {
  if (!validateForm()) return

  submitLoading.value = true
  try {
    await authApi.changePassword(oldPassword.value, newPassword.value)
    success('密码修改成功，请重新登录')
    showConfirmDialog.value = false

    setTimeout(() => {
      logout()
      router.push({ name: 'login' })
    }, 1500)
  } catch (e) {
    error(e instanceof Error ? e.message : '密码修改失败')
  } finally {
    submitLoading.value = false
  }
}

async function loadUserInfo() {
  loading.value = true
  try {
    await refreshUser()
    await loadBorrowStats()
  } catch (e) {
    error(e instanceof Error ? e.message : '加载用户信息失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<template>
  <div class="animate-fade-up">
    <div class="mb-6">
      <h1 class="font-serif text-2xl font-bold text-ink-700">个人中心</h1>
      <p class="mt-1 text-sm text-charcoal-muted">查看和管理您的账户信息</p>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Spinner size="lg" />
    </div>

    <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <div class="lg:col-span-1">
        <div class="rounded-2xl border border-paper-300 bg-paper-50 p-6 shadow-card">
          <div class="flex flex-col items-center text-center">
            <div
              class="mb-4 flex h-24 w-24 items-center justify-center rounded-full bg-ink-100 text-ink-600 shadow-lift"
            >
              <span class="font-serif text-3xl font-bold">{{ initials(user?.name || '') }}</span>
            </div>
            <h2 class="font-serif text-xl font-bold text-ink-700">{{ user?.name }}</h2>
            <p class="mt-0.5 text-sm text-charcoal-muted">@{{ user?.username }}</p>
            <span
              :class="[
                'mt-2 inline-flex items-center rounded-full border px-3 py-1 text-xs font-medium',
                roleBadgeClass,
              ]"
            >
              {{ roleLabel }}
            </span>
          </div>

          <div class="mt-6 space-y-3 border-t border-paper-300/60 pt-6">
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-ink-50 text-ink-600">
                <UserRound class="h-4 w-4" />
              </div>
              <div>
                <p class="text-xs text-charcoal-muted">姓名</p>
                <p class="text-sm font-medium text-charcoal">{{ user?.name }}</p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-ink-50 text-ink-600">
                <Mail class="h-4 w-4" />
              </div>
              <div>
                <p class="text-xs text-charcoal-muted">邮箱</p>
                <p class="text-sm font-medium text-charcoal">{{ user?.email || '未设置' }}</p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-ink-50 text-ink-600">
                <Calendar class="h-4 w-4" />
              </div>
              <div>
                <p class="text-xs text-charcoal-muted">注册时间</p>
                <p class="text-sm font-medium text-charcoal">{{ formatDate(user?.created_at || null) }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-6 rounded-2xl border border-paper-300 bg-paper-50 p-6 shadow-card">
          <h3 class="mb-4 font-serif text-lg font-bold text-ink-700">借阅统计</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-ink-50 text-ink-600">
                  <BookOpen class="h-4 w-4" />
                </div>
                <span class="text-sm text-charcoal-soft">当前借阅</span>
              </div>
              <span class="font-mono text-lg font-bold text-ink-700">{{ currentBorrows }}</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-brass-50 text-brass-600">
                  <CheckCircle class="h-4 w-4" />
                </div>
                <span class="text-sm text-charcoal-soft">已完成借阅</span>
              </div>
              <span class="font-mono text-lg font-bold text-brass-700">{{ completedBorrows }}</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-red-50 text-red-600">
                  <AlertTriangle class="h-4 w-4" />
                </div>
                <span class="text-sm text-charcoal-soft">逾期次数</span>
              </div>
              <span
                :class="[
                  'font-mono text-lg font-bold',
                  overdueCount > 0 ? 'text-red-600' : 'text-charcoal-soft',
                ]"
              >
                {{ overdueCount }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-2">
        <div class="rounded-2xl border border-paper-300 bg-paper-50 p-6 shadow-card">
          <div class="mb-6 flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-ink-50 text-ink-600">
              <Key class="h-5 w-5" />
            </div>
            <div>
              <h3 class="font-serif text-lg font-bold text-ink-700">修改密码</h3>
              <p class="text-sm text-charcoal-muted">定期更换密码以保护账户安全</p>
            </div>
          </div>

          <form @submit.prevent="openConfirmDialog" class="space-y-5">
            <div class="relative">
              <BaseInput
                v-model="oldPassword"
                :type="showOldPassword ? 'text' : 'password'"
                label="原密码"
                placeholder="请输入原密码"
                :error="formErrors.oldPassword"
                required
              >
                <template #prefix>
                  <Lock class="h-4 w-4" />
                </template>
              </BaseInput>
              <button
                type="button"
                class="absolute bottom-2.5 right-3 text-charcoal-muted transition hover:text-charcoal"
                @click="showOldPassword = !showOldPassword"
              >
                <component :is="showOldPassword ? EyeOff : Eye" class="h-4 w-4" />
              </button>
            </div>

            <div class="relative">
              <BaseInput
                v-model="newPassword"
                :type="showNewPassword ? 'text' : 'password'"
                label="新密码"
                placeholder="请输入新密码（至少6位）"
                :error="formErrors.newPassword"
                required
              >
                <template #prefix>
                  <Lock class="h-4 w-4" />
                </template>
              </BaseInput>
              <button
                type="button"
                class="absolute bottom-2.5 right-3 text-charcoal-muted transition hover:text-charcoal"
                @click="showNewPassword = !showNewPassword"
              >
                <component :is="showNewPassword ? EyeOff : Eye" class="h-4 w-4" />
              </button>
            </div>

            <div class="relative">
              <BaseInput
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                label="确认新密码"
                placeholder="请再次输入新密码"
                :error="formErrors.confirmPassword"
                required
              >
                <template #prefix>
                  <Lock class="h-4 w-4" />
                </template>
              </BaseInput>
              <button
                type="button"
                class="absolute bottom-2.5 right-3 text-charcoal-muted transition hover:text-charcoal"
                @click="showConfirmPassword = !showConfirmPassword"
              >
                <component :is="showConfirmPassword ? EyeOff : Eye" class="h-4 w-4" />
              </button>
            </div>

            <div class="pt-2">
              <BaseButton type="submit" :loading="submitLoading" block>
                <template #icon>
                  <Key class="h-4 w-4" />
                </template>
                修改密码
              </BaseButton>
            </div>
          </form>

          <div class="mt-6 rounded-xl bg-ink-50/50 p-4 text-xs text-charcoal-muted">
            <p class="mb-1 font-medium text-charcoal-soft">温馨提示：</p>
            <p>• 密码长度至少为6位</p>
            <p>• 密码修改成功后，您将被自动退出登录</p>
            <p>• 请妥善保管您的密码，不要泄露给他人</p>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      v-model="showConfirmDialog"
      title="确认修改密码"
      message="您确定要修改密码吗？修改成功后将自动退出登录，请使用新密码重新登录。"
      confirm-text="确认修改"
      :loading="submitLoading"
      danger
      @confirm="onChangePassword"
    />

    <ToastContainer />
  </div>
</template>
