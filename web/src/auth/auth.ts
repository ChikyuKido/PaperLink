import { ref } from "vue"

export const accessToken = ref<string | null>(null)

export function setAccessToken(token: string | null) {
    accessToken.value = token
}
