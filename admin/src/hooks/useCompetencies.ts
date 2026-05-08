import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type CompetencyThreshold = Readonly<{
  level: number
  xp_required: number
  damage_multiplier: number
  defense_multiplier: number
}>

export type CompetencyCategory = Readonly<{
  id: string
  name: string
  xp_multiplier: number
  thresholds: CompetencyThreshold[]
}>

export type CompetencyThresholdInput = Readonly<{
  level: number
  xp_required: number
  damage_multiplier: number
  defense_multiplier: number
}>

export type CompetencyCategoryInput = Readonly<{
  id: string
  name: string
  xp_multiplier: number
  thresholds: CompetencyThresholdInput[]
}>

export type CompetencyCategoryUpdate = Readonly<{
  name?: string
  xp_multiplier?: number
  thresholds?: CompetencyThresholdInput[]
}>

export type CharacterCompetency = Readonly<{
  category_id: string
  category_name: string
  xp: number
  level: number
  xp_multiplier: number
  damage_multiplier: number
  defense_multiplier: number
}>

export function useCompetencyCategories() {
  return useQuery({
    queryKey: ['competency-categories'],
    queryFn: async (): Promise<CompetencyCategory[]> => {
      const data = await apiGet<CompetencyCategory[]>(`${API}/api/competency-categories`)
      return Array.isArray(data) ? data : []
    },
  })
}

export function useCreateCompetencyCategory() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: CompetencyCategoryInput) =>
      apiPost<CompetencyCategory>(`${API}/api/competency-categories`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['competency-categories'] }),
  })
}

export function useUpdateCompetencyCategory() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: CompetencyCategoryUpdate }) =>
      apiPut<CompetencyCategory>(`${API}/api/competency-categories/${id}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['competency-categories'] }),
  })
}

export function useDeleteCompetencyCategory() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => apiDelete(`${API}/api/competency-categories/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['competency-categories'] }),
  })
}

export function useCharacterCompetencies(characterId: number | null) {
  return useQuery({
    queryKey: ['character-competencies', characterId],
    queryFn: async (): Promise<CharacterCompetency[]> => {
      if (!characterId) return []
      return apiGet<CharacterCompetency[]>(`${API}/api/characters/${characterId}/competencies`)
    },
    enabled: !!characterId,
  })
}