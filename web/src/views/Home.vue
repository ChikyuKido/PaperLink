<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import {
  Folder,
  FileText,
  ChevronRight,
  ArrowLeft,
  ArrowRight,
  Home as HomeIcon,
  BarChart3,
  Plus,
  X,
} from 'lucide-vue-next'
import type { Item } from '@/dto/item'

const initialTree: Item[] = [
  {
    id: '1',
    name: 'Projects',
    type: 'folder',
    children: [
      { id: '1-1', name: 'Proposal.pdf', type: 'file', size: 320_000 },
      {
        id: '1-2',
        name: 'Specs',
        type: 'folder',
        children: [
          { id: '1-2-1', name: 'Architecture.pdf', type: 'file', size: 2_400_000 },
        ],
      },
    ],
  },
  {
    id: '2',
    name: 'Invoices',
    type: 'folder',
    children: [
      { id: '2-1', name: 'Invoice-2025.pdf', type: 'file', size: 180_000 },
    ],
  },
  { id: '3', name: 'Readme.pdf', type: 'file', size: 64_000 },
]

const tree = ref(initialTree) as Ref<Item[]>
const path = ref<Item[]>([])
const history = ref<Item[][]>([[]])
const historyIndex = ref(0)

const currentItems = computed(() => {
  const last = path.value[path.value.length - 1]
  if (!last) return tree.value
  return last.children ?? []
})

const currentFolderLabel = computed(() => {
  const last = path.value[path.value.length - 1]
  return last ? last.name : 'All documents'
})

function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let idx = 0
  while (value >= 1024 && idx < units.length - 1) {
    value /= 1024
    idx++
  }
  const decimals = value >= 10 || idx === 0 ? 0 : 1
  return `${value.toFixed(decimals)} ${units[idx]}`
}

const libraryStats = computed(() => {
  const result = { folders: 0, files: 0, bytes: 0 }

  const traverse = (items: Item[]) => {
    for (const item of items) {
      if (item.type === 'folder') {
        result.folders++
        if (item.children) traverse(item.children)
      } else {
        result.files++
        if (typeof item.size === 'number') {
          result.bytes += item.size
        }
      }
    }
  }

  traverse(tree.value)
  return result
})

const currentLevelStats = computed(() => {
  let folders = 0
  let files = 0
  let bytes = 0

  for (const item of currentItems.value) {
    if (item.type === 'folder') {
      folders++
    } else {
      files++
      if (typeof item.size === 'number') {
        bytes += item.size
      }
    }
  }

  return { folders, files, bytes }
})

function updatePath(newPath: Item[], pushToHistory = true) {
  path.value = newPath
  if (!pushToHistory) return
  history.value = history.value.slice(0, historyIndex.value + 1)
  history.value.push(newPath)
  historyIndex.value++
}

function enterItem(item: Item) {
  if (item.type === 'folder') {
    updatePath([...path.value, item])
  } else {
    openFile(item)
  }
}

function breadcrumbClick(index: number) {
  const newPath = index < 0 ? [] : path.value.slice(0, index + 1)
  updatePath(newPath)
}

function goHome() {
  if (path.value.length === 0) return
  updatePath([])
}

function goBack() {
  if (historyIndex.value === 0) return
  historyIndex.value--
  path.value = history.value[historyIndex.value]
}

function goForward() {
  if (historyIndex.value >= history.value.length - 1) return
  historyIndex.value++
  path.value = history.value[historyIndex.value]
}

function openFile(item: Item) {
  console.log('Open file:', item.name)
}

function iconFor(item: Item) {
  return item.type === 'folder' ? Folder : FileText
}

/**
 * Upload modal + file handling
 */
const isUploadModalOpen = ref(false)
const uploadName = ref('')
const uploadTag = ref('')
const uploadDescription = ref('')
const uploadFile = ref<File | null>(null)
const uploadFileName = ref('')
const uploadError = ref<string | null>(null)

const fileInput = ref<HTMLInputElement | null>(null)

function resetUploadState() {
  uploadName.value = ''
  uploadTag.value = ''
  uploadDescription.value = ''
  uploadFile.value = null
  uploadFileName.value = ''
  uploadError.value = null
}

function triggerUpload() {
  resetUploadState()
  isUploadModalOpen.value = true
}

function closeUploadModal() {
  isUploadModalOpen.value = false
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
    return
  }

  const file = files[0]
  const isPdfByType = file.type === 'application/pdf'
  const isPdfByName = file.name.toLowerCase().endsWith('.pdf')

  if (!isPdfByType && !isPdfByName) {
    uploadError.value = 'Only PDF files are allowed.'
    uploadFile.value = null
    uploadFileName.value = ''
    target.value = ''
    return
  }

  uploadFile.value = file
  uploadFileName.value = file.name

  // Do NOT force-fill name; leaving it empty is allowed.
}

