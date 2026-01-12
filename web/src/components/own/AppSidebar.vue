<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
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
  BookOpen as D4SIcon,
  Plug as IntegrationsIcon,
  BarChart3 as StatsIcon,
  SlidersHorizontal as AdminSettingsIcon,
  ChevronUp,
  LogOut,
  User as UserIcon,
  Ticket,
} from 'lucide-vue-next'
import { SIDEBAR_ROUTE_RULES } from '@/router/route_config.ts'
import { checkIsAdmin } from '@/lib/admin'
import { logout } from '@/auth/logout'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { currentUser } from '@/auth/user'
import { ensureCurrentUser } from '@/auth/ensure_user'

const router = useRouter()

async function onLogout() {
  await logout()
  await router.push('/auth')
}

async function goToSettings() {
  await router.push('/settings')
}

const route = useRoute()
const userOpen = ref(true)
const sidebarRule = computed(
    () => SIDEBAR_ROUTE_RULES[route.name as string]
)

const isAdmin = ref(false)

const forceClosed = computed(
    () => sidebarRule.value?.forceClosed === true
)

const effectiveOpen = computed(() => {
  if (forceClosed.value) return false
  return userOpen.value
})

const navItems = computed(() => [
  { title: 'Home', to: '/', icon: HomeIcon },
  { title: 'Search', to: '/search', icon: SearchIcon },
  { title: 'Digi4School', to: '/d4s', icon: D4SIcon },
])

const adminItems = computed(() => [
  { title: 'Settings', to: '/admin/settings', icon: AdminSettingsIcon },
  { title: 'Integrations', to: '/admin/integrations', icon: IntegrationsIcon },
  { title: 'Statistics', to: '/admin/statistics', icon: StatsIcon },
  { title: 'Invites', to: '/admin/invites', icon: Ticket },
])

function isActive(path: string) {
  return route.path === path
}

const displayUsername = computed(() => currentUser.value?.username || 'User')
const initials = computed(() => {
  const u = displayUsername.value.trim()
  if (!u) return 'U'
  return u.slice(0, 2).toUpperCase()
})

onMounted(async () => {
  isAdmin.value = await checkIsAdmin()
  try {
    await ensureCurrentUser()
  } catch {
  }
})
</script>

<template>
  <SidebarProvider :open="effectiveOpen" :key="String(forceClosed)">

    <Sidebar
        collapsible="icon"
        class="border-r border-neutral-200 bg-neutral-50 text-neutral-900 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-50 [--sidebar-width-icon:56px]"
    >
      <SidebarHeader
          class="px-3 pt-3 pb-2"
      >
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
                  v-for="item in navItems"
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

        <SidebarGroup v-if="isAdmin">
          <SidebarGroupLabel
              class="mt-2 text-[11px] uppercase tracking-[0.16em] text-neutral-500 dark:text-neutral-400 group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden"
          >
            Admin
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem
                  v-for="item in adminItems"
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

      <SidebarFooter class="px-2 pb-2 pt-2">
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton
                  size="lg"
                  class="w-full justify-between rounded-xl border border-neutral-200 bg-white px-2 py-2 hover:bg-neutral-50 dark:border-neutral-800 dark:bg-neutral-900 dark:hover:bg-neutral-800"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <Avatar class="h-8 w-8">
                      <AvatarFallback class="bg-emerald-600 text-white dark:bg-emerald-500 dark:text-neutral-950">
                        {{ initials }}
                      </AvatarFallback>
                    </Avatar>
                    <div class="min-w-0 text-left group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden">
                      <p class="truncate text-sm font-medium">{{ displayUsername }}</p>
                      <p class="truncate text-[11px] text-neutral-500 dark:text-neutral-400">Account</p>
                    </div>
                  </div>
                  <ChevronUp class="h-4 w-4 text-neutral-500 group-has-[[data-collapsible=icon]]/sidebar-wrapper:hidden" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>

              <DropdownMenuContent side="top" align="start" class="w-56">
                <DropdownMenuItem class="cursor-pointer" @select.prevent="goToSettings">
                  <UserIcon class="mr-2 h-4 w-4" />
                  User settings
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  class="cursor-pointer text-red-600 focus:text-red-600 dark:text-red-400 dark:focus:text-red-400"
                  @select.prevent="onLogout"
                >
                  <LogOut class="mr-2 h-4 w-4" />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>

    <SidebarInset class="bg-neutral-50 dark:bg-neutral-950">
      <div class="flex min-h-screen flex-col">
        <header
            v-if="!forceClosed"
            class="flex h-12 shrink-0 items-center gap-2 border-b border-neutral-200 px-4 dark:border-neutral-800 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-10"
        >
          <SidebarTrigger
              @click="userOpen = !userOpen"
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
