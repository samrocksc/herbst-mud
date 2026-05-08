import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type Ability = Readonly<{
  id: number
  name: string
  description: string
  ability_type: string
  cost: number
  cooldown: number
  cooldown_seconds: number
  requirements: string
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
  ability_class: string
  proc_chance: number
  proc_event: string
  faction_skills: number | null
}>

export type AbilityInput = Readonly<{
  id?: number
  name: string
  description: string
  ability_type: string
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
  ability_class: string
  required_tag: string
}>

function parseForApi(input: AbilityInput): Record<string, unknown> {
  return {
    name: input.name,
    description: input.description,
    ability_type: input.ability_type,
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
    ability_class: input.ability_class,
    required_tag: input.required_tag,
  }
}

export function useAbilities(filters?: { type?: string; abilityClass?: string }) {
  return useQuery({
    queryKey: ['abilities', filters],
    queryFn: async (): Promise<Ability[]> => {
      const params = new URLSearchParams()
      if (filters?.type) params.append('type', filters.type)
      if (filters?.abilityClass) params.append('ability_class', filters.abilityClass)
      const url = `${API}/api/abilities${params.toString() ? '?' + params.toString() : ''}`
      return apiGet<Ability[]>(url)
    },
  })
}

export function useAbility(id: number | null) {
  return useQuery({
    queryKey: ['ability', id],
    queryFn: async (): Promise<Ability | null> => {
      if (!id) return null
      return apiGet<Ability>(`${API}/api/abilities/${id}`)
    },
    enabled: !!id,
  })
}

export function useCreateAbility() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: AbilityInput) =>
      apiPost<Ability>(`${API}/api/abilities`, parseForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['abilities'] }),
  })
}

export function useUpdateAbility() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: AbilityInput }) =>
      apiPut<Ability>(`${API}/api/abilities/${id}`, parseForApi(input)),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['abilities'] })
      qc.invalidateQueries({ queryKey: ['ability', id] })
    },
  })
}

export function useDeleteAbility() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/abilities/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['abilities'] }),
  })
}
