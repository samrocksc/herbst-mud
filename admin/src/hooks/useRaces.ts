/* eslint-disable functional/prefer-immutable-types, functional/immutable-data */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api/races`;

export type Race = Readonly<{
  id: number
  name: string
  display_name: string
  description: string
  stat_modifiers: Record<string, unknown> | null
  stat_growth_multipliers: Record<string, number> | null
  skill_grants: string[]
  ability_modifiers: string[]
  equipment_slots: string[]
  requirement_tags: string[]
  color: string
  tags: string[]
  resistances: Record<string, number> | null
  vulnerabilities: Record<string, number> | null
  world_id: string
}>

export type RaceInput = Readonly<{
  name: string
  display_name: string
  description: string
  stat_modifiers: string
  stat_growth_multipliers: Readonly<{ hp: number; mana: number; stamina: number }>
  skill_grants: ReadonlyArray<string>
  equipment_slots: ReadonlyArray<string>
  requirement_tags: ReadonlyArray<string>
  color: string
  tags: ReadonlyArray<string>
  resistances: Readonly<Record<string, number>>
  vulnerabilities: Readonly<Record<string, number>>
  world_id: string
}>

function parseRaceForApi(input: RaceInput, worldId?: string): Record<string, unknown> {
  const effectiveWorldId = worldId || input.world_id;
  if (!effectiveWorldId) {
    throw new Error("world_id is required but not provided");
  }

  const equipmentSlots: string[] = [...input.equipment_slots];
  const tags: string[] = [...input.tags];
  const reqTags: string[] = [...input.requirement_tags];
  const skillGrants: string[] = [...input.skill_grants];
  const body: Record<string, unknown> = {
    name: input.name,
    display_name: input.display_name || input.name,
    description: input.description,
    skill_grants: skillGrants,
    equipment_slots: equipmentSlots,
    requirement_tags: reqTags,
    color: input.color,
    tags: tags,
    resistances: { ...input.resistances },
    vulnerabilities: { ...input.vulnerabilities },
    world_id: effectiveWorldId,
  };
  if (input.stat_modifiers.trim()) {
    body.stat_modifiers = input.stat_modifiers;
  }
  body.stat_growth_multipliers = input.stat_growth_multipliers;
  return body;
}

export function useRaces() {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["races", currentWorld],
    queryFn: async (): Promise<Race[]> => {
      // Always pass world_id — backend defaults to "1" for empty/"default" but
      // some routes (e.g. /api/races) return 404 if world_id is missing entirely.
      const worldId = currentWorld && currentWorld !== "default" ? currentWorld : "1";
      const url = `${API}?world_id=${encodeURIComponent(worldId)}`;
      const data = await apiGet<Race[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCreateRace() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: (input: RaceInput) => {
      const worldId = currentWorld || input.world_id || "1";
      return apiPost<Race>(`${API}?world_id=${encodeURIComponent(worldId)}`, parseRaceForApi(input, worldId));
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["races"] }),
  });
}

export function useUpdateRace() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: RaceInput }) => {
      const worldId = input.world_id || currentWorld || "1";
      return apiPut<Race>(`${API}/${id}?world_id=${encodeURIComponent(worldId)}`, parseRaceForApi(input, worldId));
    },
    onSuccess: (_, variables) => {
      const worldId = variables.input.world_id || currentWorld || "1";
      qc.invalidateQueries({ queryKey: ["races", worldId] });
    },
  });
}

export function useDeleteRace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (race: Race) => apiDelete(`${API}/${race.id}`),
    onSuccess: (_, race) => {
      qc.invalidateQueries({ queryKey: ["races", race.world_id] });
    },
  });
}

export function useApplyRaceTags() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      apiPost<{ race: string; characters_updated: number; tags_applied: string[] }>(`${API}/${id}/apply-tags`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["races"] });
      qc.invalidateQueries({ queryKey: ["characters"] });
    },
  });
}
