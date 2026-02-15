<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { FileText, Loader2 } from 'lucide-vue-next'
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

type SubmitPayload = {
  name: string
  description: string
  tags: string[]
  fileUUID: string
  directoryId: number | null
  folderPath: string
}

const props = defineProps<{ open: boolean; directoryId: number | null; folderPath?: string }>()

const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: SubmitPayload): void }>()

const fileInput = ref<HTMLInputElement | null>(null)

const file = ref<File | null>(null)
const fileName = ref('')
const fileUUID = ref<string | null>(null)

// Truncate very long names for display (avoid dialog overflow)
const displayFileName = computed(() => {
  const n = fileName.value || ''
  const max = 42
  if (n.length <= max) return n
  return n.slice(0, max - 1) + '…'
})

const name = ref('')
const description = ref('')
const tagsRaw = ref('')

const uploading = ref(false)
const error = ref<string | null>(null)

watch(() => props.open, (open) => {
  if (!open) return
  file.value = null
  fileName.value = ''
  fileUUID.value = null
  name.value = ''
  description.value = ''
  tagsRaw.value = ''
  error.value = null
  uploading.value = false
})

function close() {
  if (!uploading.value) emit('close')
}

function pickFile() {
  if (!uploading.value) fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const selected = input.files?.[0]
  error.value = null

  if (!selected || !selected.name.toLowerCase().endsWith('.pdf')) {
    error.value = 'Only PDF files are allowed.'
    input.value = ''
    return
  }

  file.value = selected
  fileName.value = selected.name
  await upload(selected)
}

async function upload(selectedFile: File) {
  uploading.value = true

  try {
    const formData = new FormData()
    formData.append('file', selectedFile)

    const res = await apiFetch('/api/v1/document/upload', {
      method: 'POST',
      body: formData,
    })

    const json = await res.json()
    if (!res.ok || json?.code !== 200 || !json?.data?.fileUUID) {
      error.value = 'Failed to upload file.'
      file.value = null
      fileName.value = ''
      fileUUID.value = null
      return
    }

    fileUUID.value = json.data.fileUUID
  } catch {
    error.value = 'Failed to upload file.'
    file.value = null
    fileName.value = ''
    fileUUID.value = null
  } finally {
    uploading.value = false
  }
}

const documentName = computed(() =>
    name.value.trim() ||
    fileName.value.replace(/\.pdf$/i, '') ||
    'Untitled'
)

async function save() {
  if (!fileUUID.value) {
    error.value = 'File not uploaded.'
    return
  }

  uploading.value = true
  error.value = null

  try {
    const payload = {
      name: documentName.value,
      description: description.value,
      directoryId: props.directoryId ?? null,
      folderPath: props.folderPath ?? '',
      tags: tagsRaw.value.split(',').map(t => t.trim()).filter(Boolean),
      fileUUID: fileUUID.value,
    }

    const res = await apiFetch('/api/v1/document/create', {
      method: 'POST',
      body: JSON.stringify(payload),
    })

    if (!res.ok) {
      error.value = 'Failed to create document.'
      return
    }
    const resJson = await res.json()


    emit('submit', {
      name: payload.name,
      description: payload.description,
      tags: payload.tags,
      fileUUID: resJson.data.uuid,
      directoryId: payload.directoryId,
      folderPath: payload.folderPath,
    })

    emit('close')
  } catch {
    error.value = 'Failed to create document.'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="close">
    <DialogContent class="max-w-lg rounded-2xl border-neutral-200 dark:border-neutral-800">
      <DialogHeader>
        <div class="flex items-start gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl border border-emerald-600/30 bg-emerald-600/10 text-emerald-800 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-200">
            <FileText class="h-5 w-5" aria-hidden="true" />
          </div>
          <div class="min-w-0">
            <DialogTitle class="leading-tight">Create document</DialogTitle>
            <DialogDescription class="mt-0.5">
              Upload a PDF and save it <span class="font-medium text-neutral-900 dark:text-neutral-50">{{ (folderPath ?? '').length ? folderPath : 'Home' }}</span>.
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div class="space-y-5">
        <div class="space-y-1.5">
          <label class="text-xs font-medium text-neutral-700 dark:text-neutral-200">PDF file</label>
          <div class="flex items-center justify-between gap-3 rounded-xl border border-dashed border-neutral-300 bg-neutral-50/60 px-3 py-2.5 dark:border-neutral-700 dark:bg-neutral-950/40">
            <div class="flex items-center gap-3 min-w-0">
              <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-white ring-1 ring-neutral-900/10 dark:bg-neutral-100 dark:text-neutral-900">
                <Loader2 v-if="uploading" class="h-4 w-4 animate-spin" />
                <FileText v-else class="h-4 w-4" />
              </div>
              <div class="min-w-0">
                <p class="truncate text-sm font-medium">{{ displayFileName || 'No file selected' }}</p>
                <p class="truncate text-xs text-neutral-500 dark:text-neutral-400">
                  {{ fileUUID ? 'Uploaded — ready to save' : 'Choose a PDF to upload' }}
                </p>
              </div>
            </div>

            <Button size="sm" variant="outline" class="rounded-xl" :disabled="uploading" @click="pickFile">
              Choose
            </Button>
          </div>

          <input
            ref="fileInput"
            type="file"
            accept=".pdf"
            class="hidden"
            @change="onFileChange"
          />
        </div>

        <div class="grid grid-cols-1 gap-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-neutral-700 dark:text-neutral-200">Name</label>
            <input
              v-model="name"
              :disabled="uploading"
              class="w-full rounded-xl border border-neutral-200 bg-white px-3 py-2 text-sm shadow-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-emerald-400 dark:focus:ring-emerald-400/20"
              placeholder="e.g. Contract 2025"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-medium text-neutral-700 dark:text-neutral-200">Description</label>
            <textarea
              v-model="description"
              rows="3"
              :disabled="uploading"
              class="w-full rounded-xl border border-neutral-200 bg-white px-3 py-2 text-sm shadow-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-emerald-400 dark:focus:ring-emerald-400/20"
              placeholder="Optional: a short description"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-medium text-neutral-700 dark:text-neutral-200">Tags</label>
            <input
              v-model="tagsRaw"
              :disabled="uploading"
              class="w-full rounded-xl border border-neutral-200 bg-white px-3 py-2 text-sm shadow-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-emerald-400 dark:focus:ring-emerald-400/20"
              placeholder="Comma-separated (e.g. finance, 2025)"
            />
          </div>
        </div>

        <p v-if="error" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-200">
          {{ error }}
        </p>
      </div>

      <DialogFooter class="gap-2">
        <Button variant="outline" size="sm" :disabled="uploading" class="rounded-xl" @click="emit('close')">
          Cancel
        </Button>
        <Button
          size="sm"
          :disabled="uploading || !fileUUID"
          class="rounded-xl bg-emerald-600 text-white hover:bg-emerald-700 dark:bg-emerald-500 dark:hover:bg-emerald-400"
          @click="save"
        >
          Save document
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
