<script setup lang="ts">
import type { TabsTriggerProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsTrigger, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<
    TabsTriggerProps & { class?: HTMLAttributes["class"] }
>()

const delegatedProps = reactiveOmit(props, "class")
const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <TabsTrigger
      data-slot="tabs-trigger"
      :class="cn(
      `
      inline-flex
      h-full
      flex-1
      items-center
      justify-center
      gap-1.5
      rounded-md
      px-2
      py-1
      text-sm
      font-medium
      whitespace-nowrap
      transition-colors

      border-0
      ring-0
      shadow-none
      outline-none

      focus:outline-none
      focus-visible:outline-none
      focus-visible:ring-0
      focus-visible:border-0

      disabled:pointer-events-none
      disabled:opacity-50

      data-[state=active]:bg-transparent
      data-[state=active]:text-white

      data-[state=inactive]:bg-transparent
      data-[state=inactive]:text-neutral-700
      dark:data-[state=inactive]:text-neutral-200

      [&_svg]:pointer-events-none
      [&_svg]:shrink-0
      [&_svg:not([class*='size-'])]:size-4
      `,
      props.class,
    )"
      v-bind="forwardedProps"
  >
    <slot />
  </TabsTrigger>
</template>
