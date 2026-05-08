import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type CompetencyCategory = Readonly<{
  id: string
  name: string
  xp_multiplier: number
  thresholds: CompetencyThreshold[]
}>

export type CompetencyThreshold = Readonly<{
  level: number
  xp_required: number
  damage_multiplier: number
  defense_multiplier: number
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
    queryFn: async (): Promise<CompetencyCategory[]> =>
      apiGet<CompetencyCategory[]>(`${API}/api/competency-categories`),
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