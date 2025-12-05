<script setup lang="ts">
import { ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import type HomeExplorer from '@/components/own/home/HomeExplorer.vue'
import UploadPdfModal from '@/components/own/home/UploadPdfModal.vue'
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
}>

const isUploadModalOpen = ref(false)
const homeExplorerRef = ref<HomeExplorerInstance | null>(null)

function handleUploadClick() {
  isUploadModalOpen.value = true
}

function handleUploadClose() {
  isUploadModalOpen.value = false
}

function handleUploadSubmit(payload: UploadPayload) {
  homeExplorerRef.value?.addUploadedFile(payload)
  console.log('Uploading PDF with metadata (from Home.vue):', payload)
  isUploadModalOpen.value = false
}
</script>

<template>
  <div class="relative min-h-screen">
    <HomeExplorerComponent ref="homeExplorerRef" />
    <FloatingUploadButton @click="handleUploadClick" />
    <UploadPdfModal
        :open="isUploadModalOpen"
        @close="handleUploadClose"
        @submit="handleUploadSubmit"
    />
  </div>
</template>

<style scoped>
</style>
