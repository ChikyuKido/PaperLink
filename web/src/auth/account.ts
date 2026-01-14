import { setCurrentUser } from "@/auth/user"

/**
 * Frontend-only stubs for account settings.
 *
 * When the backend is ready, replace the bodies with apiFetch calls, e.g.
 *   PATCH /api/v1/auth/username  { username }
 *   PATCH /api/v1/auth/password  { oldPassword, newPassword }
 */

export async function changeUsername(username: string): Promise<void> {
  // Simulate latency
  await new Promise((r) => setTimeout(r, 350))

  // TODO (backend): validate uniqueness, enforce rules, return updated username
  setCurrentUser({ username })
}

export async function changePassword(_oldPassword: string, _newPassword: string): Promise<void> {
  // Simulate latency
  await new Promise((r) => setTimeout(r, 350))

  // TODO (backend): verify old password, update hash, return ok
}
