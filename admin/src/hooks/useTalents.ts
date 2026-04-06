import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = 'http://localhost:8080'

export interface Talent {
  id: number
  name: string
  description: string
  requirements: string
  effect_type: string
  effect_value: number
  effect_duration: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
}

export interface TalentInput {
  id?: number
  name: string
  description: string
  requirements: string
  effect_type: string
  effect_value: number
  effect_duration: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
}

function parseTalentForApi(input: TalentInput): Omit<Talent, 'id'> {
  return {
    name: input.name,
    description: input.description,
    requirements: input.requirements,
    effect_type: input.effect_type,
    effect_value: input.effect_value,
    effect_duration: input.effect_duration,
    cooldown: input.cooldown,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
  }
}

export function useTalents(filters?: { effectType?: string }) {
  return useQuery({
    queryKey: ['talents', filters],
    queryFn: async (): Promise<Talent[]> => {
      const params = new URLSearchParams()
      if (filters?.effectType) params.append('effectType', filters.effectType)

      const url = `${API_BASE}/talents${params.toString() ? '?' + params.toString() : ''}`
      const response = await fetch(url)
      if (!response.ok) throw new Error('Failed to fetch talents')
      const data = await response.json()
      return data.talents ?? []
    },
  })
}

export function useTalent(id: number | null) {
  return useQuery({
    queryKey: ['talent', id],
    queryFn: async (): Promise<Talent | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/talents/${id}`)
      if (!response.ok) throw new Error('Failed to fetch talent')
      return response.json()
    },
    enabled: !!id,
  })
}

export function useCreateTalent() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (input: TalentInput): Promise<Talent> => {
      const body = parseTalentForApi(input)
      const response = await fetch(`${API_BASE}/talents`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!response.ok) throw new Error('Failed to create talent')
      return response.json()
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['talents'] }),
  })
}

export function useUpdateTalent() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, input }: { id: number; input: TalentInput }): Promise<Talent> => {
      const body = parseTalentForApi(input)
      const response = await fetch(`${API_BASE}/talents/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!response.ok) throw new Error('Failed to update talent')
      return response.json()
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['talents'] })
      queryClient.invalidateQueries({ queryKey: ['talent', id] })
    },
  })
}

export function useDeleteTalent() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/talents/${id}`, { method: 'DELETE' })
      if (!response.ok) throw new Error('Failed to delete talent')
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['talents'] }),
  })
}
