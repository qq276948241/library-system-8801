<script setup lang="ts">
import { watch, onMounted, onUnmounted } from 'vue'
import { X } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    width?: string
  }>(),
  { width: 'max-w-lg' },
)

const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()

function close() {
  emit('update:modelValue', false)
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  },
)

onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => {
  window.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-charcoal/40 backdrop-blur-sm" @click="close" />
        <div
          :class="
            cn(
              'relative w-full rounded-2xl bg-paper-50 shadow-lift border border-paper-300 animate-fade-up',
              width,
            )
          "
        >
          <div
            v-if="title"
            class="flex items-center justify-between border-b border-paper-300 px-6 py-4"
          >
            <h3 class="font-serif text-lg font-bold text-ink-700">{{ title }}</h3>
            <button
              class="rounded-lg p-1 text-charcoal-muted transition hover:bg-ink-50 hover:text-ink-700"
              @click="close"
            >
              <X class="w-5 h-5" />
            </button>
          </div>
          <button
            v-else
            class="absolute right-3 top-3 z-10 rounded-lg p-1 text-charcoal-muted transition hover:bg-ink-50 hover:text-ink-700"
            @click="close"
          >
            <X class="w-5 h-5" />
          </button>
          <div class="px-6 py-5">
            <slot />
          </div>
          <div v-if="$slots.footer" class="border-t border-paper-300 px-6 py-4">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
