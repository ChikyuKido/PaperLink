<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { ArrowLeft, RefreshCcw, Square } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { Card } from "@/components/ui/card"
import {apiFetch} from "@/auth/api.ts";

type TaskStatus = "RUNNING" | "FAILED" | "COMPLETED" | "STOPPED"

type TaskDetail = {
  id: string
  name: string
  status: TaskStatus
  startTime: number
  endTime: number
  content: string[]
}

type ApiResponse<T> = {
  code: number
  data: T
}

const route = useRoute()
const router = useRouter()
const task = ref<TaskDetail | null>(null)
let pollTimer: number | undefined

function statusVariant(s: TaskStatus) {
  if (s === "RUNNING") return "default"
  if (s === "COMPLETED") return "success"
  if (s === "STOPPED") return "secondary"
  return "destructive"
}

function toMillis(ts: number) {
  return ts < 1_000_000_000_000 ? ts * 1000 : ts
}

function duration(task: TaskDetail) {
  const start = toMillis(task.startTime)
  const end = task.endTime ? toMillis(task.endTime) : Date.now()
  const sec = Math.max(0, Math.floor((end - start) / 1000))
  return `${Math.floor(sec / 60)}m ${sec % 60}s`
}

async function stopTask() {
  if (!task.value || task.value.status !== "RUNNING") return
  const res = await apiFetch(`/api/v1/task/stop/${task.value.id}`, { method: "POST" })
  if (!res.ok) return
  await loadTask()
}

async function loadTask() {
  const r = await apiFetch(`/api/v1/task/view/${route.params.id}`)
  const j: ApiResponse<TaskDetail> = await r.json()
  task.value = j.data
}

onMounted(async () => {
  await loadTask()
  pollTimer = window.setInterval(async () => {
    if (task.value?.status === "RUNNING") {
      await loadTask()
    }
  }, 2000)
})

onUnmounted(() => pollTimer && clearInterval(pollTimer))
</script>

<template>
  <div class="mx-auto max-w-5xl px-4 lg:px-6 py-6 space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <Button variant="ghost" @click="router.push('/admin/tasks')">
        <ArrowLeft class="h-4 w-4 mr-1" />
        Back
      </Button>

      <div class="flex items-center gap-2">
        <Button
          v-if="task?.status === 'RUNNING'"
          variant="destructive"
          @click="stopTask"
        >
          <Square class="h-4 w-4 mr-1" />
          Stop
        </Button>
        <Button variant="outline" @click="loadTask">
          <RefreshCcw class="h-4 w-4 mr-1" />
          Refresh
        </Button>
      </div>
    </div>

    <Card class="p-4 space-y-2">
      <div class="flex items-center gap-3">
        <h1 class="text-lg font-semibold">{{ task?.name }}</h1>
        <Badge v-if="task" :variant="statusVariant(task.status)">
          {{ task.status }}
        </Badge>
      </div>
      <p class="text-xs text-neutral-500">
        Started: {{ task && new Date(toMillis(task.startTime)).toLocaleString() }}
      </p>
      <p v-if="task" class="text-xs text-neutral-500">
        Duration: {{ duration(task) }}
      </p>
    </Card>

    <Separator />

    <!-- Log Viewer -->
    <div
        class="h-[65vh] overflow-auto rounded-xl bg-black text-emerald-400 p-4 font-mono text-xs space-y-1"
    >
      <div
          v-for="(line, i) in task?.content"
          :key="i"
          class="whitespace-pre-wrap"
      >
        {{ line }}
      </div>

      <div v-if="!task?.content.length" class="text-neutral-500">
        No logs yet.
      </div>
    </div>
  </div>
</template>
