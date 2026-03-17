<template>
  <div class="h-[calc(100vh-2rem)] w-full overflow-hidden">
    <section class="flex h-full min-h-0 gap-3">
      <aside class="flex h-full min-h-0 w-56 shrink-0 flex-col rounded-2xl border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900">
        <div class="border-b border-neutral-200 px-4 py-3 dark:border-neutral-800">
          <div class="text-[11px] font-semibold uppercase tracking-[0.18em] text-neutral-500">Pages</div>
          <div class="mt-1 text-sm font-medium text-neutral-900 dark:text-neutral-100">
            {{ currentPage }} / {{ pageCount || '—' }}
          </div>
        </div>

        <div
          ref="thumbnailScrollEl"
          class="min-h-0 flex-1 space-y-2 overflow-y-auto p-2"
          @scroll="onThumbnailScroll"
        >
          <button
            v-for="page in pageCount"
            :key="page"
            :data-page="page"
            :aria-current="currentPage === page ? 'page' : undefined"
            class="w-full rounded-xl border p-2 text-left transition-colors"
            :class="currentPage === page
              ? 'border-neutral-900 bg-neutral-100 shadow-sm dark:border-neutral-100 dark:bg-neutral-800'
              : 'border-neutral-200 hover:bg-neutral-50 dark:border-neutral-700 dark:hover:bg-neutral-800/70'"
            @click="go(page)"
          >
            <img
              v-if="thumbnails[page - 1]"
              :src="thumbnails[page - 1]"
              :alt="`Page ${page}`"
              class="w-full rounded-md object-contain"
              loading="lazy"
            >
            <div
              v-else
              class="flex h-24 w-full items-center justify-center rounded-md bg-neutral-100 text-[11px] text-neutral-400 dark:bg-neutral-800"
            >
              Loading preview
            </div>
            <div class="mt-2 flex items-center justify-between gap-2">
              <div class="text-[11px] font-medium text-neutral-600 dark:text-neutral-300">
                Page {{ page }}
              </div>
              <div
                v-if="currentPage === page"
                class="rounded-full bg-neutral-900 px-2 py-0.5 text-[10px] font-medium text-white dark:bg-neutral-100 dark:text-neutral-900"
              >
                Current
              </div>
            </div>
          </button>
        </div>
      </aside>

      <div class="flex min-h-0 flex-1 flex-col gap-3">
        <header class="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-neutral-200 bg-white px-4 py-3 dark:border-neutral-800 dark:bg-neutral-900">
          <div class="flex items-center gap-2">
            <Button size="icon-sm" variant="outline" :disabled="isFirstPage" @click="goFirst">
              <ChevronsLeft class="h-4 w-4" />
            </Button>
            <Button size="icon-sm" variant="outline" :disabled="isFirstPage" @click="prevPage">
              <ChevronLeft class="h-4 w-4" />
            </Button>
            <div class="min-w-[5.5rem] text-center text-sm font-medium text-neutral-900 dark:text-neutral-100">
              {{ currentPage }} / {{ pageCount || '—' }}
            </div>
            <Button size="icon-sm" variant="outline" :disabled="isLastPage" @click="nextPage">
              <ChevronRight class="h-4 w-4" />
            </Button>
            <Button size="icon-sm" variant="outline" :disabled="isLastPage" @click="goLast">
              <ChevronsRight class="h-4 w-4" />
            </Button>
          </div>

          <div class="flex flex-wrap items-center gap-2">
            <Badge
              variant="outline"
              class="gap-1.5 border-neutral-200 bg-neutral-50 text-neutral-600 dark:border-neutral-700 dark:bg-neutral-800 dark:text-neutral-200"
            >
              <LoaderCircle
                v-if="collabStatus === 'connecting'"
                class="h-3.5 w-3.5 animate-spin"
              />
              <Wifi
                v-else-if="collabStatus === 'connected'"
                class="h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400"
              />
              <WifiOff v-else class="h-3.5 w-3.5 text-neutral-400 dark:text-neutral-500" />
              {{ collabLabel }}
            </Badge>
            <div class="flex items-center gap-2">
              <Input
                v-model="pageInput"
                type="number"
                inputmode="numeric"
                min="1"
                :max="Math.max(pageCount, 1)"
                placeholder="Page"
                class="h-8 w-20"
                @keydown.enter="submitPageJump"
              />
              <Button size="sm" variant="secondary" :disabled="pageCount === 0" @click="submitPageJump">
                Go
              </Button>
            </div>
            <div class="hidden rounded-full border border-neutral-200 px-2 py-1 text-[10px] text-neutral-500 dark:border-neutral-700 dark:text-neutral-300 md:block">
              ← → ↑ ↓ navigate
            </div>
          </div>
        </header>

        <div
          v-if="collabError"
          class="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-900/80 dark:bg-amber-950/40 dark:text-amber-200"
        >
          {{ collabError }}
        </div>

        <div class="flex min-h-0 flex-1 overflow-auto rounded-2xl border border-neutral-200 bg-neutral-100 dark:border-neutral-800 dark:bg-neutral-950">
          <div
            v-if="readerError"
            class="m-4 w-full max-w-2xl rounded-lg border border-red-300 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-300"
          >
            {{ readerError }}
          </div>

          <div v-else class="flex min-h-full w-full justify-center p-4 md:p-6">
            <canvas
              ref="canvasEl"
              class="block max-w-full border border-neutral-200 bg-white shadow-sm dark:border-neutral-800 dark:bg-white"
            />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { usePdfReader } from '@/composables/usePdfReader'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, LoaderCircle, Wifi, WifiOff } from 'lucide-vue-next'

const {
  pageCount,
  currentPage,
  canvasEl,
  thumbnailScrollEl,
  thumbnails,
  readerError,
  collabStatus,
  collabError,
  onThumbnailScroll,
  go,
  goFirst,
  goLast,
  prevPage,
  nextPage,
} = usePdfReader()

const pageInput = ref('1')

const isFirstPage = computed(() => currentPage.value <= 1)
const isLastPage = computed(() => pageCount.value === 0 || currentPage.value >= pageCount.value)
const collabLabel = computed(() => {
  if (collabStatus.value === 'connecting') return 'Live sync connecting'
  if (collabStatus.value === 'connected') return 'Live sync connected'
  return 'Live sync offline'
})

watch(
  currentPage,
  (page) => {
    pageInput.value = String(page)
  },
  { immediate: true },
)

function submitPageJump() {
  const parsed = Number.parseInt(pageInput.value, 10)
  if (Number.isNaN(parsed)) return
  go(parsed)
}
</script>
