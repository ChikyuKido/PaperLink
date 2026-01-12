<script setup lang="ts">
import type { DropdownMenuSeparatorProps } from "reka-ui"
import { DropdownMenuSeparator, useForwardProps } from "reka-ui"
import { computed } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<DropdownMenuSeparatorProps & { class?: string }>()
const forwarded = useForwardProps(props)

const sepClass = computed(() => cn("-mx-1 my-1 h-px bg-neutral-200 dark:bg-neutral-800", props.class))
</script>

<template>
  <DropdownMenuSeparator v-bind="forwarded" :class="sepClass" />
</template>
