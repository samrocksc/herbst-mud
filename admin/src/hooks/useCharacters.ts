/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut, apiPost, apiDelete } from "../utils/apiFetch";

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
  lastSeenAt: string | null
}>

export type CharacterUpdate = Partial<{
  name: string
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
  return useQuery({
    queryKey: ["characters"],
    queryFn: async (): Promise<Character[]> => {
      const data = await apiGet<Character[]>(`${API}/characters`);
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