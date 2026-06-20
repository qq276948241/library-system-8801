<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  icon?: Component
  title?: string
  description?: string
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center px-6 py-14 text-center">
    <div
      v-if="icon"
      class="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-paper-200/60 text-charcoal-muted"
    >
      <component :is="icon" class="h-6 w-6" />
    </div>
    <p v-if="title" class="font-serif text-base font-semibold text-charcoal">{{ title }}</p>
    <p v-if="description" class="mt-1.5 max-w-xs text-sm text-charcoal-muted">
      {{ description }}
    </p>
    <div class="mt-4"><slot /></div>
  </div>
</template>
