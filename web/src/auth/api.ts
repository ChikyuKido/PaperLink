import { accessToken, setAccessToken } from "./auth"
import { refreshAccessToken } from "./refresh"
import router from "@/router/router.ts";

let isRefreshing = false
let refreshPromise: Promise<void> | null = null
export async function apiFetch(
    input: RequestInfo,
    init: RequestInit = {}
): Promise<Response> {
    const headers = new Headers(init.headers)
    if (accessToken.value) {
        headers.set("Authorization", `Bearer ${accessToken.value}`)
    }
    const response = await fetch(input, {
        ...init,
        headers,
        credentials: "include",
    })

    if (response.status === 403) {
        try {
            await router.push("/auth")
        } catch {
        }
        return response
    }

    if (response.status !== 401) {
        return response
    }
    if (!isRefreshing) {
        isRefreshing = true
        refreshPromise = refreshAccessToken()
            .finally(() => {
                isRefreshing = false
            })
    }

    try {
        await refreshPromise
    } catch {
        setAccessToken(null)
        await router.push("/auth")
        throw new Error("Session expired")
    }

    const retryHeaders = new Headers(init.headers)
    retryHeaders.set("Authorization", `Bearer ${accessToken.value}`)

    return fetch(input, {
        ...init,
        headers: retryHeaders,
        credentials: "include",
    })
}
