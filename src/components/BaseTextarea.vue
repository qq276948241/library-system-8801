<script setup lang="ts">
import { cn } from '@/lib/utils'

withDefaults(
  defineProps<{
    modelValue?: string
    label?: string
    placeholder?: string
    required?: boolean
    rows?: number
    error?: string
  }>(),
  { required: false, rows: 3 },
)

const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 flex items-center gap-1 text-sm font-medium text-charcoal">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </span>
    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      :required="required"
      :rows="rows"
      @input="onInput"
      :class="
        cn(
          'w-full resize-y rounded-xl border border-paper-300 bg-paper-50/80 px-3.5 py-2.5 text-sm text-charcoal placeholder:text-charcoal-muted/70 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ink-400/40 focus:border-ink-500',
          error && 'border-red-400 focus:border-red-500 focus:ring-red-300/40',
        )
      "
    />
    <span v-if="error" class="mt-1 block text-xs text-red-600">{{ error }}</span>
  </label>
</template>
