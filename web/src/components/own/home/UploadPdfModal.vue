<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { FileText } from 'lucide-vue-next'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', payload: {
    name: string
    tag: string
    description: string
    file: File
    createdAt: string | null
    modifiedAt: string | null
  }): void
}>()

function pad2(n: number): string {
  return n < 10 ? `0${n}` : `${n}`
}

function formatDateForInput(date: Date): string {
  const year = date.getFullYear()
  const month = pad2(date.getMonth() + 1)
  const day = pad2(date.getDate())
  const hours = pad2(date.getHours())
  const minutes = pad2(date.getMinutes())
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

const uploadName = ref('')
const uploadTag = ref('')
const uploadDescription = ref('')
const uploadFile = ref<File | null>(null)
const uploadFileName = ref('')
const uploadError = ref<string | null>(null)
const uploadCreatedAt = ref<string>('')
const uploadModifiedAt = ref<string>('')

const fileInput = ref<HTMLInputElement | null>(null)

function resetUploadState() {
  uploadName.value = ''
  uploadTag.value = ''
  uploadDescription.value = ''
  uploadFile.value = null
  uploadFileName.value = ''
  uploadError.value = null
  uploadCreatedAt.value = ''
  uploadModifiedAt.value = ''
}

watch(
    () => props.open,
    (val) => {
      if (val) resetUploadState()
    }
)

function handleDialogOpenChange(value: boolean) {
  if (!value) {
    emit('close')
  }
}

function openFilePicker() {
  fileInput.value?.click()
}

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const files = target.files
  uploadError.value = null

  if (!files || files.length === 0) {
    uploadFile.value = null
    uploadFileName.value = ''
    uploadCreatedAt.value = ''
    uploadModifiedAt.value = ''
    return
  }

  const file = files[0]
  const isPdfByType = file.type === 'application/pdf'
  const isPdfByName = file.name.toLowerCase().endsWith('.pdf')

  if (!isPdfByType && !isPdfByName) {
    uploadError.value = 'Only PDF files are allowed.'
    uploadFile.value = null
    uploadFileName.value = ''
    uploadCreatedAt.value = ''
    uploadModifiedAt.value = ''
    target.value = ''
    return
  }

  uploadFile.value = file
  uploadFileName.value = file.name

  const lastModifiedDate = new Date(file.lastModified)
  const formatted = formatDateForInput(lastModifiedDate)
  if (!uploadCreatedAt.value) uploadCreatedAt.value = formatted
  if (!uploadModifiedAt.value) uploadModifiedAt.value = formatted
}

const computedDocumentName = computed(() => {
  const trimmed = uploadName.value.trim()
  if (trimmed) return trimmed
  if (uploadFileName.value) {
    return uploadFileName.value.replace(/\.pdf$/i, '')
  }
  return 'Untitled'
})

function submitUpload() {
  uploadError.value = null

  if (!uploadFile.value) {
    uploadError.value = 'Please select a PDF file to upload.'
    return
  }

  const isPdfByType = uploadFile.value.type === 'application/pdf'
  const isPdfByName = uploadFile.value.name.toLowerCase().endsWith('.pdf')
  if (!isPdfByType && !isPdfByName) {
    uploadError.value = 'Only PDF files are allowed.'
    return
  }

  emit('submit', {
    name: computedDocumentName.value,
    tag: uploadTag.value,
    description: uploadDescription.value,
    file: uploadFile.value,
    createdAt: uploadCreatedAt.value || null,
    modifiedAt: uploadModifiedAt.value || null,
  })
}
</script>

