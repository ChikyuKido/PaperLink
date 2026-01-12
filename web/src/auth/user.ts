import { ref } from "vue"
import { apiFetch } from "@/auth/api"

export type CurrentUser = {
  username: string
}

const STORAGE_KEY = "paperlink.username"

function loadStoredUser(): CurrentUser | null {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    return stored ? { username: stored } : null
  } catch {
    return null
  }
}

export const currentUser = ref<CurrentUser | null>(loadStoredUser())

export function setCurrentUser(user: CurrentUser | null) {
  currentUser.value = user
  try {
    if (!user?.username) {
      localStorage.removeItem(STORAGE_KEY)
    } else {
      localStorage.setItem(STORAGE_KEY, user.username)
    }
  } catch {
  }
}

export async function fetchCurrentUser(): Promise<CurrentUser | null> {
  const res = await apiFetch("/api/v1/auth/me")

  // Backend not implemented yet: keep cached/persisted username.
  if (res.status === 404) {
    return currentUser.value
  }

  if (!res.ok) {
    // Don’t wipe local username on transient failures.
    return currentUser.value
  }

  const json = (await res.json()) as { code: number; data: CurrentUser }
  setCurrentUser(json.data)
  return json.data
}
