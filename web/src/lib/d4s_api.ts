import { apiFetch } from "@/auth/api"

export type Digi4SchoolBook = {
  id: number
  uuid: string
  bookName: string
  bookID: string
  accountID: number
}

export type ListBooksResponse = {
  code: number
  data: {
    books: Digi4SchoolBook[]
  }
}

export async function listD4SBooks(): Promise<Digi4SchoolBook[]> {
  const res = await apiFetch("/api/v1/d4s/list")
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
  const json = (await res.json()) as ListBooksResponse
  return json.data?.books ?? []
}

export async function takeD4SBook(id: number): Promise<void> {
  const res = await apiFetch(`/api/v1/d4s/takeBook/${id}`, {
    method: "POST",
  })
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
}

export type D4SAccount = {
  id: number
  username?: string
}

export type CreateAccountResponse = {
  code: number
  data: { id: number }
}

export async function createD4SAccount(username: string, password: string): Promise<number> {
  const res = await apiFetch("/api/v1/d4s/account/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
  const json = (await res.json()) as CreateAccountResponse
  return json.data.id
}

export async function deleteD4SAccount(id: number): Promise<void> {
  const res = await apiFetch(`/api/v1/d4s/account/delete/${id}`, { method: "DELETE" })
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
}

export type SyncResponse = {
  code: number
  data: { id: string }
}

export async function syncD4SAccounts(ids: "all" | number[]): Promise<string> {
  const idsParam = ids === "all" ? "all" : ids.join(",")
  const res = await apiFetch(`/api/v1/d4s/account/sync?ids=${encodeURIComponent(idsParam)}`)
  if (!res.ok) {
    const msg = await safeError(res)
    throw new Error(msg)
  }
  const json = (await res.json()) as SyncResponse
  return json.data.id
}

async function safeError(res: Response): Promise<string> {
  try {
    const json = (await res.json()) as any
    return json?.error ?? `Request failed (${res.status})`
  } catch {
    return `Request failed (${res.status})`
  }
}

