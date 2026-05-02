import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = `${window.location.origin}`

function authHeaders(): HeadersInit {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export interface Ability {
  id: number
  name: string
  description: string
  skill_type: string
  cost: number
  cooldown: number
  cooldown_seconds: number
  requirements: number
  effect_type: string
  effect_value: number
  effect_duration: number
  mana_cost: number
  stamina_cost: number
  hp_cost: number
  scaling_stat: string
  scaling_percent_per_point: number
  slug: string
  required_tag: string
  skill_class: string
  proc_chance: number
  proc_event: string
  faction_skills: number | null
}

export interface AbilityInput {
  id?: number
  name: string
  description: string
  skill_type: string
  requirements: number
  cost: number
  cooldown: number
  cooldown_seconds: number
  mana_cost: number
  stamina_cost: number
  hp_cost: number
  effect_type: string
  effect_value: number
  effect_duration: number
  scaling_stat: string
  scaling_percent_per_point: number
  proc_chance: number
  proc_event: string
  skill_class: string
  required_tag: string
}

function parseAbilityForApi(input: AbilityInput): Record<string, unknown> {
  return {
    name: input.name,
    description: input.description,
    skill_type: input.skill_type,
    requirements: input.requirements,
    cost: input.cost,
    cooldown: input.cooldown,
    cooldown_seconds: input.cooldown_seconds,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
    hp_cost: input.hp_cost,
    effect_type: input.effect_type,
    effect_value: input.effect_value,
    effect_duration: input.effect_duration,
    scaling_stat: input.scaling_stat,
    scaling_percent_per_point: input.scaling_percent_per_point,
    proc_chance: input.proc_chance,
    proc_event: input.proc_event,
    skill_class: input.skill_class,
    required_tag: input.required_tag,
  }
}

export function useAbilities(filters?: { type?: string }) {
  return useQuery({
    queryKey: ['abilities', filters],
    queryFn: async (): Promise<Ability[]> => {
      const params = new URLSearchParams()
      if (filters?.type) params.append('type', filters.type)

      const url = `${API_BASE}/skills${params.toString() ? '?' + params.toString() : ''}`
      const response = await fetch(url, { headers: authHeaders() })
      if (!response.ok) throw new Error('Failed to fetch abilities')
      const data = await response.json()
      return data.skills ?? []
    }
  })
}

export function useAbility(id: number | null) {
  return useQuery({
    queryKey: ['ability', id],
    queryFn: async (): Promise<Ability | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/skills/${id}`, { headers: authHeaders() })
      if (!response.ok) throw new Error('Failed to fetch ability')
      return response.json()
    },
    enabled: !!id
  })
}

export function useCreateAbility() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (input: AbilityInput): Promise<Ability> => {
      const body = parseAbilityForApi(input)
      const response = await fetch(`${API_BASE}/skills`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', ...authHeaders() },
        body: JSON.stringify(body)
      })
      if (!response.ok) throw new Error('Failed to create ability')
      return response.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['abilities'] })
    }
  })
}

export function useUpdateAbility() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, input }: { id: number; input: AbilityInput }): Promise<Ability> => {
      const body = parseAbilityForApi(input)
      const response = await fetch(`${API_BASE}/skills/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', ...authHeaders() },
        body: JSON.stringify(body)
      })
      if (!response.ok) throw new Error('Failed to update ability')
      return response.json()
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['abilities'] })
      queryClient.invalidateQueries({ queryKey: ['ability', id] })
    }
  })
}

export function useDeleteAbility() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/skills/${id}`, {
        method: 'DELETE',
        headers: authHeaders()
      })
      if (!response.ok) throw new Error('Failed to delete ability')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['abilities'] })
    }
  })
}
