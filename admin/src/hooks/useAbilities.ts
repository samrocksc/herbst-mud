 
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api/abilities`;

export type Ability = Readonly<{
  id: number
  name: string
  description: string
  ability_type: string
  cost: number
  cooldown: number
  cooldown_seconds: number
  requirements: string
  mana_cost: number
  stamina_cost: number
  hp_cost: number
  slug: string
  required_tag: string
  ability_class: string
  proc_chance: number
  proc_event: string
  faction_skills: number | null
  required_skill_id: number | null
  required_skill_level: number
}>

export type AbilityInput = Readonly<{
  id?: number
  name: string
  description: string
  ability_type: string
  requirements: string
  cost: number
  cooldown: number
  cooldown_seconds: number
  mana_cost: number
  stamina_cost: number
  hp_cost: number
  proc_chance: number
  proc_event: string
  ability_class: string
  required_tag: string
  slug?: string
  world_id?: string
  required_skill_id: number | null
  required_skill_level: number
}>

function parseForApi(input: AbilityInput, worldId?: string): Record<string, unknown> {
  return {
    name: input.name,
    description: input.description,
    ability_type: input.ability_type,
    requirements: input.requirements,
    cost: input.cost,
    cooldown: input.cooldown,
    cooldown_seconds: input.cooldown_seconds,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
    hp_cost: input.hp_cost,
    proc_chance: input.proc_chance,
    proc_event: input.proc_event,
    ability_class: input.ability_class,
    required_tag: input.required_tag,
    slug: input.slug || undefined,
    world_id: worldId || undefined,
    required_skill_id: input.required_skill_id,
    required_skill_level: input.required_skill_level,
  };
}

export function useAbilities(filters?: { type?: string; abilityClass?: string }) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["abilities", filters, currentWorld],
    queryFn: async (): Promise<Ability[]> => {
      const params = new URLSearchParams([
        ...(filters?.type ? [["type", filters.type]] : []),
        ...(filters?.abilityClass ? [["ability_class", filters.abilityClass]] : []),
        ...(currentWorld ? [["world_id", currentWorld]] : []),
      ]);
      const url = `${API}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<Ability[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useAbility(id: number | null) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["ability", id, currentWorld],
    queryFn: async (): Promise<Ability | null> => {
      if (!id) return null;
      const params = new URLSearchParams(currentWorld ? [["world_id", currentWorld]] : []);
      const url = `${API}/${id}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<Ability>(url);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useCreateAbility() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: (input: AbilityInput) =>
      apiPost<Ability>(API, parseForApi(input, currentWorld)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["abilities"] }),
  });
}

export function useUpdateAbility() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: ({ id, input }: { readonly id: number; readonly input: AbilityInput }) =>
      apiPut<Ability>(`${API}/${id}`, parseForApi(input, currentWorld)),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["abilities"] });
      qc.invalidateQueries({ queryKey: ["ability", id] });
    },
  });
}

export function useDeleteAbility() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["abilities"] }),
  });
}

// Effect management for abilities
export type AbilityEffect = Readonly<{
  id: number;
  effect_type: string;
  target: string;
  value: number;
  duration: number;
  scaling_stat: string;
  scaling_ratio: number;
  sort_order: number;
  ability_id?: number;
  effect_id?: number;
  applied_by_id?: number;
  duration_secs?: number;
  multiplier?: number;
  amount?: number;
  parameters?: Record<string, unknown>;
  stack_mode?: string;
  name?: string;
  description?: string;
  stack_count?: number;
  is_active?: boolean;
  expires_at?: string;
  effect?: {
    id: number;
    effect_type: string;
    parameters?: Record<string, unknown>;
  };
}>;

export type CreateAbilityEffectInput = Readonly<{
  effect_id: number;
  effect_type: string;
  target: string;
  value: number;
  duration?: number;
}>;

export function useCreateAbilityEffect() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: ({ abilityId, input }: { readonly abilityId: number; readonly input: CreateAbilityEffectInput }) =>
      apiPost<AbilityEffect>(`${API}/${abilityId}/effects`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["abilities"] }),
  });
}
