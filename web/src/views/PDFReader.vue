<template>
  <div class="h-[calc(100vh-2rem)] w-full overflow-hidden">
    <section class="flex h-full min-h-0 flex-row flex-nowrap gap-4">

      <!-- LEFT SIDEBAR -->
      <div class="flex h-full min-h-0 w-60 shrink-0 flex-col gap-6 rounded-2xl border border-neutral-200 bg-white
                  p-4 dark:border-neutral-800 dark:bg-neutral-900">

        <!-- Navigation -->
        <div class="space-y-3">
          <div class="text-xs font-semibold text-neutral-500">Navigation</div>

          <div class="flex items-center justify-between">
            <Button size="icon" variant="outline" @click="goFirst">
              <ChevronsLeft class="h-4 w-4"/>
            </Button>

            <Button size="icon" variant="outline" @click="prevPage">
              <ChevronLeft class="h-4 w-4"/>
            </Button>

            <div class="text-sm font-medium text-center px-2">
              {{ currentPage }} / {{ pageCount }}
            </div>

            <Button size="icon" variant="outline" @click="nextPage">
              <ChevronRight class="h-4 w-4"/>
            </Button>

            <Button size="icon" variant="outline" @click="goLast">
              <ChevronsRight class="h-4 w-4"/>
            </Button>
          </div>

        </div>

        <!-- Thumbnails -->
        <div class="flex min-h-0 flex-1 flex-col space-y-2">
          <div class="text-xs font-semibold text-neutral-500">Pages</div>
          <div
            ref="thumbnailScrollEl"
            class="min-h-0 flex-1 overflow-y-auto pr-1 space-y-2"
            @scroll="onThumbnailScroll"
          >
            <button
              v-for="page in pageCount"
              :key="page"
              class="w-full rounded-lg border p-1 text-left transition-colors"
              :class="currentPage === page
                ? 'border-neutral-900 bg-neutral-100 dark:border-neutral-100 dark:bg-neutral-800'
                : 'border-neutral-200 hover:bg-neutral-50 dark:border-neutral-700 dark:hover:bg-neutral-800/70'"
              @click="go(page)"
            >
              <img
                v-if="thumbnails[page - 1]"
                :src="thumbnails[page - 1]"
                :alt="`Page ${page}`"
                class="w-full rounded object-contain"
                loading="lazy"
              >
              <div
                v-else
                class="h-24 w-full rounded bg-neutral-100 dark:bg-neutral-800"
              />
              <div class="mt-1 text-[11px] text-neutral-500">Page {{ page }}</div>
            </button>
          </div>
        </div>

      </div>

      <!-- PDF VIEWER -->
      <div class="flex h-full min-h-0 flex-1 justify-center rounded-2xl border border-neutral-200 bg-white
            dark:border-neutral-800 dark:bg-neutral-900
            overflow-auto">

        <!-- padded wrapper for canvas -->
        <div class="flex justify-center h-full p-4">
          <canvas ref="canvasEl" class="block"></canvas>
        </div>

      </div>


      <!-- RIGHT SIDEBAR PLACEHOLDER -->
      <div class="h-full w-60 shrink-0 rounded-2xl border border-neutral-200 bg-white
                  dark:border-neutral-800 dark:bg-neutral-900 p-4">
        <!-- empty for now -->
      </div>

    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from "vue"
import { useRoute } from "vue-router"
import * as pdfjsLib from "pdfjs-dist/build/pdf.mjs"
import pdfWorker from "pdfjs-dist/build/pdf.worker.mjs?worker"
import { apiFetch } from "@/auth/api.ts"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-vue-next"

pdfjsLib.GlobalWorkerOptions.workerPort = new pdfWorker()

const route = useRoute()
const pdfID = route.params.id

