import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = `${window.location.origin}`

export interface WeaponSkill {
  id: number
  name: string
  description: string
  skill_category: string
  effect_type: string
  effect_value: number
  effect_duration: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  requirements: string
}

export interface WeaponSkillInput {
  id?: number
  name: string
  description: string
  skill_category: string
  effect_type: string
  effect_value: number
  effect_duration: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  requirements: string
}

function parseWeaponSkillForApi(input: WeaponSkillInput): Record<string, unknown> {
  return {
    name: input.name,
    description: input.description,
    skill_category: input.skill_category,
    effect_type: input.effect_type,
    effect_value: input.effect_value,
    effect_duration: input.effect_duration,
    cooldown: input.cooldown,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
    requirements: input.requirements,
  }
}

export function useWeaponSkills() {
  return useQuery({
    queryKey: ['weapon-skills'],
    queryFn: async (): Promise<WeaponSkill[]> => {
      const response = await fetch(`${API_BASE}/talents`)
      if (!response.ok) throw new Error('Failed to fetch weapon skills')
      const data = await response.json()
      return data.talents ?? []
    }
  })
}

export function useWeaponSkill(id: number | null) {
  return useQuery({
    queryKey: ['weapon-skill', id],
    queryFn: async (): Promise<WeaponSkill | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/talents/${id}`)
      if (!response.ok) throw new Error('Failed to fetch weapon skill')
      return response.json()
    },
    enabled: !!id,
  })
}

export function useCreateWeaponSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (input: WeaponSkillInput): Promise<WeaponSkill> => {
      const body = parseWeaponSkillForApi(input)
      const response = await fetch(`${API_BASE}/talents`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!response.ok) throw new Error('Failed to create weapon skill')
      return response.json()
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['weapon-skills'] }),
  })
}

export function useUpdateWeaponSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, input }: { id: number; input: WeaponSkillInput }): Promise<WeaponSkill> => {
      const body = parseWeaponSkillForApi(input)
      const response = await fetch(`${API_BASE}/talents/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!response.ok) throw new Error('Failed to update weapon skill')
      return response.json()
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['weapon-skills'] })
      queryClient.invalidateQueries({ queryKey: ['weapon-skill', id] })
    },
  })
}

export function useDeleteWeaponSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/talents/${id}`, { method: 'DELETE' })
      if (!response.ok) throw new Error('Failed to delete weapon skill')
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['weapon-skills'] }),
  })
}
