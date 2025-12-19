<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-4 lg:py-6 space-y-4">

    <!-- TOP BAR: Compact (Option B) -->
    <section
        class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/60
             overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
          class="px-5 py-2.5 flex items-center justify-between"
          :class="topBarGradient"
      >

        <!-- LEFT SIDE: View mode + Zoom -->
        <div class="flex items-center gap-6">

          <!-- View selector -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-neutral-800 dark:text-neutral-50">
              View
            </span>

            <Select v-model="viewMode">
              <SelectTrigger
                  class="h-8 w-32 rounded-full border-none bg-neutral-100 text-sm
                       text-neutral-900 hover:bg-neutral-50 dark:bg-neutral-800
                       dark:text-neutral-100 dark:hover:bg-neutral-700"
              >
                <SelectValue/>
              </SelectTrigger>

              <SelectContent class="bg-white dark:bg-neutral-900">
                <SelectItem value="scroll">Scroll</SelectItem>
                <SelectItem value="single">Single Page</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- Zoom -->
          <div class="flex items-center gap-3">
            <span class="text-sm font-medium text-neutral-800 dark:text-neutral-50">Zoom</span>

            <div class="flex items-center gap-3 w-48">
              <Slider
                  v-model="scale"
                  :min="0.5"
                  :max="2"
                  :step="0.1"
                  class="h-3 w-full"
              />
              <span class="text-xs text-neutral-700 dark:text-neutral-200 w-10 text-right">
                {{ Math.round(scale * 100) }}%
              </span>
            </div>
          </div>
        </div>

        <!-- CENTER: Navigation (only in single mode) -->
        <div class="flex items-center justify-center flex-1">
          <div v-if="viewMode === 'single'" class="flex items-center gap-2">
            <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full" @click="goFirst">
              <ChevronsLeft class="h-5 w-5"/>
            </Button>

            <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full" @click="prevPage">
              <ChevronLeft class="h-5 w-5"/>
            </Button>

            <div
                class="px-4 py-1.5 rounded-full bg-neutral-50 border border-neutral-200
                     text-sm font-medium text-neutral-900 dark:bg-neutral-800
                     dark:border-neutral-700 dark:text-neutral-100"
            >
              {{ currentPage }} / {{ pageCount }}
            </div>

            <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full" @click="nextPage">
              <ChevronRight class="h-5 w-5"/>
            </Button>

            <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full" @click="goLast">
              <ChevronsRight class="h-5 w-5"/>
            </Button>
          </div>
        </div>
      </div>
    </section>


    <!-- MAIN AREA -->
    <section class="grid gap-4 lg:grid-cols-[240px,minmax(0,1fr)]">



      <!-- RIGHT: UNIFIED PDF VIEWER -->
      <div class="space-y-3">
        <Card
            class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/60
                 dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none p-4"
        >
          <div ref="viewer" class="w-full h-[80vh] overflow-auto">

            <!-- SCROLL MODE -->
            <div
                v-if="viewMode === 'scroll'"
                class="flex flex-col items-center gap-6 py-4"
            >
              <div
                  v-for="idx in pageCount"
                  :key="idx"
                  :ref="el => setPageRef(el, idx)"
                  :data-page="idx"
                  class="w-full flex flex-col items-center"
                  :style="{ transform: `scale(${scale})`, transformOrigin: 'top center' }"
              >
                <Card
                    class="w-[900px] max-w-full rounded-xl overflow-hidden
                         border border-neutral-100 dark:border-neutral-800"
                >
                  <img v-if="pages[idx]" :src="pages[idx]" />
                  <div v-else class="w-full h-[1120px]
                                    flex items-center justify-center text-neutral-500">
                    Loading page {{ idx }}...
                  </div>
                </Card>

                <div class="text-xs mt-2 text-neutral-600 dark:text-neutral-400">
                  Page {{ idx }}
                </div>
              </div>
            </div>

            <!-- SINGLE MODE -->
            <div v-if="viewMode === 'single'" class="flex flex-col items-center py-4">
              <Card
                  class="w-[900px] max-w-full rounded-xl overflow-hidden
                       border border-neutral-100 dark:border-neutral-800"
                  :style="{ transform: `scale(${scale})`, transformOrigin: 'top center' }"
              >
                <img v-if="pages[currentPage]" :src="pages[currentPage]" />
                <div v-else class="w-full h-[1120px] flex items-center justify-center text-neutral-500">
                  Loading page {{ currentPage }}...
                </div>
              </Card>
            </div>

          </div>
        </Card>
      </div>

    </section>
  </div>
</template>



<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, onBeforeUnmount, watch } from "vue"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Slider } from "@/components/ui/slider"
import { ScrollArea } from "@/components/ui/scroll-area"

import {
  FileText,
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight
} from "lucide-vue-next"

