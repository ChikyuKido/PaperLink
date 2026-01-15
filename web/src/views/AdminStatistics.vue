<script setup lang="ts">
import { onMounted, ref } from "vue"
import { BarChart3, RefreshCcw } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

import { getAdminStats, type AdminStats } from "@/lib/admin_api"

type Notice = { type: "success" | "error"; message: string } | null

const notice = ref<Notice>(null)
let noticeTimer: number | undefined
function showNotice(n: Notice) {
  notice.value = n
  if (noticeTimer) window.clearTimeout(noticeTimer)
  if (n) noticeTimer = window.setTimeout(() => (notice.value = null), 3500)
}

const stats = ref<AdminStats | null>(null)
const loading = ref(false)

function formatBytes(bytes: number) {
  if (!Number.isFinite(bytes) || bytes <= 0) return "0 B"
  const units = ["B", "KB", "MB", "GB", "TB"]
  let i = 0
  let v = bytes
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

async function load() {
  loading.value = true
  try {
    stats.value = await getAdminStats()
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to load statistics" })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
    <section
      class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
        class="px-4 sm:px-6 py-4 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
      >
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-3">
            <div
              class="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-emerald-600/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200"
            >
              <BarChart3 class="h-5 w-5" />
            </div>
            <div>
              <h1 class="text-lg font-semibold tracking-tight">Statistics</h1>
              <p class="text-xs text-neutral-500 dark:text-neutral-400">Live overview of your instance.</p>
            </div>
          </div>

          <Button variant="outline" class="rounded-full" :disabled="loading" @click="load">
            <RefreshCcw class="h-4 w-4" />
            Refresh
          </Button>
        </div>

        <div
          v-if="notice"
          class="mt-4 rounded-xl border px-4 py-3 text-sm"
          :class="
            notice.type === 'success'
              ? 'border-emerald-600/30 bg-emerald-600/10 text-emerald-900 dark:text-emerald-200 dark:bg-emerald-500/10'
              : 'border-red-600/30 bg-red-600/10 text-red-900 dark:text-red-200 dark:bg-red-500/10'
          "
        >
          {{ notice.message }}
        </div>
      </div>
    </section>

    <Card class="border border-neutral-200 dark:border-neutral-800">
      <CardHeader>
        <CardTitle class="text-sm">Workspace statistics</CardTitle>
        <CardDescription class="text-[11px]">Key metrics for your Paperlink instance.</CardDescription>
      </CardHeader>

      <CardContent>
        <div v-if="loading" class="text-sm text-neutral-600 dark:text-neutral-300">Loading…</div>

        <div
          v-else-if="!stats"
          class="rounded-xl border border-dashed border-neutral-300 bg-neutral-50 p-4 text-sm text-neutral-600 dark:border-neutral-700 dark:bg-neutral-900/40 dark:text-neutral-300"
        >
          No statistics available.
        </div>

        <div v-else class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">Users</p>
            <p class="mt-1 text-2xl font-semibold">{{ stats.userCount }}</p>
          </div>

          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">Documents</p>
            <p class="mt-1 text-2xl font-semibold">{{ stats.documentCount }}</p>
          </div>

          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">Total size</p>
            <p class="mt-1 text-2xl font-semibold">{{ formatBytes(stats.totalDocSize) }}</p>
          </div>

          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">Total pages</p>
            <p class="mt-1 text-2xl font-semibold">{{ stats.totalPages }}</p>
          </div>

          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">D4S books</p>
            <p class="mt-1 text-2xl font-semibold">{{ stats.d4sBookCount }}</p>
          </div>

          <div class="rounded-2xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900">
            <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">D4S accounts</p>
            <p class="mt-1 text-2xl font-semibold">{{ stats.d4sAccountCount }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
