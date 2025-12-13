import { setAccessToken } from "./auth"

const REFRESH_ENDPOINT = "/api/v1/auth/refresh"

export async function refreshAccessToken() {
    const res = await fetch(REFRESH_ENDPOINT, {
        method: "POST",
        credentials: "include",
    })

    if (!res.ok) {
        throw new Error("Refresh failed")
    }

    const body = await res.json()
    setAccessToken(body.data.access)
}
