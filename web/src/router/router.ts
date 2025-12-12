import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/Home.vue';
import Search from "@/views/Search.vue";
import Settings from "@/views/Settings.vue";
import PDFReader from "@/views/PDFReader.vue";
import Auth from "@/views/Auth.vue";

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
        path: '/',
        name: 'Settings',
        component: Settings,
    },
    {
        path: '/pdf',
        name: 'PDF',
        component: PDFReader,
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
export default router;