const pageCount = ref(0)
const currentPage = ref(1)
const canvasEl = ref<HTMLCanvasElement | null>(null)
const thumbnailScrollEl = ref<HTMLDivElement | null>(null)
const thumbnails = ref<(string | null)[]>([])
const THUMB_BATCH_SIZE = 50
const thumbnailBatchLocks = reactive<Record<string, boolean>>({})
const thumbnailBatchCache = reactive<Record<string, boolean>>({})

// Cache and fetch locks
const pageCache = reactive<Record<number, Uint8Array>>({})
const fetchLocks = reactive<Record<string, boolean>>({}) // key = "start-end"
let keydownHandler: ((e: KeyboardEvent) => void) | null = null
let currentRenderTask: { cancel: () => void; promise: Promise<void> } | null = null

async function loadDocument() {
  const res = await apiFetch(`/api/v1/document/get/${pdfID}`)
  if (!res.ok) return
  const doc = await res.json()
  pageCount.value = doc.file.pages || 0
  thumbnails.value = Array.from({ length: pageCount.value }, () => null)
}

async function fetchThumbnailBatch(startIndex: number) {
  if (pageCount.value === 0 || startIndex >= pageCount.value) return
  const start = Math.max(0, startIndex)
  const end = Math.min(pageCount.value - 1, start + THUMB_BATCH_SIZE - 1)
  const key = `${start}-${end}`
  if (thumbnailBatchCache[key] || thumbnailBatchLocks[key]) return
  thumbnailBatchLocks[key] = true

  try {
    const res = await apiFetch(`/api/v1/pdf/thumbnails/${pdfID}/${start}-${end}`)
    if (!res.ok) return

    const buf = await res.arrayBuffer()
    const bytes = new Uint8Array(buf)
    const dv = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength)

    let offset = 0
    let pageIndex = start

    while (offset + 8 <= bytes.length && pageIndex <= end) {
      const size = Number(dv.getBigUint64(offset, true))
      offset += 8
      if (size <= 0 || offset + size > bytes.length) break

      const pngBytes = bytes.slice(offset, offset + size)
      offset += size

      const url = URL.createObjectURL(new Blob([pngBytes], { type: "image/png" }))
      if (thumbnails.value[pageIndex]) {
        URL.revokeObjectURL(thumbnails.value[pageIndex]!)
      }
      thumbnails.value[pageIndex] = url
      pageIndex++
    }

    thumbnailBatchCache[key] = true
  } finally {
    thumbnailBatchLocks[key] = false
  }
}

function ensureThumbnailBatchForPage(page: number) {
  const idx = Math.max(0, page - 1)
  const batchStart = Math.floor(idx / THUMB_BATCH_SIZE) * THUMB_BATCH_SIZE
  fetchThumbnailBatch(batchStart)
}

function ensureThumbnailBatchesForViewport() {
  const el = thumbnailScrollEl.value
  if (!el) return

  // Approximate per-item height: preview + label + padding.
  const itemHeight = 128
  const firstVisiblePage = Math.max(1, Math.floor(el.scrollTop / itemHeight) + 1)
  const lastVisiblePage = Math.min(
    pageCount.value,
    Math.ceil((el.scrollTop + el.clientHeight) / itemHeight),
  )

  ensureThumbnailBatchForPage(firstVisiblePage)
  ensureThumbnailBatchForPage(lastVisiblePage)
  ensureThumbnailBatchForPage(lastVisiblePage + THUMB_BATCH_SIZE)
}

function onThumbnailScroll() {
  ensureThumbnailBatchesForViewport()
}