/**
 * Computed document name:
 * - If user typed a name, use that.
 * - Else, use the file name (without .pdf) if available.
 * - Else, "Untitled".
 */
const computedDocumentName = computed(() => {
  const trimmed = uploadName.value.trim()
  if (trimmed) return trimmed

  if (uploadFileName.value) {
    return uploadFileName.value.replace(/\.pdf$/i, '')
  }

  return 'Untitled'
})

function addFileToCurrentFolder(item: Item) {
  const last = path.value[path.value.length - 1]
  if (!last) {
    tree.value.push(item)
    return
  }

  if (!last.children) {
    // @ts-ignore - depending on Item definition this may be optional
    last.children = []
  }
  last.children.push(item)
}

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

  const name = computedDocumentName.value

  const newItem: Item = {
    id: `uploaded-${Date.now()}`,
    name: `${name}.pdf`,
    type: 'file',
    // @ts-ignore - size may or may not exist on Item, depending on your DTO
    size: uploadFile.value.size,
  }

  // TODO: integrate with your upload API and persist metadata (tag/description)
  console.log('Uploading PDF with metadata:', {
    name,
    tag: uploadTag.value,
    description: uploadDescription.value,
    file: uploadFile.value,
  })

  addFileToCurrentFolder(newItem)
  closeUploadModal()
}
</script>

