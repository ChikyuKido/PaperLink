import { apiFetch } from "@/auth/api"

// These types mirror the backend /api/v1/structure/tree response.
export type StructureFileNode = {
  id: string
  name: string
  size: number
}

export type StructureDirNode = {
  id: number
  name: string
  files: StructureFileNode[]
  directories: StructureDirNode[]
}

export type SearchIndexItem = {
  id: string
  title: string
  path: string
  sizeBytes: number
}

export async function fetchSearchIndexFromTree(): Promise<SearchIndexItem[]> {
  const res = await apiFetch("/api/v1/structure/tree")
  const body = await res.json().catch(() => null)
  if (!res.ok) {
    throw new Error(body?.error ?? "Failed to load documents")
  }

  const root = body?.data as StructureDirNode
  if (!root) return []

  const out: SearchIndexItem[] = []

  function walk(dir: StructureDirNode, prefix: string) {
    const dirName = dir.name ?? ""
    const currentPrefix = dirName ? (prefix ? `${prefix}/${dirName}` : dirName) : prefix

    for (const f of dir.files ?? []) {
      out.push({
        id: f.id,
        title: f.name,
        path: currentPrefix,
        sizeBytes: f.size ?? 0,
      })
    }

    for (const child of dir.directories ?? []) {
      walk(child, currentPrefix)
    }
  }

  walk(root, "")
  return out
}

