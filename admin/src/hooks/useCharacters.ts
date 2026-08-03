/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut, apiPost, apiDelete, buildWorldParams } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

export type Character = Readonly<{
  id: number
  name: string
  isNPC: boolean
  currentRoomId: number
  startingRoomId: number
  respawnRoomId: number
  is_admin: boolean
  is_immortal: boolean
  is_test: boolean
  currentWorld: string
  hitpoints: number
  max_hitpoints: number
  stamina: number
  max_stamina: number
  mana: number
  max_mana: number
  race: string
  class: string
  gender: string
  description: string
  level: number
  xp: number
  gold_credits: number
  strength: number
  dexterity: number
  constitution: number
  intelligence: number
  wisdom: number
  charisma: number
  lastSeenAt: string | null
}>

export type CharacterUpdate = Partial<{
  name: string
  currentWorld: string
  isNPC: boolean
  currentRoomId: number
  startingRoomId: number
  respawnRoomId: number
  isAdmin: boolean
  isTest?: boolean
  gender: string
  description: string
  level: number
  xp: number
  gold_credits: number
  hitpoints: number
  maxHitpoints: number
  stamina: number
  maxStamina: number
  mana: number
  maxMana: number
}>

const API = `${window.location.origin}`;

export function useCharacters() {
  const { currentWorld } = useWorldStore();
  return useQuery({
    queryKey: ["characters", currentWorld],
    queryFn: async (): Promise<Character[]> => {
      // "default" is the initial localStorage value before the user picks a world.
      // Map it to world 1 (herbst-mud) so the list isn't empty on first load.
      const worldId = currentWorld === "default" ? "1" : currentWorld;
      const qs = buildWorldParams(worldId);
      const suffix = qs ? `?${qs}` : "";
      const data = await apiGet<Character[]>(`${API}/characters${suffix}`);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCharacter(id: number) {
  return useQuery({
    queryKey: ["character", id],
    queryFn: async (): Promise<Character | null> => {
      const data = await apiGet<Character>(`${API}/characters/${id}`);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useUpdateCharacter() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, update }: { id: number; update: CharacterUpdate }) =>
      apiPut(`${API}/characters/${id}`, update),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["characters"] });
      queryClient.invalidateQueries({ queryKey: ["character"] });
    },
  });
}

export function useCharacterGold(id: number) {
  return useQuery({
    queryKey: ["character-gold", id],
    queryFn: async (): Promise<{ character_id: number; gold_credits: number }> =>
      apiGet(`${API}/characters/${id}/gold`),
    enabled: !!id,
  });
}

export function useAddCharacterGold() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, amount, source }: { id: number; amount: number; source?: string }) =>
      apiPost(`${API}/characters/${id}/gold`, { amount, source }),
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: ["character-gold", vars.id] });
      queryClient.invalidateQueries({ queryKey: ["character", vars.id] });
    },
  });
}

export type CharacterSkills = Readonly<{
  id: number
  name: string
  skills: Readonly<Record<string, Readonly<{ level: number; bonus: number }>>>
  faction_abilities: ReadonlyArray<Readonly<{
    id: number
    name: string
    slug: string
    ability_type: string
    eligible: boolean
    reason: string
  }>>
}>

export function useCharacterSkills(id: number) {
  return useQuery({
    queryKey: ["character-skills", id],
    queryFn: async (): Promise<CharacterSkills | null> => {
      const data = await apiGet<CharacterSkills>(`${API}/characters/${id}/skills`);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useSpendCharacterGold() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, amount, reason }: { id: number; amount: number; reason?: string }) =>
      apiDelete(`${API}/characters/${id}/gold`, { amount, reason }),
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: ["character-gold", vars.id] });
      queryClient.invalidateQueries({ queryKey: ["character", vars.id] });
    },
  });
}

// --- Biography ---

export type BiographyEntry = Readonly<{
  id: number;
  type: string;
  description: string;
  timestamp: string;
  details: Record<string, unknown>;
}>;

export type BiographySummary = Readonly<{
  total_events: number;
  boss_kills: number;
  npc_kills: number;
  deaths: number;
  level_ups: number;
  class_changes: number;
  quest_completions: number;
  skills_learned: number;
}>;

export type BiographyResponse = Readonly<{
  character_id: number;
  character_name: string;
  entries: BiographyEntry[];
  summary: BiographySummary;
}>;

export function useCharacterBiography(id: number, typeFilter?: string) {
  return useQuery({
    queryKey: ["character-biography", id, typeFilter],
    queryFn: async (): Promise<BiographyResponse> => {
      const params = typeFilter ? `?type=${typeFilter}` : "";
      return apiGet<BiographyResponse>(`${API}/characters/${id}/biography${params}`);
    },
    enabled: !!id,
  });
}