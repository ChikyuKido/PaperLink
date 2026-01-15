<script setup lang="ts">
import { Shield, SlidersHorizontal } from "lucide-vue-next"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

</script>

<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
    <section
      class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
        class="px-4 sm:px-6 py-4 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
      >
        <div class="flex items-center gap-3">
          <div
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-emerald-600/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200"
          >
            <Shield class="h-5 w-5" />
          </div>
          <div>
            <h1 class="text-lg font-semibold tracking-tight">Admin Settings</h1>
            <p class="text-xs text-neutral-500 dark:text-neutral-400">System-wide settings for your Paperlink instance.</p>
          </div>
        </div>
      </div>
    </section>

    <Card class="border border-neutral-200 dark:border-neutral-800">
      <CardHeader>
        <div class="flex items-center gap-2">
          <span class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-emerald-700/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200">
            <SlidersHorizontal class="h-4 w-4" />
          </span>
          <div>
            <CardTitle class="text-sm">Configuration</CardTitle>
            <CardDescription class="text-[11px]">Placeholder for upcoming admin settings.</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div class="rounded-xl border border-dashed border-neutral-300 bg-neutral-50 p-4 text-sm text-neutral-600 dark:border-neutral-700 dark:bg-neutral-900/40 dark:text-neutral-300">
          Coming soon.
        </div>
      </CardContent>
    </Card>
  </div>
</template>

