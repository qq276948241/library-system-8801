<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ChevronDown, LogOut } from 'lucide-vue-next'
import { useAuth } from '@/composables/useAuth'
import { initials } from '@/lib/format'

const props = defineProps<{ align?: 'left' | 'right' }>()

const { user, logout } = useAuth()
const router = useRouter()
const open = ref(false)

function onLogout() {
  open.value = false
  logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="relative">
    <button
      class="flex items-center gap-2 rounded-full py-1 pl-1 pr-2.5 transition hover:bg-ink-50"
      @click="open = !open"
    >
      <span
        class="flex h-8 w-8 items-center justify-center rounded-full bg-ink-600 font-serif text-sm font-bold text-paper-50"
      >
        {{ initials(user?.name || '') }}
      </span>
      <span class="hidden text-left sm:block">
        <span class="block text-sm font-medium leading-tight text-charcoal">{{ user?.name }}</span>
        <span class="block text-[11px] leading-tight text-charcoal-muted">
          {{ user?.role === 'admin' ? '管理员' : '学生' }}
        </span>
      </span>
      <ChevronDown
        class="hidden h-4 w-4 text-charcoal-muted transition-transform sm:block"
        :class="{ 'rotate-180': open }"
      />
    </button>

    <div
      v-if="open"
      class="fixed inset-0 z-40"
      @click="open = false"
    />
    <Transition name="dd">
      <div
        v-if="open"
        :class="[
          'absolute z-50 mt-2 w-60 rounded-xl border border-paper-300 bg-paper-50 p-1.5 shadow-lift',
          props.align === 'left' ? 'left-0' : 'right-0',
        ]"
      >
        <div class="flex items-center gap-2.5 rounded-lg px-2.5 py-2">
          <span
            class="flex h-9 w-9 items-center justify-center rounded-full bg-ink-600 font-serif text-paper-50"
          >
            {{ initials(user?.name || '') }}
          </span>
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-charcoal">{{ user?.name }}</p>
            <p class="truncate text-xs text-charcoal-muted">@{{ user?.username }}</p>
          </div>
        </div>
        <div class="my-1 h-px bg-paper-300" />
        <slot />
        <button
          class="mt-0.5 flex w-full items-center gap-2 rounded-lg px-2.5 py-2 text-sm text-red-600 transition hover:bg-red-50"
          @click="onLogout"
        >
          <LogOut class="h-4 w-4" />
          退出登录
        </button>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.dd-enter-active,
.dd-leave-active {
  transition: all 0.18s ease;
  transform-origin: top;
}
.dd-enter-from,
.dd-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.97);
}
</style>
