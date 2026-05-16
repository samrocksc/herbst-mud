/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API_BASE = `${window.location.origin}`;

// ─── Types ──────────────────────────────────────────────────────────────────

export type ItemInstance = Readonly<{
  id: number
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
  equipment_template_id: string
  ownerId: number | null
  roomId: number
  effect_type: string
  effect_value: number
  effect_duration: number
  healing: number
  effect: string
  isContainer: boolean
  containerCapacity: number
  isLocked: boolean
  keyItemID: string
  containedItems: string
  revealCondition: string
  // Combat fields
  armor_rating: number
  armor_type: string
  stats: Record<string, number> | null
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  damage_type: string
  weapon_type: string
  is_two_handed: boolean
}>

export type ItemInstanceCreateInput = Partial<{
  equipment_template_id: string
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
  ownerId: number | null
  room_id: number
  effect_type: string
  effect_value: number
  effect_duration: number
  healing: number
  effect: string
  isContainer: boolean
  containerCapacity: number
  isLocked: boolean
  keyItemID: string
  containedItems: string
  revealCondition: string
  armor_rating: number
  armor_type: string
  stats: Record<string, number>
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  damage_type: string
  weapon_type: string
  is_two_handed: boolean
}>

export type ItemInstanceUpdateInput = Partial<{
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
  ownerId: number | null
  effect_type: string
  effect_value: number
  effect_duration: number
  healing: number
  effect: string
  isContainer: boolean
  containerCapacity: number
  isLocked: boolean
  keyItemID: string
  containedItems: string
  revealCondition: string
  armor_rating: number
  armor_type: string
  stats: Record<string, number>
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  damage_type: string
  weapon_type: string
  is_two_handed: boolean
}>

// ─── Query hooks ────────────────────────────────────────────────────────────

export function useItemInstances() {
  const { currentWorld } = useWorldStore();

  const params = new URLSearchParams(currentWorld ? [["world_id", currentWorld]] : []);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery({
    queryKey: ["item-instances", currentWorld],
    queryFn: async (): Promise<ItemInstance[]> => {
      const data = await apiGet<ItemInstance[]>(`${API_BASE}/api/item-instances${qs}`);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useItemInstance(id: number | null) {
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams(currentWorld ? [["world_id", currentWorld]] : []);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery({
    queryKey: ["item-instance", id, currentWorld],
    queryFn: () => apiGet<ItemInstance>(`${API_BASE}/api/item-instances/${id}${qs}`),
    enabled: !!id,
  });
}

// ─── Mutation hooks ─────────────────────────────────────────────────────────

export function useCreateItemInstance() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: ItemInstanceCreateInput) =>
      apiPost<ItemInstance>(`${API_BASE}/api/item-instances`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["item-instances"] }),
  });
}

export function useUpdateItemInstance() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, update }: { id: number; update: ItemInstanceUpdateInput }) =>
      apiPut<ItemInstance>(`${API_BASE}/api/item-instances/${id}`, update),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["item-instances"] });
      qc.invalidateQueries({ queryKey: ["item-instance"] });
    },
  });
}

export function useDeleteItemInstance() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API_BASE}/api/item-instances/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["item-instances"] }),
  });
}
