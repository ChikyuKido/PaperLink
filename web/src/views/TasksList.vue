<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { Activity, Clock, Eye, RefreshCcw, Square } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {apiFetch} from "@/auth/api.ts";

type TaskStatus = "RUNNING" | "FAILED" | "COMPLETED" | "STOPPED"

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
type ListTasksResponse = {
  tasks: Task[]
}

const router = useRouter()
const tasks = ref<Task[]>([])
const isLoading = ref(false)

function statusVariant(s: TaskStatus) {
  if (s === "RUNNING") return "default"
  if (s === "COMPLETED") return "success"
  if (s === "STOPPED") return "secondary"
  return "destructive"
}

function toMillis(ts: number) {
  return ts < 1_000_000_000_000 ? ts * 1000 : ts
}

function duration(t: Task) {
  const start = toMillis(t.startTime)
  const end = t.endTime ? toMillis(t.endTime) : Date.now()
  const sec = Math.max(0, Math.floor((end - start) / 1000))
  return `${Math.floor(sec / 60)}m ${sec % 60}s`
}

async function stopTask(taskId: string) {
  const res = await apiFetch(`/api/v1/task/stop/${taskId}`, { method: "POST" })
  if (!res.ok) return
  await loadTasks()
}

async function loadTasks() {
  isLoading.value = true
  try {
    const r = await apiFetch("/api/v1/task/list")
    const j: ApiResponse<ListTasksResponse> = await r.json()
    if (!j?.data?.tasks || j.data.tasks.length === 0) {
      tasks.value = []
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
    <div class="rounded-2xl border bg-white dark:bg-neutral-900 overflow-hidden">
      <div v-if="tasks.length === 0" class="px-6 py-8 text-sm text-neutral-500">
        No tasks found.
      </div>
      <div v-else class="divide-y">
        <div v-for="t in tasks" :key="t.id" class="px-4 py-3 flex items-center gap-3">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <p class="font-medium truncate">{{ t.name }}</p>
              <Badge :variant="statusVariant(t.status)">
                {{ t.status }}
              </Badge>
            </div>
            <div class="mt-1 flex items-center gap-3 text-xs text-neutral-500">
              <span>Started: {{ new Date(toMillis(t.startTime)).toLocaleString() }}</span>
              <span class="inline-flex items-center gap-1">
                <Clock class="h-3.5 w-3.5" />
                {{ duration(t) }}
              </span>
            </div>
          </div>

          <div class="flex items-center gap-2 shrink-0">
            <Button
              v-if="t.status === 'RUNNING'"
              size="sm"
              variant="destructive"
              @click="stopTask(t.id)"
            >
              <Square class="h-4 w-4 mr-1" />
              Stop
            </Button>
            <Button
              size="sm"
              variant="secondary"
              @click="router.push(`/admin/task/${t.id}`)"
            >
              <Eye class="h-4 w-4 mr-1" />
              View
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
