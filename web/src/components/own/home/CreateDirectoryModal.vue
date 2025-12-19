<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { FolderPlus, Loader2 } from 'lucide-vue-next'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { apiFetch } from '@/auth/api'

const props = defineProps<{ open: boolean; parentId: number | null; folderPath?: string }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', payload: { id: number; name: string; parentId: number | null; folderPath: string }): void
}>()

const name = ref('')
const creating = ref(false)
const error = ref<string | null>(null)

const locationLabel = computed(() => (props.folderPath ?? '') || 'Home')

watch(() => props.open, (open) => {
  if (!open) return
  name.value = ''
  creating.value = false
  error.value = null
})

function close() {
  if (!creating.value) emit('close')
}

async function create() {
  const trimmed = name.value.trim()
  if (!trimmed) {
    error.value = 'Please enter a name.'
    return
  }

  creating.value = true
  error.value = null

  try {
    const payload = {
      name: trimmed,
      parentId: props.parentId ?? null,
      folderPath: props.folderPath ?? '',
    }

    const res = await apiFetch('/api/v1/directory/create', {
      method: 'POST',
      body: JSON.stringify(payload),
    })

    const json = await res.json().catch(() => null)
    const createdId = Number(json?.data?.id)

    if (!res.ok || !Number.isFinite(createdId)) {
      error.value = 'Failed to create directory.'
      return
    }

    emit('created', { id: createdId, name: trimmed, parentId: payload.parentId, folderPath: payload.folderPath })
    emit('close')
  } catch {
    error.value = 'Failed to create directory.'
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="close">
    <DialogContent class="max-w-md rounded-2xl border-neutral-200 dark:border-neutral-800">
      <DialogHeader>
        <div class="flex items-start gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl border border-emerald-600/30 bg-emerald-600/10 text-emerald-800 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-200">
            <FolderPlus class="h-5 w-5" aria-hidden="true" />
          </div>
          <div class="min-w-0">
            <DialogTitle class="leading-tight">Create directory</DialogTitle>
            <DialogDescription class="mt-0.5">
              Create a new folder in <span class="font-medium text-neutral-900 dark:text-neutral-50">{{ locationLabel }}</span>.
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-1.5">
          <label class="text-xs font-medium text-neutral-700 dark:text-neutral-200">Name</label>
          <div class="relative">
            <FolderPlus class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-emerald-700/80 dark:text-emerald-300/80" />
            <input
              v-model="name"
              :disabled="creating"
              class="w-full rounded-xl border border-neutral-200 bg-white px-3 py-2 pl-9 text-sm shadow-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-emerald-400 dark:focus:ring-emerald-400/20"
              placeholder="e.g. Invoices"
              @keydown.enter.prevent="create"
              autofocus
            />
          </div>
        </div>

        <p v-if="error" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-200">
          {{ error }}
        </p>
      </div>

      <DialogFooter class="gap-2">
        <Button variant="outline" size="sm" :disabled="creating" class="rounded-xl" @click="emit('close')">
          Cancel
        </Button>
        <Button size="sm" :disabled="creating" class="rounded-xl bg-emerald-600 text-white hover:bg-emerald-700 dark:bg-emerald-500 dark:hover:bg-emerald-400" @click="create">
          <Loader2 v-if="creating" class="mr-2 h-4 w-4 animate-spin" />
          Create
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<style scoped>
</style>
