import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/Home.vue';
import Search from "@/views/Search.vue";
import Settings from "@/views/Settings.vue";

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
    },    {
        path: '/',
        name: 'Settings',
        component: Settings,
    },

];
const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
});
export default router;