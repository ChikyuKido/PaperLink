<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Folder,
  FileText,
  ChevronRight,
  ArrowLeft,
  ArrowRight,
  Home as HomeIcon,
  BarChart3,
} from 'lucide-vue-next'
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import CardWithoutBorder from '@/components/own/CardWithoutBorder.vue'
import type { Item } from '@/dto/item'
import { apiFetch } from '@/auth/api'

// --- Backend structure DTOs
type ApiFileNode = { id: string; name: string; size: number }
type ApiDirNode = { id: number; name: string; files: ApiFileNode[]; directories: ApiDirNode[] }

function mapDirNodeToItems(node: ApiDirNode): Item[] {
  const dirs: Item[] = (node.directories ?? []).map(d => ({
    id: String(d.id),
    name: d.name,
    type: 'folder',
    children: mapDirNodeToItems(d),
  }))

  const files: Item[] = (node.files ?? []).map(f => ({
    id: f.id,
    name: f.name.endsWith('.pdf') ? f.name : `${f.name}.pdf`,
    type: 'file',
    size: f.size,
  }))

  return [...dirs, ...files]
}

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

const tree = ref<Item[]>([])
const path = ref<Item[]>([])
const history = ref<Item[][]>([[]])
const historyIndex = ref(0)

const loading = ref(false)
const loadError = ref<string | null>(null)

async function loadTree() {
  loading.value = true
  loadError.value = null
  try {
    const res = await apiFetch('/api/v1/structure/tree', { method: 'GET' })
    const json = await res.json().catch(() => null)

    if (!res.ok || !json?.data) {
      loadError.value = json?.message || 'Failed to load documents.'
      tree.value = []
      return
    }

    const root = json.data as ApiDirNode
    tree.value = mapDirNodeToItems(root)
  } catch {
    loadError.value = 'Failed to load documents.'
    tree.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void loadTree()
})

const currentItems = computed(() => {
  const last = path.value[path.value.length - 1]
  if (!last) return tree.value
  return last.children ?? []
})

const currentFolderLabel = computed(() => {
  const last = path.value[path.value.length - 1]
  return last ? last.name : 'All documents'
})

