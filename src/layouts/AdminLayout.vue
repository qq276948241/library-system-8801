<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import {
  BookMarked,
  LayoutDashboard,
  BookCopy,
  Users,
  Menu,
  X,
  ArrowLeft,
} from 'lucide-vue-next'
import UserMenu from '@/components/UserMenu.vue'

const router = useRouter()
const sidebarOpen = ref(false)

const links = [
  { to: { name: 'admin-borrows' }, label: '借阅管理', icon: LayoutDashboard },
  { to: { name: 'admin-books' }, label: '书籍管理', icon: BookCopy },
  { to: { name: 'admin-users' }, label: '用户管理', icon: Users },
]

function goFront() {
  router.push({ name: 'books' })
}
</script>

<template>
  <div class="paper-grain flex min-h-screen">
    <!-- 侧边栏 -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-40 flex w-60 flex-col bg-ink-700 text-paper-100 transition-transform duration-300 lg:static lg:translate-x-0',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
      ]"
    >
      <div class="flex items-center justify-between px-5 py-5">
        <RouterLink :to="{ name: 'admin-borrows' }" class="flex items-center gap-2.5">
          <span class="flex h-9 w-9 items-center justify-center rounded-xl bg-paper-50/10">
            <BookMarked class="h-5 w-5 text-brass-300" />
          </span>
          <span class="flex flex-col leading-none">
            <span class="font-serif text-base font-bold text-paper-50">青檀图书馆</span>
            <span class="text-[10px] tracking-widest text-brass-300/80">管理后台</span>
          </span>
        </RouterLink>
        <button class="text-paper-100/70 lg:hidden" @click="sidebarOpen = false">
          <X class="h-5 w-5" />
        </button>
      </div>

      <div class="mx-5 mb-4 h-px bg-paper-50/10" />

      <nav class="flex-1 space-y-1 px-3">
        <RouterLink
          v-for="l in links"
          :key="l.label"
          :to="l.to"
          class="group flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-sm font-medium text-paper-100/80 transition hover:bg-paper-50/10 hover:text-paper-50"
          active-class="bg-paper-50/12 text-paper-50"
          @click="sidebarOpen = false"
        >
          <component :is="l.icon" class="h-5 w-5 shrink-0 text-brass-300" />
          {{ l.label }}
          <span
            class="ml-auto h-1.5 w-1.5 rounded-full bg-brass-300 opacity-0 transition group-[.router-link-active]:opacity-100"
          />
        </RouterLink>
      </nav>

      <div class="border-t border-paper-50/10 p-3">
        <button
          class="flex w-full items-center gap-3 rounded-xl px-3.5 py-2.5 text-sm font-medium text-paper-100/70 transition hover:bg-paper-50/10 hover:text-paper-50"
          @click="goFront"
        >
          <ArrowLeft class="h-5 w-5" />
          返回前台
        </button>
      </div>
    </aside>

    <div v-if="sidebarOpen" class="fixed inset-0 z-30 bg-charcoal/40 lg:hidden" @click="sidebarOpen = false" />

    <!-- 主区域 -->
    <div class="flex min-w-0 flex-1 flex-col">
      <header
        class="sticky top-0 z-20 flex h-16 items-center justify-between border-b border-paper-300/80 bg-paper-100/85 px-4 backdrop-blur-md sm:px-6"
      >
        <div class="flex items-center gap-3">
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-charcoal-soft transition hover:bg-ink-50 lg:hidden"
            @click="sidebarOpen = true"
          >
            <Menu class="h-5 w-5" />
          </button>
          <h1 class="font-serif text-lg font-bold text-ink-700">管理后台</h1>
        </div>
        <UserMenu />
      </header>

      <main class="flex-1 px-4 py-6 sm:px-6 sm:py-8">
        <slot />
      </main>
    </div>
  </div>
</template>
