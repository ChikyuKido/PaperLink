<script setup lang="ts">
import { ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import UploadPdfModal from '@/components/own/home/UploadPdfModal.vue'
import CreateDirectoryModal from '@/components/own/home/CreateDirectoryModal.vue'
import FloatingUploadButton from '@/components/own/home/FloatingUploadButton.vue'
import HomeExplorerComponent from '@/components/own/home/HomeExplorer.vue'

type UploadPayload = {
  name: string
  tag: string
  description: string
  file: File
  createdAt: string | null
  modifiedAt: string | null
}

type HomeExplorerInstance = ComponentPublicInstance<{
  addUploadedFile: (payload: UploadPayload) => void
  addCreatedDirectory: (name: string, id?: string) => void
  addCreatedDocument: (name: string, fileUUID?: string, size?: number) => void
  reload: () => Promise<void>
  getCurrentDirectoryId: () => number | null
  getCurrentFolderPath: () => string
}>

const isUploadModalOpen = ref(false)
const isCreateDirectoryModalOpen = ref(false)
const homeExplorerRef = ref<HomeExplorerInstance | null>(null)

function currentDirectoryId() {
  return homeExplorerRef.value?.getCurrentDirectoryId() ?? null
}

function currentFolderPath() {
  return homeExplorerRef.value?.getCurrentFolderPath?.() ?? ''
}

function handleCreateDocument() {
  isUploadModalOpen.value = true
}

function handleCreateDirectory() {
  isCreateDirectoryModalOpen.value = true
}

function handleUploadClose() {
  isUploadModalOpen.value = false
}

function handleDirectoryClose() {
  isCreateDirectoryModalOpen.value = false
}

function handleUploadSubmit(payload: { name: string; fileUUID: string; directoryId: number | null }) {
  // optimistic UI: show it immediately
  homeExplorerRef.value?.addCreatedDocument(payload.name, payload.fileUUID)
  isUploadModalOpen.value = false
}

function handleDirectoryCreated(payload: { id: number; name: string; parentId: number | null }) {
  // optimistic UI: show it immediately with real backend id
  homeExplorerRef.value?.addCreatedDirectory(payload.name, String(payload.id))
  isCreateDirectoryModalOpen.value = false
}
</script>

<template>
  <div class="relative min-h-screen">
    <HomeExplorerComponent ref="homeExplorerRef" />

    <FloatingUploadButton
      @create-document="handleCreateDocument"
      @create-directory="handleCreateDirectory"
    />

    <UploadPdfModal
      :open="isUploadModalOpen"
      :directory-id="currentDirectoryId()"
      :folder-path="currentFolderPath()"
      @close="handleUploadClose"
      @submit="handleUploadSubmit"
    />

    <CreateDirectoryModal
      :open="isCreateDirectoryModalOpen"
      :parent-id="currentDirectoryId()"
      :folder-path="currentFolderPath()"
      @close="handleDirectoryClose"
      @created="handleDirectoryCreated"
    />
  </div>
</template>

<style scoped>
</style>
