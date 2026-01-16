<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { ArrowLeft, RefreshCcw } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { Card } from "@/components/ui/card"
import {apiFetch} from "@/auth/api.ts";

type TaskStatus = "RUNNING" | "FAILED" | "COMPLETED"

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
  return "destructive"
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
      <Button variant="ghost" @click="router.push('/tasks')">
        <ArrowLeft class="h-4 w-4 mr-1" />
        Back
      </Button>

      <Button variant="outline" @click="loadTask">
        <RefreshCcw class="h-4 w-4 mr-1" />
        Refresh
      </Button>
    </div>

    <Card class="p-4 space-y-2">
      <div class="flex items-center gap-3">
        <h1 class="text-lg font-semibold">{{ task?.name }}</h1>
        <Badge v-if="task" :variant="statusVariant(task.status)">
          {{ task.status }}
        </Badge>
      </div>
      <p class="text-xs text-neutral-500">
        Started: {{ task && new Date(task.startTime).toLocaleString() }}
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
