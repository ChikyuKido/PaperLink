import { apiFetch } from "@/auth/api"

let cachedIsAdmin: boolean | null = null

export async function checkIsAdmin(): Promise<boolean> {
  // Cache across app lifetime (simple + enough until we have user state).
  if (cachedIsAdmin !== null) return cachedIsAdmin

  try {
    const res = await apiFetch("/api/v1/auth/hasAdmin")
    cachedIsAdmin = res.ok
    return cachedIsAdmin
  } catch {
    cachedIsAdmin = false
    return false
  }
}

export function clearAdminCache() {
  cachedIsAdmin = null
}

