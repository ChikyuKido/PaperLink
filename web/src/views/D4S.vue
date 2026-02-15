<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
import { MoreVertical, RefreshCcw, LibraryBig } from "lucide-vue-next"

import { apiFetch } from "@/auth/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"

import { listD4SBooks, takeD4SBook, type Digi4SchoolBook } from "@/lib/d4s_api"

type Notice = { type: "success" | "error"; message: string } | null

// Shared notice (simple in-page toast replace)
const notice = ref<Notice>(null)
let noticeTimer: number | undefined
function showNotice(n: Notice) {
  notice.value = n
  if (noticeTimer) window.clearTimeout(noticeTimer)
  if (n) noticeTimer = window.setTimeout(() => (notice.value = null), 3500)
}

// Library state
const isLoadingBooks = ref(false)
const books = ref<Digi4SchoolBook[]>([])
const bookSearch = ref("")
const takingBookIds = ref(new Set<number>())
const bookThumbnails = ref<Record<number, string>>({})

const filteredBooks = computed(() => {
  const q = bookSearch.value.trim().toLowerCase()
  if (!q) return books.value
  return books.value.filter((b) => (b.bookName ?? "").toLowerCase().includes(q))
})

function revokeThumbnailUrls() {
  for (const url of Object.values(bookThumbnails.value)) {
    URL.revokeObjectURL(url)
  }
}

async function fetchFirstThumbnail(bookID: number): Promise<string | null> {
  const res = await apiFetch(`/api/v1/d4s/thumbnail/${bookID}`)
  if (!res.ok) return null

  const blob = await res.blob()
  return URL.createObjectURL(blob)
}

async function loadBookThumbnails(nextBooks: Digi4SchoolBook[]) {
  const nextThumbnails: Record<number, string> = {}
  const workers = Math.min(6, nextBooks.length)
  let cursor = 0

  async function worker() {
    while (cursor < nextBooks.length) {
      const i = cursor++
      const book = nextBooks[i]
      const url = await fetchFirstThumbnail(book.id).catch(() => null)
      if (url) nextThumbnails[book.id] = url
    }
  }

  await Promise.all(Array.from({ length: workers }, () => worker()))
  revokeThumbnailUrls()
  bookThumbnails.value = nextThumbnails
}

async function loadBooks() {
  isLoadingBooks.value = true
  try {
    const nextBooks = await listD4SBooks()
    books.value = nextBooks
    await loadBookThumbnails(nextBooks)
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to load books" })
  } finally {
    isLoadingBooks.value = false
  }
}

async function onTakeBook(id: number) {
  takingBookIds.value.add(id)
  try {
    await takeD4SBook(id)
    showNotice({ type: "success", message: "Book added to your library." })
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to take book" })
  } finally {
    takingBookIds.value.delete(id)
  }
}

onMounted(async () => {
  await loadBooks()
})

onBeforeUnmount(() => {
  revokeThumbnailUrls()
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
    <!-- Header -->
    <section
      class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
        class="px-4 sm:px-6 py-4 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
      >
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-3">
            <div
              class="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-emerald-600/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200"
            >
              <LibraryBig class="h-5 w-5" />
            </div>
            <div>
              <h1 class="text-lg font-semibold tracking-tight">Digi4School</h1>
              <p class="text-xs text-neutral-500 dark:text-neutral-400">
                Your synced Digi4School library.
              </p>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              class="rounded-full"
              :disabled="isLoadingBooks"
              @click="loadBooks"
            >
              <RefreshCcw class="h-4 w-4" />
              Refresh
            </Button>
          </div>
        </div>

        <div
          v-if="notice"
          class="mt-4 rounded-xl border px-4 py-3 text-sm"
          :class="
            notice.type === 'success'
              ? 'border-emerald-600/30 bg-emerald-600/10 text-emerald-900 dark:text-emerald-200 dark:bg-emerald-500/10'
              : 'border-red-600/30 bg-red-600/10 text-red-900 dark:text-red-200 dark:bg-red-500/10'
          "
        >
          {{ notice.message }}
        </div>
      </div>
    </section>

    <!-- Library -->
    <div class="space-y-4">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="w-full sm:w-80">
          <Input v-model="bookSearch" placeholder="Search books by title…" class="rounded-full" />
        </div>
        <div class="text-xs text-neutral-500 dark:text-neutral-400">
          <span v-if="!isLoadingBooks">
            {{ filteredBooks.length }} book<span v-if="filteredBooks.length !== 1">s</span>
          </span>
          <span v-else>Loading…</span>
        </div>
      </div>

      <div v-if="isLoadingBooks" class="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <Card v-for="i in 8" :key="i" class="overflow-hidden">
          <div class="aspect-[3/4] bg-neutral-200/60 dark:bg-neutral-800/60" />
          <CardContent class="p-3">
            <div class="h-4 w-3/4 bg-neutral-200/70 rounded dark:bg-neutral-800/70" />
            <div class="mt-2 h-3 w-1/2 bg-neutral-200/60 rounded dark:bg-neutral-800/60" />
          </CardContent>
        </Card>
      </div>

      <div v-else class="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <Card
          v-for="book in filteredBooks"
          :key="book.id"
          class="overflow-hidden border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm"
        >
          <div class="relative aspect-[3/4]">
            <img
              v-if="bookThumbnails[book.id]"
              :src="bookThumbnails[book.id]"
              :alt="book.bookName"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div
              v-else
              class="h-full w-full bg-gradient-to-br from-neutral-900 via-neutral-800 to-emerald-900/60 text-white"
            >
              <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_30%_30%,#34d399,transparent_55%)]" />
              <div class="absolute bottom-0 left-0 right-0 p-3">
                <div class="text-[11px] uppercase tracking-[0.16em] text-emerald-200/90">Digi4School</div>
              </div>
            </div>

            <div class="absolute top-2 right-2">
              <Button
                variant="secondary"
                size="icon-sm"
                class="rounded-full bg-black/30 text-white hover:bg-black/40"
                :disabled="takingBookIds.has(book.id)"
                @click="onTakeBook(book.id)"
                :aria-label="`Take ${book.bookName}`"
              >
                <MoreVertical class="h-4 w-4" />
              </Button>
            </div>
          </div>

          <CardContent class="p-3">
            <div class="space-y-1">
              <p class="text-sm font-semibold leading-snug line-clamp-2">
                {{ book.bookName }}
              </p>
              <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
                BookID: {{ book.bookID }}
              </p>
            </div>

            <Separator class="my-3" />

            <Button
              class="w-full rounded-xl bg-emerald-700 text-white hover:bg-emerald-700/90"
              :disabled="takingBookIds.has(book.id)"
              @click="onTakeBook(book.id)"
            >
              <span v-if="!takingBookIds.has(book.id)">Take Book</span>
              <span v-else>Working…</span>
            </Button>
          </CardContent>
        </Card>

        <div
          v-if="!filteredBooks.length"
          class="col-span-full rounded-2xl border border-dashed border-neutral-300 bg-neutral-50 p-8 text-center text-sm text-neutral-600 dark:border-neutral-700 dark:bg-neutral-900/40 dark:text-neutral-300"
        >
          No books found. Ask an admin to run a sync from the Admin page.
        </div>
      </div>
    </div>
  </div>
</template>