<template>
  <Dialog :open="open" @update:open="handleDialogOpenChange">
    <DialogContent class="max-w-lg rounded-2xl border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900">
      <DialogHeader class="px-0 pb-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-700/10 border border-emerald-700/40 dark:bg-emerald-500/10 dark:border-emerald-500/50">
              <FileText class="w-4 h-4 text-emerald-800 dark:text-emerald-300" />
            </div>
            <div>
              <DialogTitle class="text-sm font-medium text-neutral-900 dark:text-neutral-50">
                Upload PDF
              </DialogTitle>
              <DialogDescription class="text-xs text-neutral-500 dark:text-neutral-400">
                Add a new document with metadata and tags.
              </DialogDescription>
            </div>
          </div>
        </div>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-1.5">
          <label class="block text-xs font-medium text-neutral-700 dark:text-neutral-200">
            Name
          </label>
          <input
              v-model="uploadName"
              type="text"
              placeholder="Proposal for Q2 Roadmap"
              class="w-full rounded-xl border border-neutral-300 bg-neutral-50 px-3 py-2 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-50"
          />
          <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
            The document will be saved as
            <span class="font-mono">{{ computedDocumentName }}.pdf</span>
          </p>
        </div>

        <div class="space-y-1.5">
          <label class="block text-xs font-medium text-neutral-700 dark:text-neutral-200">
            Tag
          </label>
          <input
              v-model="uploadTag"
              type="text"
              placeholder="project-x, contract, finance"
              class="w-full rounded-xl border border-neutral-300 bg-neutral-50 px-3 py-2 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-50"
          />
          <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
            Use tags to quickly filter and search later.
          </p>
        </div>

        <div class="space-y-2">
          <p class="text-xs font-medium text-neutral-700 dark:text-neutral-200">
            Document metadata
          </p>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="block text-[11px] font-medium text-neutral-600 dark:text-neutral-300">
                Created at
              </label>
              <input
                  v-model="uploadCreatedAt"
                  type="datetime-local"
                  class="w-full rounded-xl border border-neutral-300 bg-neutral-50 px-2.5 py-2 text-xs text-neutral-900 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-50"
              />
            </div>
            <div class="space-y-1">
              <label class="block text-[11px] font-medium text-neutral-600 dark:text-neutral-300">
                Modified at
              </label>
              <input
                  v-model="uploadModifiedAt"
                  type="datetime-local"
                  class="w-full rounded-xl border border-neutral-300 bg-neutral-50 px-2.5 py-2 text-xs text-neutral-900 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-50"
              />
            </div>
          </div>
          <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
            Pre-filled from the file where possible. You can adjust these before uploading.
          </p>
        </div>

        <div class="space-y-1.5">
          <label class="block text-xs font-medium text-neutral-700 dark:text-neutral-200">
            Description / notes
          </label>
          <textarea
              v-model="uploadDescription"
              rows="3"
              placeholder="Short description, relevant people, version info, etc."
              class="w-full rounded-xl border border-neutral-300 bg-neutral-50 px-3 py-2 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-50"
          />
        </div>

        <div class="space-y-2">
          <label class="block text-xs font-medium text-neutral-700 dark:text-neutral-200">
            PDF file
          </label>
          <div class="flex items-center justify-between rounded-xl border border-dashed border-neutral-300 bg-neutral-50 px-3 py-2.5 text-sm dark:border-neutral-700 dark:bg-neutral-900/70">
            <div class="flex min-w-0 items-center gap-2">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-neutral-900 text-neutral-50 dark:bg-neutral-200 dark:text-neutral-900">
                <FileText class="w-4 h-4" aria-hidden="true" />
              </div>
              <div class="min-w-0">
                <p class="truncate text-xs text-neutral-700 dark:text-neutral-100">
                  {{ uploadFileName || 'No file selected' }}
                </p>
                <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
                  Only PDF files are supported.
                </p>
              </div>
            </div>
            <Button
                variant="outline"
                size="sm"
                class="rounded-full px-3 py-1.5 text-xs font-medium"
                type="button"
                @click="openFilePicker"
            >
              Choose file
            </Button>
          </div>

          <input
              ref="fileInput"
              type="file"
              accept="application/pdf,.pdf"
              class="hidden"
              @change="onFileChange"
          />

          <p v-if="uploadError" class="text-xs text-red-600 dark:text-red-400">
            {{ uploadError }}
          </p>
        </div>
      </div>

      <DialogFooter class="pt-4">
        <Button
            variant="outline"
            size="sm"
            class="rounded-full px-3.5"
            type="button"
            @click="emit('close')"
        >
          Cancel
        </Button>
        <Button
            size="sm"
            class="rounded-full px-4"
            type="button"
            @click="submitUpload"
        >
          Upload PDF
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<style scoped>
</style>
