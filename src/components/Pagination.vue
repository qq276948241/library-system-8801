<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    page: number
    pageSize: number
    total: number
  }>(),
  {},
)

const emit = defineEmits<{ (e: 'update:page', v: number): void }>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const start = computed(() =>
  props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1,
)
const end = computed(() => Math.min(props.page * props.pageSize, props.total))

const pages = computed(() => {
  const tp = totalPages.value
  const cur = props.page
  const res: (number | '…')[] = []
  if (tp <= 7) {
    for (let i = 1; i <= tp; i++) res.push(i)
    return res
  }
  res.push(1)
  if (cur > 3) res.push('…')
  const from = Math.max(2, cur - 1)
  const to = Math.min(tp - 1, cur + 1)
  for (let i = from; i <= to; i++) res.push(i)
  if (cur < tp - 2) res.push('…')
  res.push(tp)
  return res
})

function go(p: number) {
  if (p < 1 || p > totalPages.value || p === props.page) return
  emit('update:page', p)
}
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3">
    <p class="text-xs text-charcoal-muted">
      共 <span class="font-mono font-medium text-charcoal">{{ total }}</span> 条 ·
      第 <span class="font-mono text-charcoal">{{ start }}-{{ end }}</span> 条
    </p>
    <div class="flex items-center gap-1">
      <button
        class="flex h-8 w-8 items-center justify-center rounded-lg text-charcoal-soft transition hover:bg-ink-50 hover:text-ink-700 disabled:opacity-40 disabled:hover:bg-transparent"
        :disabled="page <= 1"
        @click="go(page - 1)"
      >
        <ChevronLeft class="h-4 w-4" />
      </button>
      <template v-for="(p, i) in pages" :key="i">
        <span v-if="p === '…'" class="px-1.5 text-charcoal-muted">…</span>
        <button
          v-else
          @click="go(p)"
          :class="[
            'flex h-8 min-w-8 items-center justify-center rounded-lg px-2 font-mono text-xs transition',
            p === page
              ? 'bg-ink-600 text-paper-50 shadow-card'
              : 'text-charcoal-soft hover:bg-ink-50 hover:text-ink-700',
          ]"
        >
          {{ p }}
        </button>
      </template>
      <button
        class="flex h-8 w-8 items-center justify-center rounded-lg text-charcoal-soft transition hover:bg-ink-50 hover:text-ink-700 disabled:opacity-40 disabled:hover:bg-transparent"
        :disabled="page >= totalPages"
        @click="go(page + 1)"
      >
        <ChevronRight class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>
