/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

const API_BASE = `${window.location.origin}`;

// ─── Types ──────────────────────────────────────────────────────────────────

export type NPCInstance = Readonly<{
  id: number
  name: string
  npc_template_id: string
  instance_number: number
  room_id: number
  starting_room_id: number
  level: number
  race: string
  hitpoints: number
  max_hitpoints: number
  stamina: number
  max_stamina: number
  mana: number
  max_mana: number
  isNPC: boolean
  is_instance: boolean
}>

export type NPCInstanceInput = Readonly<{
  template_id: string
  room_id: number
  instance_number?: number
}>

export type NPCInstanceUpdate = Readonly<{
  room_id?: number
  starting_room_id?: number
  hitpoints?: number
  is_instance?: boolean
  instance_number?: number
}>

// ─── Query hooks ────────────────────────────────────────────────────────────

export function useNPCInstances(roomId?: number) {
  const params = roomId ? `?roomId=${roomId}` : "";
  return useQuery({
    queryKey: roomId ? ["npc-instances", roomId] : ["npc-instances"],
    queryFn: (): Promise<NPCInstance[]> => {
      return apiGet<NPCInstance[]>(`${API_BASE}/api/npc-instances${params}`);
    },
  });
}

// ─── Mutation hooks ─────────────────────────────────────────────────────────

export function useCreateNPCInstance() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: NPCInstanceInput): Promise<NPCInstance> => {
      return apiPost<NPCInstance>(`${API_BASE}/api/npc-instances`, input);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-instances"] });
    },
  });
}

export function useUpdateNPCInstance() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, update }: { id: number; update: NPCInstanceUpdate }): Promise<NPCInstance> => {
      return apiPut<NPCInstance>(`${API_BASE}/api/npc-instances/${id}`, update);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-instances"] });
    },
  });
}

export function useDeleteNPCInstance() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number): Promise<NPCInstance> => {
      return apiDelete<NPCInstance>(`${API_BASE}/api/npc-instances/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-instances"] });
    },
  });
}