import { createApp } from 'vue';
import App from './App.vue';
import '@/style.css';
import {refreshAccessToken} from "@/auth/refresh.ts";
import router from "@/router/router.ts";
import { ensureCurrentUser } from "@/auth/ensure_user";

const app = createApp(App)
app.use(router)
app.mount("#app")

// Auth bootstrap runs after mount so the UI never stays blank due to bootstrap errors.
;(async function bootstrap() {
  try {
    await refreshAccessToken()
    await ensureCurrentUser()
  } catch {
    await router.push("/auth")
  }
})()
