<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppSidebar from '@/components/own/AppSidebar.vue'
import { SIDEBAR_ROUTE_RULES } from '@/router/route_config.ts'

const route = useRoute()

const sidebarRule = computed(
    () => SIDEBAR_ROUTE_RULES[route.name as string]
)

const hideSidebar = computed(
    () => sidebarRule.value?.hide === true
)
</script>

<template>
  <div class="flex min-h-screen">
    <AppSidebar v-if="!hideSidebar">
      <main class="p-4 flex-1">
        <RouterView />
      </main>
    </AppSidebar>

    <main v-else class="flex-1">
      <RouterView />
    </main>
  </div>
</template>
