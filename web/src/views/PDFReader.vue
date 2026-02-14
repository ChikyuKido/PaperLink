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

        <div
          v-if="readerError"
          class="m-4 w-full max-w-2xl rounded-lg border border-red-300 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-300"
        >
          {{ readerError }}
        </div>

        <!-- padded wrapper for canvas -->
        <div v-else class="flex justify-center h-full p-4">
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
import { ref, reactive, onMounted, onBeforeUnmount, shallowRef, markRaw } from "vue"
import { useRoute } from "vue-router"
import * as pdfjsLib from "pdfjs-dist/build/pdf.mjs"
import pdfWorker from "pdfjs-dist/build/pdf.worker.mjs?worker"
import { apiFetch } from "@/auth/api.ts"
import { accessToken } from "@/auth/auth.ts"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-vue-next"

pdfjsLib.GlobalWorkerOptions.workerPort = new pdfWorker()

const route = useRoute()
const pdfID = String(route.params.id ?? "")

const pageCount = ref(0)
const currentPage = ref(1)
const canvasEl = ref<HTMLCanvasElement | null>(null)
const thumbnailScrollEl = ref<HTMLDivElement | null>(null)
const thumbnails = ref<(string | null)[]>([])
const THUMB_BATCH_SIZE = 50
const thumbnailBatchLocks = reactive<Record<string, boolean>>({})
const thumbnailBatchCache = reactive<Record<string, boolean>>({})
const readerError = ref<string | null>(null)

let keydownHandler: ((e: KeyboardEvent) => void) | null = null
let currentRenderTask: { cancel: () => void; promise: Promise<void> } | null = null
const pdfDocument = shallowRef<pdfjsLib.PDFDocumentProxy | null>(null)

async function loadDocument() {
  const res = await apiFetch(`/api/v1/document/get/${pdfID}`)
  if (!res.ok) return
  const doc = await res.json()
  pageCount.value = doc.file.pages || 0
  thumbnails.value = Array.from({ length: pageCount.value }, () => null)
}

async function loadPDFDocument() {
  const headers: Record<string, string> = {}
  if (accessToken.value) {
    headers.Authorization = `Bearer ${accessToken.value}`
  }
  const task = pdfjsLib.getDocument({
    url: `/api/v1/pdf/${pdfID}`,
    httpHeaders: headers,
    withCredentials: true,
    rangeChunkSize: 512 * 1024,
    disableAutoFetch: true,
    disableStream: true,
  })
  pdfDocument.value = markRaw(await task.promise)

  if (pdfDocument.value && pageCount.value !== pdfDocument.value.numPages) {
    pageCount.value = pdfDocument.value.numPages
    thumbnails.value = Array.from({ length: pageCount.value }, () => null)
  }
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
}

function onThumbnailScroll() {
  ensureThumbnailBatchesForViewport()
}

function ensureSurrounding(n: number) {
  const pdf = pdfDocument.value
  if (!pdf) return

  const preloadBefore = n
  const preloadAfter = Math.min(pageCount.value, n + 1)
  for (let i = preloadBefore; i <= preloadAfter; i++) {
    void pdf.getPage(i).catch(() => {
      // Best-effort warm cache.
    })
  }
}


let renderToken = 0
async function renderPage(n: number) {
  const token = ++renderToken
  const pdf = pdfDocument.value
  if (!pdf) return
  const canvas = canvasEl.value!
  const ctx = canvas.getContext("2d")!
  if (token !== renderToken) return

  if (currentRenderTask) {
    try {
      currentRenderTask.cancel()
    } catch {
    }
    currentRenderTask = null
  }

  const page = await pdf.getPage(n)
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
  if (pageCount.value === 0) return
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
  try {
    await loadDocument()
    await loadPDFDocument()
    if (pageCount.value === 0) return
    await fetchThumbnailBatch(0)
    ensureThumbnailBatchesForViewport()
    await renderPage(1)
  } catch (err) {
    console.error("Failed to initialize PDF reader", err)
    readerError.value = "Failed to load this PDF."
    return
  }

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
  if (pdfDocument.value) {
    void pdfDocument.value.destroy()
    pdfDocument.value = null
  }
  thumbnails.value.forEach((url) => {
    if (url) URL.revokeObjectURL(url)
  })
})
</script>
