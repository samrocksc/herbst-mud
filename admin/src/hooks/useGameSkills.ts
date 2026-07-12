/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api/skills`;

export type GameSkill = Readonly<{
  id: number
  world_id: string
  name: string
  display_name: string
  description: string
  category: string
  max_level: number
  xp_curve_mode: string
  xp_curve_data: Record<string, unknown>
}>

export type GameSkillInput = Readonly<{
  id?: number
  world_id?: string
  name: string
  display_name: string
  description: string
  category: string
  max_level: number
  xp_curve_mode: string
  xp_curve_data: Record<string, unknown>
}>

function parseForApi(input: GameSkillInput, worldId?: string): Record<string, unknown> {
  return {
    name: input.name,
    display_name: input.display_name,
    description: input.description,
    category: input.category,
    max_level: input.max_level,
    xp_curve_mode: input.xp_curve_mode,
    xp_curve_data: input.xp_curve_data,
    world_id: worldId || undefined,
  };
}

export function useGameSkills(filters?: { category?: string }) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["game-skills", filters, currentWorld],
    queryFn: async (): Promise<GameSkill[]> => {
      const params = new URLSearchParams([
        ...(filters?.category ? [["category", filters.category]] : []),
        ...(currentWorld ? [["world_id", currentWorld]] : []),
      ]);
      const url = `${API}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<GameSkill[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useGameSkill(id: number | null) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["game-skill", id, currentWorld],
    queryFn: async (): Promise<GameSkill | null> => {
      if (!id) return null;
      const params = new URLSearchParams(currentWorld ? [["world_id", currentWorld]] : []);
      const url = `${API}/${id}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<GameSkill>(url);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useCreateGameSkill() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: (input: GameSkillInput) =>
      apiPost<GameSkill>(API, parseForApi(input, currentWorld)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["game-skills"] }),
  });
}

export function useUpdateGameSkill() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: ({ id, input }: { readonly id: number; readonly input: GameSkillInput }) =>
      apiPut<GameSkill>(`${API}/${id}`, parseForApi(input, currentWorld)),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["game-skills"] });
      qc.invalidateQueries({ queryKey: ["game-skill", id] });
    },
  });
}

export function useDeleteGameSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["game-skills"] }),
  });
}