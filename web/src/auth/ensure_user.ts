import { accessToken } from "@/auth/auth"
import { currentUser, fetchCurrentUser } from "@/auth/user"

let inflight: Promise<void> | null = null
let validatedOnce = false

export async function ensureCurrentUser(): Promise<void> {
  if (!accessToken.value) return
  if (inflight) return inflight

  // Even if we have a cached user, revalidate at least once per app load.
  if (currentUser.value && validatedOnce) return

  inflight = (async () => {
    try {
      await fetchCurrentUser()
      validatedOnce = true
    } finally {
      inflight = null
    }
  })()

  return inflight
}
