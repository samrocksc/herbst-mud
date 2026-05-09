import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type QuestObjective = Readonly<{
  type: string
  target_id: string
  count: number
  label: string
  hint: string
}>

export type QuestRewards = Readonly<{
  xp: number
  item_ids: string[]
  effect_ids: number[]
  tag_adds: string[]
  tag_removes: string[]
  achievement_ids: number[]
}>

export type Quest = Readonly<{
  id: number
  name: string
  description: string
  prerequisite_quest_ids: string[]
  objectives: QuestObjective[]
  rewards: QuestRewards
  repeat_mode: string
  cooldown_hours: number
  is_active: boolean
}>

export type QuestInput = Readonly<{
  name?: string
  description?: string
  prerequisite_quest_ids?: string[]
  objectives?: QuestObjective[]
  rewards?: QuestRewards
  repeat_mode?: string
  cooldown_hours?: number
  is_active?: boolean
}>

const EMPTY_REWARDS: QuestRewards = {
  xp: 0, item_ids: [], effect_ids: [],
  tag_adds: [], tag_removes: [], achievement_ids: [],
}

export { EMPTY_REWARDS }

export function useQuests() {
  return useQuery({
    queryKey: ['quests'],
    queryFn: async (): Promise<Quest[]> => {
      const data = await apiGet<{ quests: Quest[] }>(`${API}/api/quests`)
      return data.quests ?? []
    },
  })
}

export function useQuest(id: number | null) {
  return useQuery({
    queryKey: ['quest', id],
    queryFn: async (): Promise<Quest | null> => {
      if (!id) return null
      return apiGet<Quest>(`${API}/api/quests/${id}`)
    },
    enabled: !!id,
  })
}

export function useCreateQuest() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: QuestInput) =>
      apiPost<Quest>(`${API}/api/quests`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['quests'] }),
  })
}

export function useUpdateQuest() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: QuestInput }) =>
      apiPut<Quest>(`${API}/api/quests/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['quests'] })
      qc.invalidateQueries({ queryKey: ['quest', id] })
    },
  })
}

export function useDeleteQuest() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/quests/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['quests'] }),
  })
}