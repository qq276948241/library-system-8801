<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Option {
  value: string | number
  label: string
}

withDefaults(
  defineProps<{
    modelValue?: string | number
    label?: string
    options: Option[]
    placeholder?: string
    required?: boolean
    disabled?: boolean
    error?: string
  }>(),
  { required: false, disabled: false },
)

const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

function onChange(e: Event) {
  emit('update:modelValue', (e.target as HTMLSelectElement).value)
}
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 flex items-center gap-1 text-sm font-medium text-charcoal">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </span>
    <select
      :value="modelValue"
      :required="required"
      :disabled="disabled"
      @change="onChange"
      :class="
        cn(
          'w-full appearance-none rounded-xl border bg-paper-50/80 bg-none px-3.5 py-2.5 pr-9 text-sm text-charcoal transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ink-400/40 focus:border-ink-500 disabled:bg-paper-200/50 disabled:cursor-not-allowed',
          'bg-[url(\'data:image/svg+xml;utf8,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%2216%22 height=%2216%22 viewBox=%220 0 24 24%22 fill=%22none%22 stroke=%22%232A2620%22 stroke-width=%222%22 stroke-linecap=%22round%22 stroke-linejoin=%22round%22><polyline points=%226 9 12 15 18 9%22></polyline></svg>\')] bg-[length:16px] bg-[right_0.75rem_center] bg-no-repeat',
          error ? 'border-red-400' : 'border-paper-300',
        )
      "
    >
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option>
    </select>
    <span v-if="error" class="mt-1 block text-xs text-red-600">{{ error }}</span>
  </label>
</template>
