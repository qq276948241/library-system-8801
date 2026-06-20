<script setup lang="ts">
import { computed, useSlots } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    label?: string
    type?: string
    placeholder?: string
    error?: string
    required?: boolean
    disabled?: boolean
    hint?: string
  }>(),
  {
    type: 'text',
    required: false,
    disabled: false,
  },
)

const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()
const slots = useSlots()

const hasIcon = computed(() => !!slots.prefix)

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
}
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 flex items-center gap-1 text-sm font-medium text-charcoal">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </span>
    <div class="relative">
      <span
        v-if="hasIcon"
        class="absolute inset-y-0 left-0 flex items-center pl-3 text-charcoal-muted"
      >
        <slot name="prefix" />
      </span>
      <input
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :required="required"
        :disabled="disabled"
        @input="onInput"
        :class="
          cn(
            'w-full rounded-xl border bg-paper-50/80 px-3.5 py-2.5 text-sm text-charcoal placeholder:text-charcoal-muted/70 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ink-400/40 focus:border-ink-500 disabled:bg-paper-200/50 disabled:cursor-not-allowed',
            hasIcon && 'pl-10',
            error ? 'border-red-400 focus:border-red-500 focus:ring-red-300/40' : 'border-paper-300',
          )
        "
      />
    </div>
    <span v-if="error" class="mt-1 block text-xs text-red-600">{{ error }}</span>
    <span v-else-if="hint" class="mt-1 block text-xs text-charcoal-muted">{{ hint }}</span>
  </label>
</template>
