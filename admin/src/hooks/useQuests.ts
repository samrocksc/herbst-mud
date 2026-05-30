/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}`;

export type QuestObjective = Readonly<{
  type: string
  target_id: string
  tag_filter: string
  count: number
  labels: string[]
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
  main_type: string
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
};

export { EMPTY_REWARDS };

export function useQuests() {
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery({
    queryKey: ["quests", currentWorld],
    queryFn: async (): Promise<Quest[]> => {
      return apiGet<Quest[]>(`${API}/api/quests${qs}`);
    },
  });
}

export type QuestLookups = Readonly<{
  quest_types: { id: string; name: string }[]
  npcs: { id: string; name: string }[]
  rooms: { id: string; name: string }[]
  items: { id: string; name: string }[]
  effects: { id: string; name: string }[]
  tags: { id: string; name: string }[]
  achievements: { id: string; name: string }[]
  prerequisite_quests: { id: string; name: string }[]
}>

export function useQuestLookups() {
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery({
    queryKey: ["quest-lookups", currentWorld],
    queryFn: async (): Promise<QuestLookups> => {
      return apiGet<QuestLookups>(`${API}/api/quests/lookups${qs}`);
    },
  });
}

export function useQuest(id: number | null) {
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery({
    queryKey: ["quest", id, currentWorld],
    queryFn: async (): Promise<Quest | null> => {
      if (!id) return null;
      return apiGet<Quest>(`${API}/api/quests/${id}${qs}`);
    },
    enabled: !!id,
  });
}

export function useCreateQuest() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: QuestInput) =>
      apiPost<Quest>(`${API}/api/quests`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["quests"] }),
  });
}

export function useUpdateQuest() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: QuestInput }) =>
      apiPut<Quest>(`${API}/api/quests/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["quests"] });
      qc.invalidateQueries({ queryKey: ["quest", id] });
    },
  });
}

export function useDeleteQuest() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/quests/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["quests"] }),
  });
}
