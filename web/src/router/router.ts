import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/Home.vue';
import Search from "@/views/Search.vue";
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
import TaskView from "@/views/TaskView.vue";
import TasksList from "@/views/TasksList.vue";
import UserSettings from "@/views/UserSettings.vue";

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
        component: UserSettings,
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
    {
        path: '/admin/tasks',
        name: 'Task List',
        component: TasksList,
    },
    {
        path: '/admin/task/:id',
        name: 'Task View',
        component: TaskView,
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