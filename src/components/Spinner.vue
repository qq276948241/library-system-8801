<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{ size?: number | 'sm' | 'md' | 'lg'; class?: string }>(),
  { size: 24 },
)

const sizeMap: Record<string, number> = {
  sm: 16,
  md: 24,
  lg: 40,
}

const pixelSize = computed(() => {
  if (typeof props.size === 'number') return props.size
  return sizeMap[props.size] ?? 24
})
</script>

<template>
  <Loader2
    :class="cn('animate-spin text-ink-500', $props.class)"
    :style="{ width: pixelSize + 'px', height: pixelSize + 'px' }"
  />
</template>
