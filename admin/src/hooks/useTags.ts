import { useState, useEffect, useCallback } from 'react'

export type Tag = Readonly<{
  id: number
  name: string
  color: string
}>

export type TagInput = Readonly<{
  name: string
  color: string
}>

export type TagUsage = Readonly<{
  id: number
  name: string
  type: string
}>

export type TagUsageReport = Readonly<{
  tag_name: string
  total_usages: number
  skills: TagUsage[]
  factions: TagUsage[]
  characters: TagUsage[]
}>

const API = '/api/tags'

function getToken() {
  return localStorage.getItem('token') ?? ''
}

async function apiFetch(path: string, opts?: RequestInit) {
  const res = await fetch(path, {
    ...opts,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`,
      ...opts?.headers,
    },
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

export function useTags() {
  const [tags, setTags] = useState<Tag[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTags = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await apiFetch(API) as { tags: Tag[] }
      setTags(data.tags ?? [])
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load tags')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { void fetchTags() }, [fetchTags])

  const createTag = useCallback(async (input: TagInput): Promise<Tag> => {
    const data = await apiFetch(API, {
      method: 'POST',
      body: JSON.stringify(input),
    }) as Tag
    setTags(prev => [...prev, data])
    return data
  }, [])

  const updateTag = useCallback(async (id: number, input: Partial<TagInput>): Promise<Tag> => {
    const data = await apiFetch(`${API}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    }) as Tag
    setTags(prev => prev.map(t => t.id === id ? data : t))
    return data
  }, [])

  const deleteTag = useCallback(async (id: number): Promise<void> => {
    await apiFetch(`${API}/${id}`, { method: 'DELETE' })
    setTags(prev => prev.filter(t => t.id !== id))
  }, [])

  return { tags, loading, error, createTag, updateTag, deleteTag, refetch: fetchTags }
}

export async function fetchTagUsages(id: number): Promise<TagUsageReport> {
  const res = await fetch(`${API}/${id}/usages`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`,
    },
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json() as Promise<TagUsageReport>
}
