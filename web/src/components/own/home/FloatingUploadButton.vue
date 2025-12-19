<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Plus, FileText, FolderPlus } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'create-document'): void
  (e: 'create-directory'): void
}>()

const open = ref(false)

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

function onGlobalPointerDown(e: PointerEvent) {
  const target = e.target as HTMLElement | null
  if (!target) return
  if (target.closest('[data-create-menu-root]')) return
  close()
}

function onGlobalKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

onMounted(() => {
  window.addEventListener('pointerdown', onGlobalPointerDown)
  window.addEventListener('keydown', onGlobalKeyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener('pointerdown', onGlobalPointerDown)
  window.removeEventListener('keydown', onGlobalKeyDown)
})

function chooseDocument() {
  close()
  emit('create-document')
}

function chooseDirectory() {
  close()
  emit('create-directory')
}
</script>

<template>
  <div class="fixed bottom-6 right-6 z-30" data-create-menu-root>
    <!-- Backdrop for a clearer focus state + easier outside click target -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <button
        v-if="open"
        type="button"
        aria-label="Close create menu"
        class="fixed inset-0 z-20 cursor-default bg-neutral-950/10 dark:bg-neutral-950/40"
        @click="close"
      />
    </transition>

    <div class="relative z-30 flex flex-col items-end">
      <transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 translate-y-2 scale-95"
        enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 translate-y-0 scale-100"
        leave-to-class="opacity-0 translate-y-2 scale-95"
      >
        <div v-if="open" class="mb-3 w-[260px]">
          <div class="rounded-2xl border border-neutral-200 bg-white/95 p-2 shadow-xl shadow-neutral-900/10 backdrop-blur dark:border-neutral-800 dark:bg-neutral-900/95">
            <p class="px-2.5 pb-2 pt-1 text-[11px] font-medium uppercase tracking-[0.14em] text-neutral-500 dark:text-neutral-400">
              Create
            </p>

            <button
              type="button"
              class="group flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition hover:bg-emerald-50/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 dark:hover:bg-emerald-500/10"
              @click="chooseDocument"
            >
              <span class="flex h-9 w-9 items-center justify-center rounded-xl border border-emerald-600/30 bg-emerald-600/10 text-emerald-800 transition group-hover:bg-emerald-600/15 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-200">
                <FileText class="h-4.5 w-4.5" aria-hidden="true" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block truncate text-sm font-semibold text-neutral-900 dark:text-neutral-50">Document</span>
                <span class="block truncate text-xs text-neutral-500 dark:text-neutral-400">Upload a PDF and add metadata</span>
              </span>
            </button>

            <button
              type="button"
              class="group flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition hover:bg-emerald-50/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 dark:hover:bg-emerald-500/10"
              @click="chooseDirectory"
            >
              <span class="flex h-9 w-9 items-center justify-center rounded-xl border border-emerald-600/30 bg-emerald-600/10 text-emerald-800 transition group-hover:bg-emerald-600/15 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-200">
                <FolderPlus class="h-4.5 w-4.5" aria-hidden="true" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block truncate text-sm font-semibold text-neutral-900 dark:text-neutral-50">Directory</span>
                <span class="block truncate text-xs text-neutral-500 dark:text-neutral-400">Create a new folder here</span>
              </span>
            </button>
          </div>
        </div>
      </transition>

      <!-- Main FAB -->
      <button
        type="button"
        :aria-label="open ? 'Close create menu' : 'Open create menu'"
        :aria-expanded="open"
        class="flex h-14 w-14 items-center justify-center rounded-full bg-neutral-900 text-neutral-50 border border-neutral-800 shadow-lg shadow-neutral-900/40 hover:bg-emerald-700 hover:border-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-neutral-50 dark:bg-neutral-100 dark:text-neutral-900 dark:border-neutral-300 dark:shadow-neutral-950/40 dark:hover:bg-emerald-500 dark:hover:border-emerald-400 dark:focus-visible:ring-offset-neutral-950 transition-all"
        @click="toggle"
      >
        <Plus class="w-6 h-6 transition-transform duration-200" :class="open ? 'rotate-45' : ''" aria-hidden="true" />
      </button>
    </div>
  </div>
</template>

<style scoped>
</style>
