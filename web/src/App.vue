<script setup>
import { RouterView } from 'vue-router'
import AppSidebar from '@/components/own/AppSidebar.vue'
</script>

<template>
  <div class="flex min-h-screen">
    <AppSidebar>
      <main class="p-4 flex-1">
        <RouterView />
      </main>
    </AppSidebar>
  </div>
</template>

<style scoped>
</style>
