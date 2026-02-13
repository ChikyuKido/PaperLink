<template>
  <div class="w-full h-screen overflow-hidden px-4 py-0">
    <section class="flex flex-row flex-nowrap gap-4 h-full">

      <!-- LEFT SIDEBAR -->
      <div class="w-60 shrink-0 h-full rounded-2xl border border-neutral-200 bg-white
                  dark:border-neutral-800 dark:bg-neutral-900 p-4 space-y-6">

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

      </div>

      <!-- PDF VIEWER -->
      <div class="flex-1 h-full rounded-2xl border border-neutral-200 bg-white
            dark:border-neutral-800 dark:bg-neutral-900
            overflow-auto flex justify-center">

        <!-- padded wrapper for canvas -->
        <div class="flex justify-center h-full p-4">
          <canvas ref="canvasEl" class="block"></canvas>
        </div>

      </div>


      <!-- RIGHT SIDEBAR PLACEHOLDER -->
      <div class="w-60 shrink-0 h-full rounded-2xl border border-neutral-200 bg-white
                  dark:border-neutral-800 dark:bg-neutral-900 p-4">
        <!-- empty for now -->
      </div>

    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
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

// Cache and fetch locks
const pageCache = reactive<Record<number, Uint8Array>>({})
const fetchLocks = reactive<Record<string, boolean>>({}) // key = "start-end"

async function loadDocument() {
  const res = await apiFetch(`/api/v1/document/get/${pdfID}`)
  if (!res.ok) return
  const doc = await res.json()
  pageCount.value = doc.file.pages || 0
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
let rendering = false
async function renderPage(n: number) {
  const canvas = canvasEl.value!
  const ctx = canvas.getContext("2d")!
  if (!pageCache[n]) {
    await fetchPages(n, n)
  }

  const pdf = await pdfjsLib.getDocument({ data: pageCache[n] }).promise
  const page = await pdf.getPage(1)
  const viewport = page.getViewport({ scale: 1.5 }) // default sharp scale

  canvas.width = viewport.width
  canvas.height = viewport.height

  await page.render({ canvasContext: ctx, viewport }).promise

  // Preload surrounding pages
  ensureSurrounding(n)
}

// Navigation
function go(n: number) {
  n = Math.min(pageCount.value, Math.max(1, n))
  currentPage.value = n
  renderPage(n)
}

const goFirst = () => go(1)
const goLast = () => go(pageCount.value)
const prevPage = () => go(Math.max(currentPage.value - 1, 1))
const nextPage = () => go(Math.min(currentPage.value + 1, pageCount.value))

onMounted(async () => {
  await loadDocument()
  await renderPage(1)

  window.addEventListener("keydown", e => {
    if (e.key === "ArrowRight" || e.key === "PageDown") nextPage()
    if (e.key === "ArrowLeft" || e.key === "PageUp") prevPage()
  })
})
</script>
