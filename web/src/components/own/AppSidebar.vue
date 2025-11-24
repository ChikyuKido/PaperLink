<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Sidebar,
  SidebarProvider,
  SidebarHeader,
  SidebarContent,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from '@/components/ui/sidebar'
import {
  Home as HomeIcon,
  Search as SearchIcon,
  Settings as SettingsIcon,
  Shield as AdminIcon,
} from 'lucide-vue-next'

const isAdmin = computed(() => {
  return true
})

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
</script>

<template>
  <SidebarProvider default-open>
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" as-child>
              <RouterLink to="/">
                <div class="flex aspect-square h-8 w-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <HomeIcon class="h-4 w-4" />
                </div>
                <div class="ml-2 text-sm font-semibold">MyApp</div>
              </RouterLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in menuItems" :key="item.title">
                <SidebarMenuButton as-child>
                  <RouterLink :to="item.to" class="flex items-center space-x-2">
                    <item.icon class="h-4 w-4" />
                    <span>{{ item.title }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
    <div class="flex-1">
      <slot />
    </div>
  </SidebarProvider>
</template>

<style scoped>
</style>
