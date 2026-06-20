<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    color?: string
    title?: string
    author?: string
    category?: string
    size?: 'sm' | 'md' | 'lg'
  }>(),
  { color: '#1F3D2B', size: 'md' },
)

const bg = computed(() => {
  const c = props.color || '#1F3D2B'
  return `linear-gradient(150deg, ${c} 0%, ${c} 55%, rgba(0,0,0,0.35) 100%)`
})

const heights: Record<string, string> = {
  sm: 'h-28',
  md: 'h-40',
  lg: 'h-56',
}

function initials(s: string) {
  return s ? s.slice(0, 1) : ''
}
</script>

<template>
  <div
    :class="[
      'group relative flex w-full flex-col justify-between overflow-hidden rounded-xl p-3.5 text-paper-50 shadow-card transition-transform duration-300 group-hover:-translate-y-0.5',
      heights[size],
    ]"
    :style="{ backgroundImage: bg }"
  >
    <!-- 纸张噪点纹理 -->
    <div
      class="pointer-events-none absolute inset-0 opacity-30 mix-blend-overlay"
      style="background-image: radial-gradient(rgba(255,255,255,0.25) 0.5px, transparent 0.5px); background-size: 3px 3px;"
    />
    <!-- 黄铜书脊 -->
    <div class="absolute left-0 top-0 h-full w-1 bg-brass-400/80" />
    <div class="absolute inset-y-0 left-1 w-px bg-paper-50/15" />

    <div class="relative flex items-start justify-between">
      <span
        v-if="category"
        class="rounded-md bg-paper-50/15 px-1.5 py-0.5 text-[10px] font-medium tracking-wide backdrop-blur-sm"
      >
        {{ category }}
      </span>
      <span
        class="flex h-7 w-7 items-center justify-center rounded-full bg-paper-50/15 font-serif text-sm font-bold backdrop-blur-sm"
        >{{ initials(title) }}</span
      >
    </div>

    <div class="relative">
      <p class="line-clamp-2 font-serif text-[15px] font-bold leading-snug drop-shadow-sm">
        {{ title }}
      </p>
      <p v-if="author" class="mt-1 line-clamp-1 text-[11px] text-paper-50/75">
        {{ author }}
      </p>
    </div>
  </div>
</template>
