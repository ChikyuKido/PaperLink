<script setup lang="ts">
import { cn } from "@/lib/utils"

defineProps<{ class?: string }>()
</script>

<template>
  <span
    :class="
      cn(
        'relative flex h-9 w-9 shrink-0 overflow-hidden rounded-full border border-neutral-200 bg-neutral-100 text-neutral-900',
        'dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-50',
        $props.class,
      )
    "
  >
    <slot />
  </span>
</template>
