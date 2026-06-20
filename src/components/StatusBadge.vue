<script setup lang="ts">
import { computed } from 'vue'
import { statusMeta, type StatusMeta } from '@/lib/format'
import type { BorrowStatus } from '@/lib/api'
import { cn } from '@/lib/utils'

const props = defineProps<{ status: BorrowStatus }>()
const meta = computed<StatusMeta>(() => statusMeta(props.status))
</script>

<template>
  <span
    :class="
      cn(
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium',
        meta.bg,
        meta.text,
      )
    "
  >
    <span :class="['h-1.5 w-1.5 rounded-full', meta.dot]" />
    {{ meta.label }}
  </span>
</template>
