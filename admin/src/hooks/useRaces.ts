import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}/api/races`

export type Race = Readonly<{
  id: number
  name: string
  display_name: string
  description: string
  stat_modifiers: Record<string, unknown> | null
  skill_grants: string[]
  ability_modifiers: string[]
  equipment_slots: string[]
  is_playable: boolean
  color: string
}>

export type RaceInput = Readonly<{
  name: string
  display_name: string
  description: string
  stat_modifiers: string
  equipment_slots: string[]
  is_playable: boolean
  color: string
}>

function parseRaceForApi(input: RaceInput) {
  const body: Record<string, unknown> = {
    name: input.name,
    display_name: input.display_name || input.name,
    description: input.description,
    equipment_slots: input.equipment_slots,
    is_playable: input.is_playable,
    color: input.color,
  }
  if (input.stat_modifiers.trim()) {
    body.stat_modifiers = input.stat_modifiers
  }
  return body
}

export function useRaces() {
  return useQuery({
    queryKey: ['races'],
    queryFn: async (): Promise<Race[]> => {
      const data = await apiGet<Race[]>(API)
      return Array.isArray(data) ? data : []
    },
  })
}

export function useCreateRace() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: RaceInput) =>
      apiPost<Race>(API, parseRaceForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['races'] }),
  })
}

export function useUpdateRace() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: RaceInput }) =>
      apiPut<Race>(`${API}/${id}`, parseRaceForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['races'] }),
  })
}

export function useDeleteRace() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['races'] }),
  })
}