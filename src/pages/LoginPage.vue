<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { BookMarked, User, Lock, LogIn } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput.vue'
import BaseButton from '@/components/BaseButton.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { useAuth } from '@/composables/useAuth'
import { useToast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'

const { login, isLogged, isAdmin } = useAuth()
const { success, error } = useToast()
const router = useRouter()
const route = useRoute()
useTheme()

const username = ref('')
const password = ref('')
const loading = ref(false)

onMounted(() => {
  if (isLogged.value) {
    router.push(isAdmin.value ? { name: 'admin-borrows' } : { name: 'books' })
  }
})

async function onSubmit() {
  if (!username.value || !password.value) {
    error('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const user = await login(username.value, password.value)
    success(`欢迎回来，${user.name}！`)
    const redirect = route.query.redirect as string | undefined
    if (redirect) {
      router.push(redirect)
    } else {
      router.push(user.role === 'admin' ? { name: 'admin-borrows' } : { name: 'books' })
    }
  } catch (e) {
    error(e instanceof Error ? e.message : '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="paper-grain flex min-h-screen items-center justify-center px-4 py-10">
    <div class="w-full max-w-md animate-fade-up">
      <div class="mb-8 text-center">
        <div class="mb-4 flex justify-center">
          <span class="flex h-16 w-16 items-center justify-center rounded-2xl bg-ink-600 text-paper-50 shadow-lift">
            <BookMarked class="h-8 w-8" />
          </span>
        </div>
        <h1 class="mb-1.5 font-serif text-2xl font-bold text-ink-700">青檀图书馆</h1>
        <p class="text-sm text-charcoal-muted">请登录以继续使用借阅服务</p>
      </div>

      <div class="rounded-2xl border border-paper-300 bg-paper-50 p-6 shadow-lift">
        <form @submit.prevent="onSubmit" class="space-y-4">
          <BaseInput
            v-model="username"
            label="用户名"
            placeholder="请输入用户名"
            required
          >
            <template #prefix>
              <User class="h-4 w-4" />
            </template>
          </BaseInput>

          <BaseInput
            v-model="password"
            label="密码"
            type="password"
            placeholder="请输入密码"
            required
          >
            <template #prefix>
              <Lock class="h-4 w-4" />
            </template>
          </BaseInput>

          <BaseButton type="submit" block :loading="loading">
            <template #icon>
              <LogIn class="h-4 w-4" />
            </template>
            登 录
          </BaseButton>
        </form>

        <div class="mt-5 rounded-xl bg-ink-50/50 p-3 text-xs text-charcoal-muted">
          <p class="mb-1 font-medium text-charcoal-soft">测试账号：</p>
          <p>管理员：<code class="rounded bg-paper-200 px-1.5 py-0.5 font-mono">admin / admin123</code></p>
          <p>学生：<code class="rounded bg-paper-200 px-1.5 py-0.5 font-mono">zhangsan / student123</code></p>
        </div>
      </div>

      <p class="mt-6 text-center text-xs text-charcoal-muted">
        青檀图书馆借阅管理系统 · 校园数字化借阅服务
      </p>
    </div>
    <ToastContainer />
  </div>
</template>
