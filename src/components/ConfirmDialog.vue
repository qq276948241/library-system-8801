<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'
import BaseButton from './BaseButton.vue'

withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    message?: string
    confirmText?: string
    cancelText?: string
    loading?: boolean
    danger?: boolean
  }>(),
  {
    title: '确认操作',
    confirmText: '确认',
    cancelText: '取消',
    loading: false,
    danger: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'confirm'): void
}>()
</script>

<template>
  <BaseModal
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    width="max-w-md"
  >
    <div class="flex gap-4">
      <div
        :class="[
          'flex h-11 w-11 shrink-0 items-center justify-center rounded-full',
          danger ? 'bg-red-100 text-red-600' : 'bg-brass-300/30 text-brass-600',
        ]"
      >
        <AlertTriangle class="h-5 w-5" />
      </div>
      <div class="flex-1 pt-0.5">
        <h3 class="font-serif text-lg font-bold text-ink-700">{{ title }}</h3>
        <p v-if="message" class="mt-1.5 text-sm leading-relaxed text-charcoal-soft">
          {{ message }}
        </p>
      </div>
    </div>
    <template #footer>
      <div class="flex justify-end gap-2.5">
        <BaseButton variant="ghost" size="md" @click="emit('update:modelValue', false)">
          {{ cancelText }}
        </BaseButton>
        <BaseButton
          :variant="danger ? 'danger' : 'primary'"
          size="md"
          :loading="loading"
          @click="emit('confirm')"
        >
          {{ confirmText }}
        </BaseButton>
      </div>
    </template>
  </BaseModal>
</template>
