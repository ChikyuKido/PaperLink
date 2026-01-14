import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/Home.vue';
import Search from "@/views/Search.vue";
import Settings from "@/views/Settings.vue";
import PDFReader from "@/views/PDFReader.vue";
import Auth from "@/views/Auth.vue";
import D4S from "@/views/D4S.vue";
import AdminSettings from "@/views/AdminSettings.vue";
import AdminIntegrations from "@/views/AdminIntegrations.vue";
import AdminStatistics from "@/views/AdminStatistics.vue";
import AdminInvites from "@/views/AdminInvites.vue";
import { refreshAccessToken } from "@/auth/refresh";
import { accessToken } from "@/auth/auth";
import { ensureCurrentUser } from "@/auth/ensure_user";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: HomeView,
    },
    {
        path: '/search',
        name: 'Search',
        component: Search,
    },
    {
        path: '/settings',
        name: 'Settings',
        component: Settings,
    },
    {
        path: '/pdf/:id',
        name: 'PDF',
        component: PDFReader,
    },
    {
        path: '/d4s',
        name: 'D4S',
        component: D4S,
    },
    {
        path: '/admin',
        name: 'Admin',
        redirect: '/admin/settings',
    },
    {
        path: '/admin/settings',
        name: 'AdminSettings',
        component: AdminSettings,
    },
    {
        path: '/admin/integrations',
        name: 'AdminIntegrations',
        component: AdminIntegrations,
    },
    {
        path: '/admin/statistics',
        name: 'AdminStatistics',
        component: AdminStatistics,
    },
    {
        path: '/admin/invites',
        name: 'AdminInvites',
        component: AdminInvites,
    },
    {
        path: '/auth',
        name: 'Auth',
        component: Auth,
    },


];
const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
});

const PUBLIC_ROUTES = new Set(["Auth"])

router.beforeEach(async (to) => {
  if (PUBLIC_ROUTES.has(String(to.name))) return true

  if (!accessToken.value) {
    try {
      await refreshAccessToken()
    } catch {
      return { name: "Auth" }
    }
  }

  try {
    await ensureCurrentUser()
  } catch {
  }

  return true
})

export default router;