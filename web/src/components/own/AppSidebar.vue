<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import {
  Home as HomeIcon,
  Search as SearchIcon,
  Settings as SettingsIcon,
  Shield as AdminIcon,
} from 'lucide-vue-next'

const route = useRoute()

const isAdmin = computed(() => true)

const menuItems = computed(() => {
  const base = [
    { title: 'Home', to: '/', icon: HomeIcon },
    { title: 'Search', to: '/search', icon: SearchIcon },
    { title: 'Settings', to: '/settings', icon: SettingsIcon },
  ]
  if (isAdmin.value) {
    base.push({ title: 'Admin', to: '/admin', icon: AdminIcon })
  }
  return base
})

function isActive(path: string) {
  return route.path === path
}
</script>

<template>
  <SidebarProvider default-open>
    <Sidebar
        collapsible="icon"
        class="border-r border-neutral-200 bg-neutral-50 text-neutral-900 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-50 [--sidebar-width-icon:56px]"
    >
      <SidebarHeader class="px-3 pt-3 pb-2">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
                size="lg"
                as-child
                class="hover:bg-neutral-100/80 dark:hover:bg-neutral-900/80 group-has-[[data-collapsible=icon]]/sidebar-wrapper:justify-center group-has-[[data-collapsible=icon]]/sidebar-wrapper:px-0"
            >
              <RouterLink
                  to="/"
                  class="flex items-center gap-3 group-has-[[data-collapsible=icon]]/sidebar-wrapper:justify-center group-has-[[data-collapsible=icon]]/sidebar-wrapper:gap-0"
              >
                <div
                    class="flex h-9 w-8 items-center justify-center rounded-md bg-neutral-900 text-neutral-50 text-[10px] font-semibold tracking-[0.2em] dark:bg-neutral-100 dark:text-neutral-900 overflow-hidden"
                >
                  <img
                      src="/logo.png"
                      alt="Paperlink logo"
                      class="h-7 w-7 object-contain"
                  />
                </div>
                <div
                    class="grid text-left text-sm leading-tight group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden"
                >
                  <span class="truncate font-semibold">Paperlink</span>
                  <span class="truncate text-[11px] text-neutral-500 dark:text-neutral-400">
                    Library
                  </span>
                </div>
              </RouterLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent class="px-1">
        <SidebarGroup>
          <SidebarGroupLabel
              class="text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400 group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden"
          >
            Navigation
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem
                  v-for="item in menuItems"
                  :key="item.title"
              >
                <SidebarMenuButton
                    as-child
                    :class="[
                    'flex items-center gap-2 rounded-lg px-2 py-1.5 text-sm font-medium transition-colors',
                    'group-has-[[data-collapsible=icon]]/sidebar-wrapper:px-0',
                    isActive(item.to)
                      ? 'bg-emerald-600 text-white hover:bg-emerald-600/90'
                      : 'text-neutral-800 hover:bg-neutral-100 dark:text-neutral-100 dark:hover:bg-neutral-900',
                  ]"
                >
                  <RouterLink
                      :to="item.to"
                      class="flex items-center gap-2 group-has-[[data-collapsible=icon]]/sidebar-wrapper:justify-center group-has-[[data-collapsible=icon]]/sidebar-wrapper:gap-0"
                  >
                    <component
                        :is="item.icon"
                        class="h-4 w-4 shrink-0"
                    />
                    <span
                        class="truncate group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden"
                    >
                      {{ item.title }}
                    </span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter class="px-3 pb-3 pt-2">
        <div
            class="flex items-center justify-between rounded-lg border border-neutral-200 bg-neutral-50 px-3 py-2 text-[11px] text-neutral-600 dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-400 group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden"
        >
          <span class="truncate">Self-hosted</span>
          <span class="text-[10px] uppercase tracking-[0.16em] text-neutral-400 dark:text-neutral-500">
            Paperlink
          </span>
        </div>
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>

    <SidebarInset class="bg-neutral-50 dark:bg-neutral-950">
      <div class="flex min-h-screen flex-col">
        <header
            class="flex h-12 shrink-0 items-center gap-2 border-b border-neutral-200 px-4 dark:border-neutral-800 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-10"
        >
          <SidebarTrigger
              class="-ml-1 rounded-full border border-neutral-300 bg-white px-2 py-1 text-neutral-800 hover:border-neutral-400 hover:bg-neutral-50 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800 dark:hover:border-neutral-500"
          />
        </header>
        <main class="flex-1">
          <slot />
        </main>
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>

<style scoped>
</style>
