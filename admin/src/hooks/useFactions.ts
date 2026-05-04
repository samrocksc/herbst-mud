import { useState, useEffect, useCallback } from 'react'

export type FactionCategory = Readonly<{
  id: number
  name: string
  description?: string
}>

export type Faction = Readonly<{
  id: number
  name: string
  description?: string
  category_id?: number
  standing: number
  members?: number[]
  is_universal?: boolean
  created_at?: string
}>

export type FactionInput = Readonly<{
  name: string
  description: string
  category_id?: number | null
  standing: number
  is_universal: boolean
}>

const FACTIONS_API = '/api/factions'
const CATEGORIES_API = '/api/faction-categories'

function getToken(): string {
  return localStorage.getItem('token') ?? ''
}

async function apiFetch<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const res = await fetch(`${window.location.origin}${path}`, {
    ...opts,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getToken()}`,
      ...opts.headers,
    },
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export function useFactions() {
  const [factions, setFactions] = useState<Faction[]>([])
  const [categories, setCategories] = useState<FactionCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchFactions = useCallback(async () => {
    setFactions(await apiFetch<Faction[]>(FACTIONS_API))
  }, [])

  const fetchCategories = useCallback(async () => {
    setCategories(await apiFetch<FactionCategory[]>(CATEGORIES_API))
  }, [])

  const refetch = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      await Promise.all([fetchFactions(), fetchCategories()])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load factions')
    } finally {
      setLoading(false)
    }
  }, [fetchFactions, fetchCategories])

  useEffect(() => { void refetch() }, [refetch])

  const createFaction = useCallback(async (input: FactionInput): Promise<Faction> => {
    const data = await apiFetch<Faction>(FACTIONS_API, {
      method: 'POST',
      body: JSON.stringify(input),
    })
    setFactions(prev => [...prev, data])
    return data
  }, [])

  const updateFaction = useCallback(async (id: number, input: FactionInput): Promise<Faction> => {
    const data = await apiFetch<Faction>(`${FACTIONS_API}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    })
    setFactions(prev => prev.map(f => (f.id === id ? data : f)))
    return data
  }, [])

  const deleteFaction = useCallback(async (id: number): Promise<void> => {
    await apiFetch(`${FACTIONS_API}/${id}`, { method: 'DELETE' })
    setFactions(prev => prev.filter(f => f.id !== id))
  }, [])

  const createCategory = useCallback(async (name: string, description?: string): Promise<FactionCategory> => {
    const data = await apiFetch<FactionCategory>(CATEGORIES_API, {
      method: 'POST',
      body: JSON.stringify({ name, description }),
    })
    setCategories(prev => [...prev, data])
    return data
  }, [])

  const deleteCategory = useCallback(async (id: number): Promise<void> => {
    await apiFetch(`${CATEGORIES_API}/${id}`, { method: 'DELETE' })
    setCategories(prev => prev.filter(c => c.id !== id))
  }, [])

  return {
    factions,
    categories,
    loading,
    error,
    refetch,
    createFaction,
    updateFaction,
    deleteFaction,
    createCategory,
    deleteCategory,
  }
}
