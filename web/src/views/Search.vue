<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
    <!-- Header / Search -->
    <section
        class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
          class="px-4 sm:px-6 py-4 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
      >
        <form
            class="w-full"
            @submit.prevent="onSearch"
        >
          <div
              class="flex w-full items-center gap-2 rounded-full border border-neutral-300 bg-white px-3 py-2 shadow-sm transition-colors dark:border-neutral-700 dark:bg-neutral-900"
          >
            <SearchIcon class="h-4 w-4 text-neutral-500 dark:text-neutral-400" />

            <Input
                v-model="searchQuery"
                type="search"
                placeholder="Search documents, text inside PDFs, annotations, tags..."
                class="flex-1 border-0 bg-transparent shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 text-sm"
            />

            <Select v-model="selectedScope">
              <SelectTrigger
                  class="h-9 w-[140px] rounded-full border-none bg-neutral-100 text-[11px] sm:text-xs text-neutral-800 hover:bg-neutral-50 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700"
              >
                <SelectValue placeholder="Scope" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All documents</SelectItem>
                <SelectItem value="mine">My PDFs</SelectItem>
                <SelectItem value="shared">Shared</SelectItem>
              </SelectContent>
            </Select>

            <Button
                type="submit"
                class="h-9 rounded-full px-4 text-xs sm:text-sm"
            >
              Search
            </Button>
          </div>
        </form>
      </div>
    </section>

    <!-- Content -->
    <section class="grid gap-4 lg:grid-cols-[260px,minmax(0,1fr)]">
      <!-- Left: Filters -->
      <Card
          class="h-fit border border-neutral-200 bg-white shadow-sm shadow-neutral-200/60 dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
      >
        <CardHeader class="pb-2">
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span
                  class="inline-flex h-7 w-7 items-center justify-center rounded-full bg-emerald-700/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300"
              >
                <FilterIcon class="h-3.5 w-3.5" />
              </span>
              <div>
                <CardTitle class="text-sm">Filters</CardTitle>
                <CardDescription class="text-[11px]">
                  Narrow down your results
                </CardDescription>
              </div>
            </div>
            <Button
                variant="ghost"
                size="icon"
                class="h-7 w-7 rounded-full text-neutral-500 hover:bg-neutral-100 dark:text-neutral-300 dark:hover:bg-neutral-800"
                @click="resetFilters"
            >
              <FilterIcon class="h-4 w-4" />
            </Button>
          </div>
        </CardHeader>

        <CardContent class="space-y-3 pt-1">
          <!-- Sort -->
          <div class="space-y-1.5">
            <p class="text-[11px] font-medium text-neutral-500 dark:text-neutral-400">
              Sort by
            </p>
            <Tabs v-model="selectedSort" class="w-full">
              <TabsList
                  class="grid w-full grid-cols-3 bg-neutral-100 text-[11px] dark:bg-neutral-800"
              >
                <TabsTrigger
                    value="relevance"
                    class="text-neutral-600 dark:text-neutral-300 data-[state=active]:bg-emerald-600/15 data-[state=active]:text-emerald-800 data-[state=active]:shadow-sm dark:data-[state=active]:bg-emerald-500/20 dark:data-[state=active]:text-emerald-200"
                >
                  Relevance
                </TabsTrigger>
                <TabsTrigger
                    value="recent"
                    class="text-neutral-600 dark:text-neutral-300 data-[state=active]:bg-emerald-600/15 data-[state=active]:text-emerald-800 data-[state=active]:shadow-sm dark:data-[state=active]:bg-emerald-500/20 dark:data-[state=active]:text-emerald-200"
                >
                  Recent
                </TabsTrigger>
                <TabsTrigger
                    value="az"
                    class="text-neutral-600 dark:text-neutral-300 data-[state=active]:bg-emerald-600/15 data-[state=active]:text-emerald-800 data-[state=active]:shadow-sm dark:data-[state=active]:bg-emerald-500/20 dark:data-[state=active]:text-emerald-200"
                >
                  A–Z
                </TabsTrigger>
              </TabsList>
            </Tabs>
          </div>

          <Separator />

          <!-- Tags with single-row collapsed + expandable full list -->
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-2">
              <p
                  class="flex items-center gap-1.5 text-xs font-medium text-neutral-600 dark:text-neutral-300"
              >
                <TagIcon class="h-3.5 w-3.5 text-emerald-700 dark:text-emerald-300" />
                Tags
              </p>
              <Button
                  variant="ghost"
                  size="icon"
                  class="h-6 w-6 rounded-full text-neutral-500 hover:bg-neutral-100 dark:text-neutral-300 dark:hover:bg-neutral-800"
                  @click="showAllTags = !showAllTags"
              >
                <ChevronDown
                    class="h-3.5 w-3.5 transition-transform"
                    :class="showAllTags ? 'rotate-180' : ''"
                />
              </Button>
            </div>

            <div class="flex items-start gap-2">
              <!-- Left: tag list -->
              <div class="flex-1">
                <!-- Collapsed single-row mode (same height as search bar) -->
                <div
                    v-if="!showAllTags"
                    class="h-9 flex items-center rounded-full border border-neutral-200 bg-neutral-50 px-2 overflow-x-auto whitespace-nowrap dark:border-neutral-700 dark:bg-neutral-900/60"
                >
                  <Badge
                      v-for="tag in filteredTags"
                      :key="tag"
                      variant="outline"
                      class="mr-1 shrink-0 cursor-pointer rounded-full border border-dashed px-2 py-0.5 text-[11px] text-neutral-700 dark:text-neutral-200 dark:border-neutral-600"
                      :class="selectedTags.includes(tag)
                      ? 'border-emerald-600 bg-emerald-600/10 text-emerald-800 dark:border-emerald-400 dark:bg-emerald-500/15 dark:text-emerald-200'
                      : ''"
                      @click="toggleTag(tag)"
                  >
                    {{ tag }}
                  </Badge>

                  <span
                      v-if="!filteredTags.length"
                      class="text-[11px] text-neutral-500 dark:text-neutral-400"
                  >
                    No tags match your search.
                  </span>
                </div>

                <!-- Expanded multi-row mode -->
                <ScrollArea
                    v-else
                    class="h-32 rounded-2xl border border-neutral-200 bg-neutral-50 px-2 py-2 dark:border-neutral-700 dark:bg-neutral-900/60"
                >
                  <div class="flex flex-wrap gap-1.5">
                    <Badge
                        v-for="tag in filteredTags"
                        :key="tag"
                        variant="outline"
                        class="cursor-pointer rounded-full border border-dashed px-2 py-0.5 text-[11px] text-neutral-700 dark:text-neutral-200 dark:border-neutral-600"
                        :class="selectedTags.includes(tag)
                        ? 'border-emerald-600 bg-emerald-600/10 text-emerald-800 dark:border-emerald-400 dark:bg-emerald-500/15 dark:text-emerald-200'
                        : ''"
                        @click="toggleTag(tag)"
                    >
                      {{ tag }}
                    </Badge>

                    <p
                        v-if="!filteredTags.length"
                        class="py-1 text-[11px] text-neutral-500 dark:text-neutral-400"
                    >
                      No tags match your search.
                    </p>
                  </div>
                </ScrollArea>
              </div>

              <!-- Right: tag search (same height as tag row) -->
              <div class="w-40 shrink-0">
                <Input
                    v-model="tagSearch"
                    type="search"
                    placeholder="Search tags"
                    class="h-9 text-xs"
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Right: Results -->
      <div class="space-y-2">
        <!-- Meta / status line -->
        <div
            class="flex flex-wrap items-center justify-between gap-2 text-xs text-neutral-500 dark:text-neutral-400"
        >
          <div class="flex items-center gap-2">
            <span v-if="!isLoading">
              Showing
              <span class="font-medium text-neutral-900 dark:text-neutral-50">
                {{ filteredResults.length }}
              </span>
              result<span v-if="filteredResults.length !== 1">s</span>
            </span>
            <span
                v-else
                class="flex items-center gap-1.5"
            >
              <Loader2 class="h-3.5 w-3.5 animate-spin text-emerald-600" />
              Searching your workspace…
            </span>
          </div>
          <div class="flex items-center gap-2 text-[11px]">
            <Clock class="h-3 w-3 text-emerald-600 dark:text-emerald-400" />
            <span>Updated live as you type</span>
          </div>
        </div>

        <!-- Results list -->
        <div class="space-y-1.5">
          <!-- Loading skeleton -->
          <template v-if="isLoading">
            <Card
                v-for="n in 3"
                :key="n"
                class="border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900"
            >
              <CardContent class="flex gap-3 p-3">
                <Skeleton class="mt-1 h-9 w-7 rounded-md" />
                <div class="flex-1 space-y-1.5">
                  <Skeleton class="h-4 w-[55%]" />
                  <Skeleton class="h-3 w-[80%]" />
                  <div class="flex gap-2">
                    <Skeleton class="h-4 w-16" />
                    <Skeleton class="h-4 w-20" />
                  </div>
                </div>
              </CardContent>
            </Card>
          </template>

          <!-- Empty state -->
          <template v-else-if="!filteredResults.length">
            <Card
                class="border border-dashed border-neutral-300 bg-neutral-50 py-4 dark:border-neutral-700 dark:bg-neutral-900"
            >
              <CardContent class="flex items-center gap-3 p-3">
                <div
                    class="flex h-9 w-9 items-center justify-center rounded-full bg-neutral-900 text-neutral-50 dark:bg-neutral-200 dark:text-neutral-900"
                >
                  <SearchIcon class="h-5 w-5" />
                </div>
                <div class="space-y-0.5">
                  <p class="text-sm font-medium text-neutral-800 dark:text-neutral-100">
                    No results found
                  </p>
                  <p class="text-xs text-neutral-500 dark:text-neutral-400">
                    Try a different keyword or adjust your filters.
                  </p>
                </div>
              </CardContent>
            </Card>
          </template>

          <!-- Result cards: use CardWithoutBorder -->
          <template v-else>
            <CardWithoutBorder
                v-for="result in filteredResults"
                :key="result.id"
                class="group border border-neutral-200 bg-white transition hover:-translate-y-[1px] hover:border-emerald-600/80 hover:shadow-md hover:shadow-emerald-900/10 dark:border-neutral-800 dark:bg-neutral-900 dark:hover:border-emerald-500/80"
            >
              <div class="flex gap-3 p-3">
                <!-- Icon column -->
                <div
                    class="mt-0.5 flex h-9 w-7 items-center justify-center rounded-lg bg-neutral-900 text-neutral-50 group-hover:bg-emerald-700 transition-colors dark:bg-neutral-200 dark:text-neutral-900 dark:group-hover:bg-emerald-500"
                >
                  <FileText class="h-4 w-4" />
                </div>

                <!-- Content column -->
                <div class="flex-1 space-y-1">
                  <div class="flex items-start justify-between gap-2">
                    <div>
                      <h2
                          class="line-clamp-1 text-sm font-medium tracking-tight text-neutral-900 dark:text-neutral-50"
                      >
                        {{ result.title }}
                      </h2>
                      <p
                          class="mt-0.5 line-clamp-2 text-xs text-neutral-500 dark:text-neutral-400"
                      >
                        {{ result.description }}
                      </p>
                    </div>
                  </div>

                  <div class="flex flex-wrap items-center gap-1 pt-0.5">
                    <Badge
                        variant="outline"
                        class="border-dashed text-[10px] text-neutral-600 dark:text-neutral-300"
                    >
                      {{ result.pages }} pages
                    </Badge>
                    <Badge
                        variant="outline"
                        class="text-[10px] text-neutral-600 dark:text-neutral-300"
                    >
                      {{ result.size }}
                    </Badge>
                    <Badge
                        variant="outline"
                        class="text-[10px] text-neutral-600 dark:text-neutral-300"
                    >
                      <Clock class="mr-1 h-3 w-3" />
                      {{ result.updatedAt }}
                    </Badge>
                  </div>

                  <div class="flex flex-wrap items-center gap-2 pt-0.5">
                    <div
                        class="flex items-center gap-1.5 text-[11px] text-neutral-500 dark:text-neutral-400"
                    >
                      <Users class="h-3.5 w-3.5" />
                      <span>{{ result.owner }}</span>
                      <span v-if="result.shared">· shared</span>
                    </div>

                    <div class="flex flex-wrap gap-1">
                      <Badge
                          v-for="tag in result.tags"
                          :key="tag"
                          variant="secondary"
                          class="rounded-full px-2 py-0.5 text-[10px]"
                      >
                        {{ tag }}
                      </Badge>
                    </div>
                  </div>
                </div>
              </div>
            </CardWithoutBorder>
          </template>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import CardWithoutBorder from '@/components/own/CardWithoutBorder.vue'
