<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import { BookMarked, Library, Menu, X, ClipboardList, UserRound } from 'lucide-vue-next'
import UserMenu from '@/components/UserMenu.vue'



const navOpen = ref(false)

const links = [
  { to: { name: 'books' }, label: '图书列表', icon: Library },
  { to: { name: 'borrows' }, label: '借阅记录', icon: ClipboardList },
  { to: { name: 'profile' }, label: '个人中心', icon: UserRound },
]
</script>

<template>
  <div class="paper-grain flex min-h-screen flex-col">
    <header
      class="sticky top-0 z-30 border-b border-paper-300/80 bg-paper-100/85 backdrop-blur-md"
    >
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6">
        <RouterLink :to="{ name: 'books' }" class="flex items-center gap-2.5">
          <span
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-ink-600 text-paper-50 shadow-card"
          >
            <BookMarked class="h-5 w-5" />
          </span>
          <span class="flex flex-col leading-none">
            <span class="font-serif text-lg font-bold text-ink-700">青檀图书馆</span>
            <span class="text-[10px] tracking-widest text-brass-600">QINGTAN LIBRARY</span>
          </span>
        </RouterLink>

        <nav class="hidden items-center gap-1 md:flex">
          <RouterLink
            v-for="l in links"
            :key="l.label"
            :to="l.to"
            class="group relative flex items-center gap-1.5 rounded-lg px-3.5 py-2 text-sm font-medium text-charcoal-soft transition hover:bg-ink-50 hover:text-ink-700"
            active-class="text-ink-700 bg-ink-50/60"
          >
            <component :is="l.icon" class="h-4 w-4" />
            {{ l.label }}
            <span
              class="absolute -bottom-px left-3.5 right-3.5 h-0.5 origin-left scale-x-0 rounded-full bg-brass-500 transition-transform duration-200 group-[.router-link-active]:scale-x-100"
            />
          </RouterLink>
        </nav>

        <div class="flex items-center gap-2">
          <div class="hidden md:block">
            <UserMenu />
          </div>
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-charcoal-soft transition hover:bg-ink-50 md:hidden"
            @click="navOpen = !navOpen"
          >
            <component :is="navOpen ? X : Menu" class="h-5 w-5" />
          </button>
        </div>
      </div>

      <Transition name="slide">
        <div v-if="navOpen" class="border-t border-paper-300/80 px-4 py-3 md:hidden">
          <RouterLink
            v-for="l in links"
            :key="l.label"
            :to="l.to"
            class="flex items-center gap-2.5 rounded-lg px-3 py-2.5 text-sm font-medium text-charcoal-soft transition hover:bg-ink-50 hover:text-ink-700"
            active-class="bg-ink-50/70 text-ink-700"
            @click="navOpen = false"
          >
            <component :is="l.icon" class="h-4 w-4" />
            {{ l.label }}
          </RouterLink>
          <div class="mt-2 border-t border-paper-300 pt-3">
            <UserMenu align="left" />
          </div>
        </div>
      </Transition>
    </header>

    <main class="mx-auto w-full max-w-6xl flex-1 px-4 py-6 sm:px-6 sm:py-8">
      <RouterView />
    </main>

    <footer class="border-t border-paper-300/60 py-5 text-center text-xs text-charcoal-muted">
      青檀图书馆借阅管理系统 · 校园数字化借阅服务
    </footer>
  </div>
</template>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.22s ease;
  overflow: hidden;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
}
.slide-enter-to,
.slide-leave-from {
  max-height: 320px;
}
</style>
