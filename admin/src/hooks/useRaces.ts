/** Race data hook refactored to use the centralized apiFetch. */

import { useState, useEffect, useCallback } from 'react'
import { apiFetch } from '../api'

export type Race = Readonly<{
  id: number
  name: string
  display_name: string
  description: string
  stat_modifiers: Record<string, unknown> | null
  skill_grants: string[]
  ability_modifiers: string[]
  is_playable: boolean
  color: string
}>

export type RaceInput = Readonly<{
  name: string
  display_name: string
  description: string
  stat_modifiers: string
  is_playable: boolean
  color: string
}>

const API = '/api/races'

function parseRaceForApi(input: RaceInput) {
  const body: Record<string, unknown> = {
    name: input.name,
    display_name: input.display_name || input.name,
    description: input.description,
    is_playable: input.is_playable,
    color: input.color,
  }
  if (input.stat_modifiers.trim()) {
    body.stat_modifiers = input.stat_modifiers
  }
  return body
}

export function useRaces() {
  const [races, setRaces] = useState<Race[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchRaces = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await apiFetch(API) as { races: Race[] }
      setRaces(data.races ?? [])
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load races')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { void fetchRaces() }, [fetchRaces])

  const createRace = useCallback(async (input: RaceInput): Promise<Race> => {
    const body = parseRaceForApi(input)
    const data = (await apiFetch(API, {
      method: 'POST',
      body: JSON.stringify(body),
    })) as Race
    setRaces(prev => [...prev, data])
    return data
  }, [])

  const updateRace = useCallback(async (id: number, input: RaceInput): Promise<Race> => {
    const body = parseRaceForApi(input)
    const data = (await apiFetch(`${API}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(body),
    })) as Race
    setRaces(prev => prev.map(r => (r.id === id ? data : r)))
    return data
  }, [])

  const deleteRace = useCallback(async (id: number): Promise<void> => {
    await apiFetch(`${API}/${id}`, { method: 'DELETE' })
    setRaces(prev => prev.filter(r => r.id !== id))
  }, [])

  return { races, loading, error, createRace, updateRace, deleteRace, refetch: fetchRaces }
}
