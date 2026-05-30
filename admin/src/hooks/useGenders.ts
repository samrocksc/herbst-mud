/* eslint-disable functional/prefer-immutable-types, functional/immutable-data */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

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

function parseGenderForApi(input: GenderInput): Record<string, unknown> {
  return {
    name: input.name,
    display_name: input.display_name || input.name,
    subject_pronoun: input.subject_pronoun,
    object_pronoun: input.object_pronoun,
    possessive_pronoun: input.possessive_pronoun,
    world_id: input.world_id,
  };
}

export function useGenders() {
  return useQuery({
    queryKey: ["genders"],
    queryFn: async (): Promise<Gender[]> => {
      const data = await apiGet<Gender[]>(API);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCreateGender() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: GenderInput) =>
      apiPost<Gender>(API, parseGenderForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["genders"] }),
  });
}

export function useUpdateGender() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: GenderInput }) =>
      apiPut<Gender>(`${API}/${id}`, parseGenderForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["genders"] }),
  });
}

export function useDeleteGender() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["genders"] }),
  });
}