import { fetchSearchIndexFromTree, type SearchIndexItem } from '@/lib/search_api'

import {
  Search as SearchIcon,
  Filter as FilterIcon,
  FileText,
  Clock,
  Users,
  Tag as TagIcon,
  Loader2,
  ChevronDown,
} from 'lucide-vue-next'

type Scope = 'all' | 'mine' | 'shared'
type Sort = 'relevance' | 'recent' | 'az'

interface SearchResult {
  id: string
  title: string
  description: string
  tags: string[]
  pages: number
  size: string
  updatedAt: string
  owner: string
  shared: boolean
  path: string
}

const searchQuery = ref('')
const selectedScope = ref<Scope>('all')
const selectedSort = ref<Sort>('relevance')
const selectedTags = ref<string[]>([])
const tagSearch = ref('')
const isLoading = ref(false)
const showAllTags = ref(false)
const loadError = ref<string | null>(null)

const results = ref<SearchResult[]>([])

function formatBytes(bytes: number) {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let v = bytes
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

function pathToTags(path: string) {
  // Use folder segments as tags (simple + useful with current backend data)
  const normalized = (path || '').trim().replace(/^\/+|\/+$/g, '')
  if (!normalized) return ['Root']
  return normalized.split('/').filter(Boolean)
}

async function loadFromBackend() {
  isLoading.value = true
  loadError.value = null
  try {
    const items: SearchIndexItem[] = await fetchSearchIndexFromTree()
    results.value = items.map((it) => ({
      id: it.id,
      title: it.title,
      description: it.path ? `Folder: ${it.path}` : 'Folder: /',
      tags: pathToTags(it.path),
      pages: 0,
      size: formatBytes(it.sizeBytes),
      updatedAt: '',
      owner: 'You',
      shared: false,
      path: it.path,
    }))
  } catch (e: any) {
    loadError.value = e?.message ?? 'Failed to load documents'
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  await loadFromBackend()
})

const tags = computed(() => {
  const set = new Set<string>()
  for (const r of results.value) {
    for (const t of r.tags) set.add(t)
  }
  return Array.from(set).sort((a, b) => a.localeCompare(b))
})

const filteredTags = computed(() => {
  const q = tagSearch.value.trim().toLowerCase()
  const all = tags.value
  if (!q) return all
  return all.filter((t) => t.toLowerCase().includes(q))
})

const filteredResults = computed(() => {
  let list = [...results.value]

  const q = searchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter(
      (item) =>
        item.title.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q) ||
        item.tags.some((t) => t.toLowerCase().includes(q)),
    )
  }

  // scope placeholders until backend provides sharing/ownership info
  if (selectedScope.value === 'shared') {
    list = list.filter((item) => item.shared)
  }

  if (selectedTags.value.length) {
    list = list.filter((item) => selectedTags.value.every((tag) => item.tags.includes(tag)))
  }

  if (selectedSort.value === 'az') {
    list.sort((a, b) => a.title.localeCompare(b.title))
  }

  return list
})

function onSearch() {
  // This is now an instant local filter; keep the small loading UX.
  isLoading.value = true
  window.setTimeout(() => {
    isLoading.value = false
  }, 200)
}

function toggleTag(tag: string) {
  if (selectedTags.value.includes(tag)) {
    selectedTags.value = selectedTags.value.filter((t) => t !== tag)
  } else {
    selectedTags.value = [...selectedTags.value, tag]
  }
}

function resetFilters() {
  selectedScope.value = 'all'
  selectedSort.value = 'relevance'
  selectedTags.value = []
  tagSearch.value = ''
  showAllTags.value = false
}
</script>