import * as pdfjsLib from "pdfjs-dist/build/pdf.mjs"
import pdfWorker from "pdfjs-dist/build/pdf.worker.mjs?worker"
import { useRoute} from 'vue-router';
import {apiFetch} from "@/auth/api.ts";
pdfjsLib.GlobalWorkerOptions.workerPort = new pdfWorker()

const pageCount = 280


const route = useRoute()
const pdfID = route.params.id

const pages = reactive<Record<number, string>>({})
const thumbs = reactive<Record<number, string>>({})
const fetchLocks = reactive<Record<number, boolean>>({})
const pageRefs = reactive<Record<number, HTMLElement | null>>({})

const currentPage = ref(1)
const scale = ref(1)
const viewMode = ref("scroll")
const gotoPageInput = ref(1)

const viewer = ref<HTMLElement | null>(null)

const topBarGradient =
    "bg-gradient-to-r from-neutral-50 via-white to-emerald-50/30 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/10"

let io: IntersectionObserver | null = null

function setPageRef(el: HTMLElement | null, idx: number) {
  pageRefs[idx] = el
}

async function fetchPage(n: number) {
  if (pages[n] || fetchLocks[n]) return
  fetchLocks[n] = true

  try {
    const res = await apiFetch(`/api/v1/pdf/${pdfID}/${n}`)
    if (!res.ok) return
    const buf = await res.arrayBuffer()
    const fixed = buf.slice(1)

    const pdf = await pdfjsLib.getDocument({ data: fixed }).promise
    const page = await pdf.getPage(1)
    const viewport = page.getViewport({ scale: 1 })

    const canvas = document.createElement("canvas")
    const ctx = canvas.getContext("2d")!
    canvas.width = viewport.width
    canvas.height = viewport.height

    await page.render({ canvasContext: ctx, viewport }).promise
    pages[n] = canvas.toDataURL()

    // thumbnail
    canvas.toBlob(blob => {
      if (!blob) return
      const img = new Image()
      const url = URL.createObjectURL(blob)
      img.onload = () => {
        const c = document.createElement("canvas")
        const maxW = 150
        const s = maxW / img.width
        c.width = img.width * s
        c.height = img.height * s
        c.getContext("2d")!.drawImage(img, 0, 0, c.width, c.height)
        thumbs[n] = c.toDataURL()
        URL.revokeObjectURL(url)
      }
      img.src = url
    })
  } finally {
    fetchLocks[n] = false
  }
}

function ensureSurrounding(n: number) {
  ;[n - 1, n, n + 1, n + 2].forEach(p => {
    if (p >= 1 && p <= pageCount) fetchPage(p)
  })
}

function scrollToPage(n: number) {
  n = Math.min(pageCount, Math.max(1, n))

  if (viewMode.value === "single") {
    currentPage.value = n
    gotoPageInput.value = n
    ensureSurrounding(n)
    return
  }

  const el = pageRefs[n]
  if (el) el.scrollIntoView({ behavior: "smooth", block: "center" })

  currentPage.value = n
  gotoPageInput.value = n
  ensureSurrounding(n)
}

function prevPage() { scrollToPage(currentPage.value - 1) }
function nextPage() { scrollToPage(currentPage.value + 1) }
function goFirst() { scrollToPage(1) }
function goLast() { scrollToPage(pageCount) }

function gotoPage() {
  scrollToPage(Number(gotoPageInput.value) || currentPage.value)
}

function onScroll() {
  if (!viewer.value || viewMode.value !== "scroll") return

  const rect = viewer.value.getBoundingClientRect()
  let best = 0
  let bestOverlap = 0

  for (let i = 1; i <= pageCount; i++) {
    const el = pageRefs[i]
    if (!el) continue

    const r = el.getBoundingClientRect()
    const overlap = Math.max(0, Math.min(r.bottom, rect.bottom) - Math.max(r.top, rect.top))

    if (overlap > bestOverlap) {
      bestOverlap = overlap
      best = i
    }
  }

  if (best) {
    currentPage.value = best
    gotoPageInput.value = best
    ensureSurrounding(best)
  }
}

onMounted(async () => {
  ensureSurrounding(1)
  window.addEventListener("keydown", e => {
    if (e.key === "ArrowRight" || e.key === "PageDown") nextPage()
    if (e.key === "ArrowLeft" || e.key === "PageUp") prevPage()
  })

  io = new IntersectionObserver(
      entries => {
        entries.forEach(ent => {
          if (ent.isIntersecting) {
            const p = Number((ent.target as HTMLElement).dataset.page)
            if (p) fetchPage(p)
          }
        })
      },
      { root: viewer.value, rootMargin: "700px" }
  )

  await nextTick()
  for (let i = 1; i <= pageCount; i++) {
    const el = pageRefs[i]
    if (el) io.observe(el)
  }

  viewer.value?.addEventListener("scroll", onScroll)
})

onBeforeUnmount(() => {
  viewer.value?.removeEventListener("scroll", onScroll)
  io?.disconnect()
})
</script>
