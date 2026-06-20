<script setup lang="ts">
import { CheckCircle2, XCircle, Info } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'
import { cn } from '@/lib/utils'

const { toasts, remove } = useToast()

const iconMap = { success: CheckCircle2, error: XCircle, info: Info }
const tone: Record<string, string> = {
  success: 'text-ink-600',
  error: 'text-red-600',
  info: 'text-brass-600',
}
</script>

<template>
  <Teleport to="body">
    <div class="pointer-events-none fixed right-4 top-4 z-[60] flex w-full max-w-xs flex-col gap-2">
      <TransitionGroup name="toast">
        <div
          v-for="t in toasts"
          :key="t.id"
          @click="remove(t.id)"
          :class="
            cn(
              'pointer-events-auto flex cursor-pointer items-start gap-2.5 rounded-xl border border-paper-300 bg-paper-50 px-4 py-3 shadow-lift animate-fade-up',
            )
          "
        >
          <component :is="iconMap[t.type]" :class="['mt-0.5 h-4 w-4 shrink-0', tone[t.type]]" />
          <p class="text-sm leading-relaxed text-charcoal">{{ t.message }}</p>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.28s cubic-bezier(0.22, 1, 0.36, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(24px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(24px);
}
.toast-move {
  transition: transform 0.28s ease;
}
</style>
