<script setup lang="ts">
import { cn } from "@/lib/utils"

defineProps<{ class?: string }>()
</script>

<template>
  <span
    :class="
      cn(
        'flex h-full w-full items-center justify-center rounded-full bg-neutral-900 text-[11px] font-semibold text-white',
        'dark:bg-neutral-100 dark:text-neutral-900',
        $props.class,
      )
    "
  >
    <slot />
  </span>
</template>
