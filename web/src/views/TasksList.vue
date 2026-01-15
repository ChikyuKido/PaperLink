<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { Activity, Clock, Eye, RefreshCcw } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import {apiFetch} from "@/auth/api.ts";

type TaskStatus = "RUNNING" | "FAILED" | "COMPLETED"

type Task = {
  id: string
  name: string
  status: TaskStatus
  startTime: number
  endTime: number
}

type ApiResponse<T> = {
  code: number
  data: T
}

const router = useRouter()
const tasks = ref<Task[]>([])
const isLoading = ref(false)

function statusVariant(s: TaskStatus) {
  if (s === "RUNNING") return "default"
  if (s === "COMPLETED") return "success"
  return "destructive"
}

function duration(t: Task) {
  const end = t.endTime || Date.now()
  const sec = Math.max(0, Math.floor((end - t.startTime) / 1000))
  return `${Math.floor(sec / 60)}m ${sec % 60}s`
}

async function loadTasks() {
  isLoading.value = true
  try {
    const r = await apiFetch("/api/v1/task/list")
    const j: ApiResponse<Task[]> = await r.json()
    if (j.data.tasks.length === 0) {
      return
    }
    tasks.value = j.data.tasks
  } finally {
    isLoading.value = false
  }
}

onMounted(loadTasks)
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-6 space-y-4">
    <!-- Header -->
    <section class="rounded-2xl border bg-white dark:bg-neutral-900">
      <div class="px-6 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 rounded-2xl bg-emerald-600/10 text-emerald-700 flex items-center justify-center">
            <Activity class="h-5 w-5" />
          </div>
          <div>
            <h1 class="text-lg font-semibold">Tasks</h1>
            <p class="text-xs text-neutral-500">Background executions</p>
          </div>
        </div>

        <Button variant="outline" class="rounded-full" @click="loadTasks">
          <RefreshCcw class="h-4 w-4 mr-1" />
          Refresh
        </Button>
      </div>
    </section>

    <!-- List -->
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="t in tasks" :key="t.id">
        <CardContent class="p-4 space-y-3">
          <div class="flex justify-between items-start gap-2">
            <div>
              <p class="font-semibold">{{ t.name }}</p>
              <p class="text-xs text-neutral-500">
                {{ new Date(t.startTime).toLocaleString() }}
              </p>
            </div>
            <Badge :variant="statusVariant(t.status)">
              {{ t.status }}
            </Badge>
          </div>

          <div class="flex items-center gap-2 text-xs text-neutral-500">
            <Clock class="h-4 w-4" />
            {{ duration(t) }}
          </div>

          <Separator />

          <Button
              class="w-full"
              variant="secondary"
              @click="router.push(`/admin/task/${t.id}`)"
          >
            <Eye class="h-4 w-4 mr-1" />
            View
          </Button>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
