import { apiFetch } from "@/auth/api"

export type AdminStats = {
  userCount: number
  documentCount: number
  totalDocSize: number
  totalPages: number
  d4sBookCount: number
  d4sAccountCount: number
}

type AdminStatsEnvelope = {
  code: number
  data: AdminStats
}

export async function getAdminStats(): Promise<AdminStats> {
  const res = await apiFetch("/api/v1/admin/stats")
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
  const json = (await res.json()) as AdminStatsEnvelope
  return json.data
}

async function safeError(res: Response): Promise<string> {
  try {
    const json = (await res.json()) as any
    return json?.error ?? `Request failed (${res.status})`
  } catch {
    return `Request failed (${res.status})`
  }
}

