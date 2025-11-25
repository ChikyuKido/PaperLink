<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'click'): void
}>()

function handleClick() {
  emit('click')
}
</script>

<template>
  <button
      type="button"
      aria-label="Upload PDF"
      class="fixed bottom-6 right-6 z-30 flex h-14 w-14 items-center justify-center rounded-full bg-neutral-900 text-neutral-50 border border-neutral-800 shadow-lg shadow-neutral-900/40 hover:bg-emerald-700 hover:border-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-neutral-50 dark:bg-neutral-100 dark:text-neutral-900 dark:border-neutral-300 dark:shadow-neutral-950/40 dark:hover:bg-emerald-500 dark:hover:border-emerald-400 dark:focus-visible:ring-offset-neutral-950 transition-all"
      @click="handleClick"
  >
    <Plus class="w-6 h-6" aria-hidden="true" />
  </button>
</template>

<style scoped>
</style>