const libraryStats = computed(() => {
  const result = { folders: 0, files: 0, bytes: 0 }

  const traverse = (items: Item[]) => {
    for (const item of items) {
      if (item.type === 'folder') {
        result.folders++
        if (item.children) traverse(item.children)
      } else {
        result.files++
        if (typeof (item as any).size === 'number') {
          result.bytes += (item as any).size
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
      if (typeof (item as any).size === 'number') {
        bytes += (item as any).size
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
  const next = history.value[historyIndex.value]
  if (next) path.value = next
}

function goForward() {
  if (historyIndex.value >= history.value.length - 1) return
  historyIndex.value++
  const next = history.value[historyIndex.value]
  if (next) path.value = next
}

function openFile(item: Item) {
  console.log('Open file:', item.name)
}

function handleItemClick(item: Item) {
  if (item.type === 'folder') {
    enterItem(item)
  } else {
    openFile(item)
  }
}

function iconFor(item: Item) {
  return item.type === 'folder' ? Folder : FileText
}

// Local payload type for optimistic addUploadedFile (used by Home.vue)
type UploadPayload = {
  name: string
  file: File
}

function addUploadedFile(payload: UploadPayload) {
  const newItem: Item = {
    id: `uploaded-${Date.now()}`,
    name: `${payload.name}.pdf`,
    type: 'file',
    size: payload.file.size as any,
  }

  const last = path.value[path.value.length - 1]
  if (!last) {
    tree.value.push(newItem)
  } else {
    if (!last.children) {
      ;(last as any).children = []
    }
    last.children!.push(newItem)
  }
}

function addCreatedDirectory(name: string, id?: string) {
  const newItem: Item = {
    id: id ?? `dir-${Date.now()}`,
    name,
    type: 'folder',
    children: [],
  }

  const last = path.value[path.value.length - 1]
  if (!last) {
    tree.value.push(newItem)
  } else {
    if (!last.children) (last as any).children = []
    last.children!.push(newItem)
  }
}

function addCreatedDocument(name: string, fileUUID?: string, size?: number) {
  const newItem: Item = {
    id: fileUUID ?? `doc-${Date.now()}`,
    name: name.endsWith('.pdf') ? name : `${name}.pdf`,
    type: 'file',
    size: size,
  }

  const last = path.value[path.value.length - 1]
  if (!last) {
    tree.value.push(newItem)
  } else {
    if (!last.children) (last as any).children = []
    last.children!.push(newItem)
  }
}

defineExpose({
  addUploadedFile,
  addCreatedDirectory,
  addCreatedDocument,
  reload: loadTree,
  getCurrentDirectoryId: () => {
    const last = path.value[path.value.length - 1]
    if (!last) return null
    const n = Number(last.id)
    return Number.isFinite(n) ? n : null
  },
  getCurrentFolderPath: () => path.value.map(p => p.name).join('/'),
})
</script>

<template>
  <div class="min-h-screen bg-neutral-50 text-neutral-900 dark:bg-neutral-950 dark:text-neutral-50 transition-colors">
    <div class="min-h-screen flex flex-col">
      <header class="bg-neutral-50/90 dark:bg-neutral-950/90 backdrop-blur-sm">
        <div class="mx-auto max-w-6xl px-4 lg:px-6 py-3.5 flex items-center justify-between gap-4">
          <TooltipProvider>
            <div class="flex items-center gap-1.5 sm:gap-2">
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button
                      variant="outline"
                      size="sm"
                      class="rounded-full px-3 text-xs sm:text-sm"
                      :disabled="path.length === 0"
                      @click="goHome"
                  >
                    <HomeIcon class="w-4 h-4 mr-1" />
                    <span class="hidden sm:inline">Home</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Home</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger as-child>
                  <Button
                      variant="outline"
                      size="sm"
                      class="rounded-full px-3 text-xs sm:text-sm"
                      :disabled="historyIndex === 0"
                      @click="goBack"
                  >
                    <ArrowLeft class="w-4 h-4 mr-1" />
                    <span class="hidden sm:inline">Back</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Back</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger as-child>
                  <Button
                      variant="outline"
                      size="sm"
                      class="rounded-full px-3 text-xs sm:text-sm"
                      :disabled="historyIndex >= history.length - 1"
                      @click="goForward"
                  >
                    <ArrowRight class="w-4 h-4 mr-1" />
                    <span class="hidden sm:inline">Forward</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Forward</TooltipContent>
              </Tooltip>
            </div>
          </TooltipProvider>

          <div class="flex items-center justify-end">
            <div class="inline-flex h-9 items-center rounded-full border border-neutral-300 bg-white px-3 sm:px-4 text-[11px] sm:text-xs text-neutral-800 shadow-sm dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100">
              <div class="flex h-7 w-7 items-center justify-center rounded-full bg-emerald-700/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300">
                <BarChart3 class="w-4 h-4" aria-hidden="true" />
              </div>

              <div class="ml-2 flex items-center gap-2 sm:gap-3">
                <span class="hidden sm:inline text-[10px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">
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
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink as-child>
                  <button
                      type="button"
                      class="inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs sm:text-sm text-neutral-600 hover:bg-neutral-200/80 hover:text-neutral-900 transition-colors dark:text-neutral-300 dark:hover:bg-neutral-800 dark:hover:text-neutral-50"
                      @click="breadcrumbClick(-1)"
                  >
                    <span class="font-medium">Home</span>
                  </button>
                </BreadcrumbLink>
              </BreadcrumbItem>

              <template v-for="(node, idx) in path" :key="node.id">
                <BreadcrumbSeparator>
                  <ChevronRight class="w-3.5 h-3.5 text-neutral-400 dark:text-neutral-600" />
                </BreadcrumbSeparator>
                <BreadcrumbItem>
                  <BreadcrumbLink as-child>
                    <button
                        type="button"
                        class="inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs sm:text-sm text-neutral-600 hover:bg-neutral-200/80 hover:text-neutral-900 transition-colors dark:text-neutral-300 dark:hover:bg-neutral-800 dark:hover:text-neutral-50"
                        @click="breadcrumbClick(idx)"
                    >
                      <Folder
                          v-if="node.type === 'folder'"
                          class="w-3.5 h-3.5 text-neutral-500 dark:text-neutral-400"
                      />
                      <FileText
                          v-else
                          class="w-3.5 h-3.5 text-neutral-500 dark:text-neutral-400"
                      />
                      <span class="truncate max-w-[140px] sm:max-w-[200px]">
                        {{ node.name }}
                      </span>
                    </button>
                  </BreadcrumbLink>
                </BreadcrumbItem>
              </template>
            </BreadcrumbList>
          </Breadcrumb>

          <section class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none">
            <div class="border-b border-neutral-200 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 px-4 sm:px-6 py-3.5 flex items-center gap-3 dark:border-neutral-800 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30">
              <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-700/10 border border-emerald-700/40 dark:bg-emerald-500/10 dark:border-emerald-500/50">
                <Folder class="w-5 h-5 text-emerald-800 dark:text-emerald-300" aria-hidden="true" />
              </div>
              <div>
                <p class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400">
                  Explorer
                </p>
                <p class="text-sm font-medium">
                  {{ currentFolderLabel }}
                </p>
              </div>
            </div>

            <div class="px-4 sm:px-6 py-5 sm:py-6">
              <div v-if="loadError" class="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-200">
                {{ loadError }}
              </div>

              <div v-if="loading" class="py-10 text-center text-sm text-neutral-500 dark:text-neutral-400">
                Loading...
              </div>

              <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-5">
                <CardWithoutBorder
                    v-for="item in currentItems"
                    :key="item.id"
                    class="group relative flex flex-col rounded-xl bg-neutral-50/80 hover:bg-white hover:border-emerald-600/80 transition-all hover:-translate-y-[1px] hover:shadow-md hover:shadow-emerald-900/10 cursor-pointer dark:bg-neutral-900/80 dark:hover:bg-neutral-800"
                    @click="handleItemClick(item)"
                >
                  <div class="flex items-start gap-3 p-4">
                    <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-neutral-50 group-hover:bg-emerald-800 transition-colors dark:bg-neutral-200 dark:text-neutral-900 dark:group-hover:bg-emerald-500">
                      <component :is="iconFor(item)" class="w-5 h-5" aria-hidden="true" />
                    </div>
                    <div class="flex-1 overflow-hidden">
                      <p class="text-sm font-medium truncate" :title="item.name">
                        {{ item.name }}
                      </p>
                      <p class="text-xs text-neutral-500 dark:text-neutral-400 mt-0.5">
                        {{ item.type === 'folder' ? 'Folder' : 'PDF document' }}
                      </p>
                    </div>
                  </div>
                </CardWithoutBorder>

                <div
                    v-if="currentItems.length === 0"
                    class="col-span-full flex flex-col items-center justify-center rounded-xl border border-dashed border-neutral-300 bg-neutral-50 py-10 text-center dark:border-neutral-700 dark:bg-neutral-900"
                >
                  <Folder class="w-6 h-6 text-neutral-400 mb-2 dark:text-neutral-500" aria-hidden="true" />
                  <p class="text-sm text-neutral-600 dark:text-neutral-300">
                    This folder is empty.
                  </p>
                  <p class="text-xs text-neutral-500 dark:text-neutral-400 mt-1">
                    You will be able to add or upload documents here.
                  </p>
                </div>
              </div>

            </div>
          </section>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
</style>