<template>
  <div
      class="min-h-screen bg-neutral-50 text-neutral-900 dark:bg-neutral-950 dark:text-neutral-50 transition-colors"
  >
    <div class="min-h-screen flex flex-col">
      <header class="bg-neutral-50/90 dark:bg-neutral-950/90 backdrop-blur-sm">
        <div
            class="mx-auto max-w-6xl px-4 lg:px-6 py-3.5 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-1.5 sm:gap-2">
            <button
                type="button"
                class="inline-flex h-9 items-center gap-1.5 rounded-full border border-neutral-300 bg-white px-3 text-xs sm:text-sm text-neutral-800 hover:border-neutral-400 hover:bg-neutral-50 transition-colors disabled:opacity-40 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800 dark:hover:border-neutral-500"
                :disabled="path.length === 0"
                @click="goHome"
            >
              <HomeIcon class="w-4 h-4" aria-hidden="true" />
              <span class="hidden sm:inline">Home</span>
            </button>

            <button
                type="button"
                class="inline-flex h-9 items-center gap-1.5 rounded-full border border-neutral-300 bg-white px-3 text-xs sm:text-sm text-neutral-800 hover:border-neutral-400 hover:bg-neutral-50 transition-colors disabled:opacity-40 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800 dark:hover:border-neutral-500"
                :disabled="historyIndex === 0"
                @click="goBack"
            >
              <ArrowLeft class="w-4 h-4" aria-hidden="true" />
              <span class="hidden sm:inline">Back</span>
            </button>

            <button
                type="button"
                class="inline-flex h-9 items-center gap-1.5 rounded-full border border-neutral-300 bg-white px-3 text-xs sm:text-sm text-neutral-800 hover:border-neutral-400 hover:bg-neutral-50 transition-colors disabled:opacity-40 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800 dark:hover:border-neutral-500"
                :disabled="historyIndex >= history.length - 1"
                @click="goForward"
            >
              <ArrowRight class="w-4 h-4" aria-hidden="true" />
              <span class="hidden sm:inline">Forward</span>
            </button>
          </div>

          <div class="flex items-center justify-end">
            <div
                class="inline-flex h-9 items-center rounded-full border border-neutral-300 bg-white px-3 sm:px-4 text-[11px] sm:text-xs text-neutral-800 shadow-sm dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100"
            >
              <div
                  class="flex h-7 w-7 items-center justify-center rounded-full bg-emerald-700/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300"
              >
                <BarChart3 class="w-4 h-4" aria-hidden="true" />
              </div>

              <div class="ml-2 flex items-center gap-2 sm:gap-3">
                <span
                    class="hidden sm:inline text-[10px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400"
                >
                  Stats
                </span>

                <span class="whitespace-nowrap">
                  <span class="font-medium">Library:</span>
                  {{ libraryStats.folders }} folders ·
                  {{ libraryStats.files }} files ·
                  {{ formatBytes(libraryStats.bytes) }}
                </span>

                <span class="hidden sm:inline text-neutral-400 dark:text-neutral-600">
                  |
                </span>

                <span class="hidden sm:inline whitespace-nowrap">
                  <span class="font-medium">Level:</span>
                  {{ currentLevelStats.folders }} folders ·
                  {{ currentLevelStats.files }} files ·
                  {{ formatBytes(currentLevelStats.bytes) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </header>

      <main class="flex-1">
        <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
          <nav
              class="text-xs sm:text-sm text-neutral-500 dark:text-neutral-400"
              aria-label="Folder path"
          >
            <ol class="flex flex-wrap items-center gap-1.5">
              <li>
                <button
                    type="button"
                    class="inline-flex items-center gap-1 rounded-full px-2 py-1 hover:bg-neutral-200/80 hover:text-neutral-900 transition-colors dark:hover:bg-neutral-800 dark:hover:text-neutral-50"
                    @click="breadcrumbClick(-1)"
                >
                  <span class="font-medium">Home</span>
                </button>
              </li>

              <template v-for="(node, idx) in path" :key="node.id">
                <li class="text-neutral-400 dark:text-neutral-600">
                  <ChevronRight class="w-3.5 h-3.5" aria-hidden="true" />
                </li>
                <li>
                  <button
                      type="button"
                      class="inline-flex items-center gap-1 rounded-full px-2 py-1 hover:bg-neutral-200/80 hover:text-neutral-900 transition-colors dark:hover:bg-neutral-800 dark:hover:text-neutral-50"
                      @click="breadcrumbClick(idx)"
                  >
                    <Folder
                        v-if="node.type === 'folder'"
                        class="w-3.5 h-3.5 text-neutral-500 dark:text-neutral-400"
                        aria-hidden="true"
                    />
                    <FileText
                        v-else
                        class="w-3.5 h-3.5 text-neutral-500 dark:text-neutral-400"
                        aria-hidden="true"
                    />
                    <span class="truncate max-w-[140px] sm:max-w-[200px]">
                      {{ node.name }}
                    </span>
                  </button>
                </li>
              </template>
            </ol>
          </nav>

          <section
              class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
          >
            <div
                class="border-b border-neutral-200 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 px-4 sm:px-6 py-3.5 flex items-center gap-3 dark:border-neutral-800 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
            >
              <div
                  class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-700/10 border border-emerald-700/40 dark:bg-emerald-500/10 dark:border-emerald-500/50"
              >
                <Folder
                    class="w-5 h-5 text-emerald-800 dark:text-emerald-300"
                    aria-hidden="true"
                />
              </div>
              <div>
                <p
                    class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400"
                >
                  Explorer
                </p>
                <p class="text-sm font-medium">
                  {{ currentFolderLabel }}
                </p>
              </div>
            </div>

            <div class="px-4 sm:px-6 py-5 sm:py-6">
              <div
                  class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-5"
              >
                <div
                    v-for="item in currentItems"
                    :key="item.id"
                    class="group relative flex flex-col rounded-xl border border-neutral-200 bg-neutral-50/80 hover:bg-white hover:border-emerald-600/80 transition-all hover:-translate-y-[1px] hover:shadow-md hover:shadow-emerald-900/10 cursor-pointer dark:border-neutral-700 dark:bg-neutral-900/80 dark:hover:bg-neutral-800 dark:hover:border-emerald-500/80"
                    @click="enterItem(item)"
                >
                  <div class="flex items-start gap-3 p-4">
                    <div
                        class="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-neutral-50 group-hover:bg-emerald-800 transition-colors dark:bg-neutral-200 dark:text-neutral-900 dark:group-hover:bg-emerald-500"
                    >
                      <component
                          :is="iconFor(item)"
                          class="w-5 h-5"
                          aria-hidden="true"
                      />
                    </div>
                    <div class="flex-1 overflow-hidden">
                      <p
                          class="text-sm font-medium truncate"
                          :title="item.name"
                      >
                        {{ item.name }}
                      </p>
                      <p
                          class="text-xs text-neutral-500 dark:text-neutral-400 mt-0.5"
                      >
                        {{ item.type === 'folder' ? 'Folder' : 'PDF document' }}
                      </p>
                    </div>
                  </div>
                </div>

                <div
                    v-if="currentItems.length === 0"
                    class="col-span-full flex flex-col items-center justify-center rounded-xl border border-dashed border-neutral-300 bg-neutral-50 py-10 text-center dark:border-neutral-700 dark:bg-neutral-900"
                >
                  <Folder
                      class="w-6 h-6 text-neutral-400 mb-2 dark:text-neutral-500"
                      aria-hidden="true"
                  />
                  <p class="text-sm text-neutral-600 dark:text-neutral-300">
                    This folder is empty.
                  </p>
                  <p
                      class="text-xs text-neutral-500 dark:text-neutral-400 mt-1"
                  >
                    You will be able to add or upload documents here.
                  </p>
                </div>
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>

    <!-- Floating upload button: big plus, bottom-right -->
    <button
        type="button"
        aria-label="Upload PDF"
        class="fixed bottom-6 right-6 z-30 flex h-14 w-14 items-center justify-center rounded-full bg-neutral-900 text-neutral-50 border border-neutral-800 shadow-lg shadow-neutral-900/40 hover:bg-emerald-700 hover:border-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-neutral-50 dark:bg-neutral-100 dark:text-neutral-900 dark:border-neutral-300 dark:shadow-neutral-950/40 dark:hover:bg-emerald-500 dark:hover:border-emerald-400 dark:focus-visible:ring-offset-neutral-950 transition-all"
        @click="triggerUpload"
    >
      <Plus class="w-6 h-6" aria-hidden="true" />
    </button>

    <!-- Upload modal -->
    <div
        v-if="isUploadModalOpen"
        class="fixed inset-0 z-40 flex items-center justify-center px-4 py-6 bg-black/40 backdrop-blur-sm"
    >
      <div
          class="w-full max-w-lg rounded-2xl border border-neutral-200 bg-white shadow-2xl shadow-neutral-900/30 dark:border-neutral-800 dark:bg-neutral-900"
      >
        <div
            class="flex items-center justify-between border-b border-neutral-200 px-4 sm:px-6 py-3.5 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/60 dark:border-neutral-800 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30 rounded-t-2xl"
        >
          <div class="flex items-center gap-2">
            <div
                class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-700/10 border border-emerald-700/40 dark:bg-emerald-500/10 dark:border-emerald-500/50"
            >
              <FileText class="w-4 h-4 text-emerald-800 dark:text-emerald-300" />
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">
                Upload PDF
              </p>
              <p class="text-sm font-medium text-neutral-900 dark:text-neutral-50">
                Add a new document
              </p>
            </div>
          </div>

          <button
              type="button"
              class="inline-flex h-8 w-8 items-center justify-center rounded-full hover:bg-neutral-200/70 text-neutral-500 hover:text-neutral-900 transition-colors dark:hover:bg-neutral-800 dark:text-neutral-400 dark:hover:text-neutral-100"
              @click="closeUploadModal"
          >
            <X class="w-4 h-4" aria-hidden="true" />
          </button>
        </div>

        <div class="px-4 sm:px-6 py-4 space-y-4">
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

          <div class="space-y-1.5">
            <label class="block text-xs font-medium text-neutral-700 dark:text-neutral-200">
              Metadata / notes
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
            <div
                class="flex items-center justify-between rounded-xl border border-dashed border-neutral-300 bg-neutral-50 px-3 py-2.5 text-sm dark:border-neutral-700 dark:bg-neutral-900/70"
            >
              <div class="flex min-w-0 items-center gap-2">
                <div
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-neutral-900 text-neutral-50 dark:bg-neutral-200 dark:text-neutral-900"
                >
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
              <button
                  type="button"
                  class="ml-3 inline-flex items-center rounded-full border border-neutral-300 bg-white px-3 py-1.5 text-xs font-medium text-neutral-800 hover:border-emerald-500 hover:text-emerald-700 hover:bg-emerald-50 transition-colors dark:border-neutral-600 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:border-emerald-500 dark:hover:bg-emerald-950"
                  @click="openFilePicker"
              >
                Choose file
              </button>
            </div>

            <input
                ref="fileInput"
                type="file"
                accept="application/pdf,.pdf"
                class="hidden"
                @change="onFileChange"
            />

            <p
                v-if="uploadError"
                class="text-xs text-red-600 dark:text-red-400"
            >
              {{ uploadError }}
            </p>
          </div>
        </div>

        <div
            class="flex items-center justify-end gap-2 border-t border-neutral-200 px-4 sm:px-6 py-3.5 bg-neutral-50/80 dark:border-neutral-800 dark:bg-neutral-900/80 rounded-b-2xl"
        >
          <button
              type="button"
              class="inline-flex items-center justify-center rounded-full border border-neutral-300 bg-white px-3.5 py-1.5 text-xs font-medium text-neutral-700 hover:bg-neutral-100 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800"
              @click="closeUploadModal"
          >
            Cancel
          </button>
          <button
              type="button"
              class="inline-flex items-center justify-center rounded-full bg-emerald-700 px-4 py-1.5 text-xs font-medium text-white shadow-sm shadow-emerald-900/30 hover:bg-emerald-800 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-1 focus-visible:ring-offset-neutral-50 dark:focus-visible:ring-offset-neutral-900"
              @click="submitUpload"
          >
            Upload PDF
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
