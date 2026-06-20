<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    size?: 'sm' | 'md'
    type?: 'button' | 'submit'
    loading?: boolean
    disabled?: boolean
    block?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
    type: 'button',
    loading: false,
    disabled: false,
    block: false,
  },
)

const base =
  'inline-flex items-center justify-center gap-1.5 font-medium rounded-xl transition-all duration-200 select-none focus:outline-none focus-visible:ring-2 focus-visible:ring-ink-400/60 focus-visible:ring-offset-2 focus-visible:ring-offset-paper disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.98]'

const variants: Record<string, string> = {
  primary: 'bg-ink-600 text-paper-50 shadow-card hover:bg-ink-700 hover:shadow-lift',
  secondary:
    'bg-paper-50 text-brass-700 border border-brass-400/70 hover:bg-brass-300/15 hover:border-brass-500',
  ghost: 'bg-transparent text-charcoal-soft hover:bg-ink-50/70 hover:text-ink-700',
  danger: 'bg-red-600 text-white shadow-card hover:bg-red-700',
}

const sizes: Record<string, string> = {
  sm: 'text-xs px-3 py-1.5',
  md: 'text-sm px-4 py-2.5',
}

const classes = computed(() =>
  cn(base, variants[props.variant], sizes[props.size], props.block && 'w-full'),
)
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
    <slot v-else name="icon" />
    <slot />
  </button>
</template>
