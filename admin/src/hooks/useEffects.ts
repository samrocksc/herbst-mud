import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type AbilityEffect = Readonly<{
  id: number
  effect_type: string
  damage_subtype: string
  target: string
  value: number
  duration: number
  scaling_stat: string
  scaling_ratio: number
  sort_order: number
  ability_id?: number
}>

export type EffectInput = Readonly<{
  effect_type: string
  damage_subtype?: string
  target?: string
  value?: number
  duration?: number
  scaling_stat?: string
  scaling_ratio?: number
  sort_order?: number
}>

export function useEffects(abilityId: number | null) {
  return useQuery({
    queryKey: ['effects', abilityId],
    queryFn: async (): Promise<AbilityEffect[]> => {
      if (!abilityId) return []
      return apiGet<AbilityEffect[]>(`${API}/api/abilities/${abilityId}/effects`)
    },
    enabled: !!abilityId,
  })
}

export function useCreateEffect() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ abilityId, input }: { abilityId: number; input: EffectInput }) =>
      apiPost<AbilityEffect>(`${API}/api/abilities/${abilityId}/effects`, input),
    onSuccess: (_, { abilityId }) => {
      qc.invalidateQueries({ queryKey: ['effects', abilityId] })
      qc.invalidateQueries({ queryKey: ['abilities'] })
    },
  })
}

export function useUpdateEffect() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; abilityId: number; input: Partial<EffectInput> }) =>
      apiPut<AbilityEffect>(`${API}/api/effects/${id}`, input),
    onSuccess: (_, { abilityId }) => {
      qc.invalidateQueries({ queryKey: ['effects', abilityId] })
      qc.invalidateQueries({ queryKey: ['abilities'] })
    },
  })
}

export function useDeleteEffect() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id }: { id: number; abilityId: number }) =>
      apiDelete(`${API}/api/effects/${id}`),
    onSuccess: (_, { abilityId }) => {
      qc.invalidateQueries({ queryKey: ['effects', abilityId] })
      qc.invalidateQueries({ queryKey: ['abilities'] })
    },
  })
}