// Fetch multiple pages at once
async function fetchPages(start: number, end: number) {
  start = Math.max(1, start)
  end = Math.min(pageCount.value, end)
  const key = `${start}-${end}`
  if (fetchLocks[key]) return
  fetchLocks[key] = true

  try {
    const res = await apiFetch(`/api/v1/pdf/${pdfID}/${start}-${end}`)
    if (!res.ok) return
    const buf = await res.arrayBuffer()
    const bytes = new Uint8Array(buf)

    if (bytes[0] === 0) {
      // single page
      pageCache[start] = bytes.slice(1)
    } else {
      // multi-page
      let offset = 1
      let pageNum = start
      while (offset < bytes.length && pageNum <= end) {
        const size = Number(new DataView(bytes.buffer, offset, 8).getBigUint64(0))
        offset += 8
        pageCache[pageNum] = bytes.slice(offset, offset + size)
        offset += size
        pageNum++
      }
    }
  } finally {
    fetchLocks[key] = false
  }
}

// Preload surrounding pages
// Preload surrounding pages without fetching already cached pages
function ensureSurrounding(n: number) {
  const preloadBefore = Math.max(1, n - 1)
  const preloadAfter = Math.min(pageCount.value, n + 2)

  // Compute contiguous ranges of pages that are **not cached**
  const ranges: Array<[number, number]> = []
  let rangeStart: number | null = null

  for (let i = preloadBefore; i <= preloadAfter; i++) {
    if (!pageCache[i]) {
      if (rangeStart === null) rangeStart = i
    } else {
      if (rangeStart !== null) {
        ranges.push([rangeStart, i - 1])
        rangeStart = null
      }
    }
  }
  if (rangeStart !== null) ranges.push([rangeStart, preloadAfter])

  // Fetch only missing ranges
  for (const [start, end] of ranges) {
    fetchPages(start, end)
  }
}


// Render a single page from cache
let renderToken = 0
async function renderPage(n: number) {
  const token = ++renderToken
  const canvas = canvasEl.value!
  const ctx = canvas.getContext("2d")!
  if (!pageCache[n]) {
    await fetchPages(n, n)
  }
  if (token !== renderToken || !pageCache[n]) return

  if (currentRenderTask) {
    try {
      currentRenderTask.cancel()
    } catch {
    }
    currentRenderTask = null
  }

  // Use a copy so repeated renders from cache stay reliable.
  const pageBytes = pageCache[n].slice()
  const pdf = await pdfjsLib.getDocument({ data: pageBytes }).promise
  if (token !== renderToken) return
  const page = await pdf.getPage(1)
  if (token !== renderToken) return
  const viewport = page.getViewport({ scale: 1.5 }) // default sharp scale

  canvas.width = viewport.width
  canvas.height = viewport.height

  const renderTask = page.render({ canvasContext: ctx, viewport })
  currentRenderTask = renderTask as { cancel: () => void; promise: Promise<void> }
  await renderTask.promise
  if (token !== renderToken) return

  // Preload surrounding pages
  ensureSurrounding(n)
}

// Navigation
function go(n: number) {
  n = Math.min(pageCount.value, Math.max(1, n))
  currentPage.value = n
  ensureThumbnailBatchForPage(n)
  renderPage(n)
}

const goFirst = () => go(1)
const goLast = () => go(pageCount.value)
const prevPage = () => go(Math.max(currentPage.value - 1, 1))
const nextPage = () => go(Math.min(currentPage.value + 1, pageCount.value))

onMounted(async () => {
  await loadDocument()
  await fetchThumbnailBatch(0)
  ensureThumbnailBatchesForViewport()
  await renderPage(1)

  keydownHandler = (e: KeyboardEvent) => {
    if (e.key === "ArrowRight" || e.key === "PageDown") nextPage()
    if (e.key === "ArrowLeft" || e.key === "PageUp") prevPage()
  }
  window.addEventListener("keydown", keydownHandler)
})

onBeforeUnmount(() => {
  renderToken++
  if (currentRenderTask) {
    try {
      currentRenderTask.cancel()
    } catch {
    }
    currentRenderTask = null
  }
  if (keydownHandler) {
    window.removeEventListener("keydown", keydownHandler)
    keydownHandler = null
  }
  thumbnails.value.forEach((url) => {
    if (url) URL.revokeObjectURL(url)
  })
})
</script>
