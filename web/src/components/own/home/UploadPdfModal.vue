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

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ (e: 'close'): void }>()

const fileInput = ref<HTMLInputElement | null>(null)

const file = ref<File | null>(null)
const fileName = ref('')
const fileUUID = ref<string | null>(null)

const name = ref('')
const description = ref('')
const path = ref('')
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
  path.value = ''
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

async function upload(file: File) {
  uploading.value = true

  try {
    const formData = new FormData()
    formData.append('file', file)

    const res = await apiFetch('/api/v1/document/upload', {
      method: 'POST',
      body: formData,
    })

    const json = await res.json()
    if (!res.ok || json?.code !== 200 || !json?.data?.id) {
      throw new Error()
    }

    fileUUID.value = json.data.id
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
      path: path.value,
      tags: tagsRaw.value.split(',').map(t => t.trim()).filter(Boolean),
      fileUUID: fileUUID.value,
    }

    const res = await apiFetch('/api/v1/document/create', {
      method: 'POST',
      body: JSON.stringify(payload),
    })

    if (!res.ok) throw new Error()
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
    <DialogContent class="max-w-lg rounded-2xl">
      <DialogHeader>
        <DialogTitle>Upload PDF</DialogTitle>
        <DialogDescription>File uploads immediately. Metadata is saved after.</DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div>
          <label class="text-xs font-medium">PDF file</label>
          <div class="flex items-center justify-between rounded-xl border border-dashed px-3 py-2">
            <div class="flex items-center gap-2 min-w-0">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-neutral-900 text-white">
                <Loader2 v-if="uploading" class="h-4 w-4 animate-spin" />
                <FileText v-else class="h-4 w-4" />
              </div>
              <p class="truncate text-xs">
                {{ fileName || 'No file selected' }}
              </p>
            </div>
            <Button size="sm" variant="outline" :disabled="uploading" @click="pickFile">
              Choose file
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

        <div>
          <label class="text-xs font-medium">Name</label>
          <input v-model="name" :disabled="uploading" class="w-full rounded-xl border px-3 py-2 text-sm" />
        </div>

        <div>
          <label class="text-xs font-medium">Description</label>
          <textarea v-model="description" rows="3" :disabled="uploading" class="w-full rounded-xl border px-3 py-2 text-sm" />
        </div>

        <div>
          <label class="text-xs font-medium">Path</label>
          <input v-model="path" :disabled="uploading" class="w-full rounded-xl border px-3 py-2 text-sm" />
        </div>

        <div>
          <label class="text-xs font-medium">Tags</label>
          <input v-model="tagsRaw" :disabled="uploading" class="w-full rounded-xl border px-3 py-2 text-sm" />
        </div>

        <p v-if="error" class="text-xs text-red-600">{{ error }}</p>
      </div>

      <DialogFooter>
        <Button variant="outline" size="sm" :disabled="uploading" @click="emit('close')">
          Cancel
        </Button>
        <Button size="sm" :disabled="uploading || !fileUUID" @click="save">
          Save document
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
