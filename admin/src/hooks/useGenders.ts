/* eslint-disable functional/prefer-immutable-types, functional/immutable-data */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api/genders`;

export type Gender = Readonly<{
  id: number
  name: string
  display_name: string
  subject_pronoun: string
  object_pronoun: string
  possessive_pronoun: string
  world_id: string
}>

export type GenderInput = Readonly<{
  name: string
  display_name: string
  subject_pronoun: string
  object_pronoun: string
  possessive_pronoun: string
  world_id: string
}>

function parseGenderForApi(input: GenderInput, worldId?: string): Record<string, unknown> {
  const effectiveWorldId = worldId || input.world_id;
  if (!effectiveWorldId) {
    throw new Error("world_id is required but not provided");
  }
  return {
    name: input.name,
    display_name: input.display_name || input.name,
    subject_pronoun: input.subject_pronoun,
    object_pronoun: input.object_pronoun,
    possessive_pronoun: input.possessive_pronoun,
    world_id: effectiveWorldId,
  };
}

export function useGenders(worldId?: string) {
  const { currentWorld } = useWorldStore();
  const effectiveWorldId = worldId || currentWorld;

  if (!effectiveWorldId) {
    throw new Error("World ID not available - must be set via WorldStoreProvider");
  }

  return useQuery({
    queryKey: ["genders", effectiveWorldId],
    queryFn: async (): Promise<Gender[]> => {
      const url = effectiveWorldId !== "default" ? `${API}?world_id=${effectiveWorldId}` : API;
      const data = await apiGet<Gender[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCreateGender() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: (input: GenderInput) =>
      apiPost<Gender>(API, parseGenderForApi(input, currentWorld)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["genders"] }),
  });
}

export function useUpdateGender() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: GenderInput }) =>
      apiPut<Gender>(`${API}/${id}`, parseGenderForApi(input, currentWorld)),
    onSuccess: (_, variables) => {
      const worldId = variables.input.world_id || currentWorld;
      qc.invalidateQueries({ queryKey: ["genders", worldId] });
    },
  });
}

export function useDeleteGender() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (gender: Gender) => apiDelete(`${API}/${gender.id}`),
    onSuccess: (_, gender) => {
      qc.invalidateQueries({ queryKey: ["genders", gender.world_id] });
    },
  });
}
