import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = 'http://localhost:8080'

export interface Skill {
  id: number
  name: string
  description: string
  type: 'weapon' | 'magic' | 'utility'
  tags: string[]
  primary_stat: 'STR' | 'DEX' | 'INT' | 'WIS'
  level_req: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  classless: boolean
  effects?: Record<string, unknown>
}

export interface SkillInput {
  id?: number
  name: string
  description: string
  type: 'weapon' | 'magic' | 'utility'
  tags: string
  primary_stat: 'STR' | 'DEX' | 'INT' | 'WIS'
  level_req: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  classless: boolean
  effects?: string
}

function parseSkillForApi(input: SkillInput): Omit<Skill, 'id'> {
  return {
    name: input.name,
    description: input.description,
    type: input.type,
    tags: input.tags.split(',').map(t => t.trim()).filter(Boolean),
    primary_stat: input.primary_stat,
    level_req: input.level_req,
    cooldown: input.cooldown,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
    classless: input.classless,
    effects: input.effects ? JSON.parse(input.effects) : undefined
  }
}

export function useSkills(filters?: { type?: string; tag?: string }) {
  return useQuery({
    queryKey: ['skills', filters],
    queryFn: async (): Promise<Skill[]> => {
      const params = new URLSearchParams()
      if (filters?.type) params.append('type', filters.type)
      if (filters?.tag) params.append('tag', filters.tag)
      
      const url = `${API_BASE}/skills${params.toString() ? '?' + params.toString() : ''}`
      const response = await fetch(url)
      if (!response.ok) throw new Error('Failed to fetch skills')
      const data = await response.json()
      return data.skills ?? []
    }
  })
}

export function useSkill(id: number | null) {
  return useQuery({
    queryKey: ['skill', id],
    queryFn: async (): Promise<Skill | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/skills/${id}`)
      if (!response.ok) throw new Error('Failed to fetch skill')
      return response.json()
    },
    enabled: !!id
  })
}

export function useCreateSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (input: SkillInput): Promise<Skill> => {
      const body = parseSkillForApi(input)
      const response = await fetch(`${API_BASE}/skills`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      })
      if (!response.ok) throw new Error('Failed to create skill')
      return response.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['skills'] })
    }
  })
}

export function useUpdateSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, input }: { id: number; input: SkillInput }): Promise<Skill> => {
      const body = parseSkillForApi(input)
      const response = await fetch(`${API_BASE}/skills/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      })
      if (!response.ok) throw new Error('Failed to update skill')
      return response.json()
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['skills'] })
      queryClient.invalidateQueries({ queryKey: ['skill', id] })
    }
  })
}

export function useDeleteSkill() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/skills/${id}`, {
        method: 'DELETE'
      })
      if (!response.ok) throw new Error('Failed to delete skill')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['skills'] })
    }
  })
}
