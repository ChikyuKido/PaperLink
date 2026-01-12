<script setup lang="ts">
import type { DropdownMenuContentEmits, DropdownMenuContentProps } from "reka-ui"
import { DropdownMenuContent, DropdownMenuPortal, useForwardPropsEmits } from "reka-ui"
import { computed } from "vue"
import { cn } from "@/lib/utils"

const props = withDefaults(defineProps<DropdownMenuContentProps & { class?: string }>(), {
  sideOffset: 4,
})
const emits = defineEmits<DropdownMenuContentEmits>()

const forwarded = useForwardPropsEmits(props, emits)

const contentClass = computed(() =>
  cn(
    "z-50 min-w-[12rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md",
    "border-neutral-200 bg-white text-neutral-900 shadow-neutral-200/60",
    "dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-50 dark:shadow-none",
    "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
    props.class,
  ),
)
</script>

<template>
  <DropdownMenuPortal>
    <DropdownMenuContent v-bind="forwarded" :class="contentClass">
      <slot />
    </DropdownMenuContent>
  </DropdownMenuPortal>
</template>
