<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

withDefaults(defineProps<{ size?: number; class?: string }>(), { size: 24 })
</script>

<template>
  <Loader2 :class="cn('animate-spin text-ink-500', $props.class)" :style="{ width: size + 'px', height: size + 'px' }" />
</template>
