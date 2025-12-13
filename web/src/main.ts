import { createApp } from 'vue';
import App from './App.vue';
import '@/style.css';
import {refreshAccessToken} from "@/auth/refresh.ts";
import router from "@/router/router.ts";
async function bootstrap() {
    try {
        await refreshAccessToken()
    } catch {
        await router.push("/auth")
    }

    const app = createApp(App)
    app.use(router)
    app.mount("#app")
}

await bootstrap()
