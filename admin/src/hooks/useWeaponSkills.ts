/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

const API = `${window.location.origin}`;

// Trainable skill multiplier (e.g. "Pizza Making", "One-Handed", "Fist Fights")
// Name + description + category only. Effects/costs/cooldowns belong on Abilities.
export type TrainableSkill = Readonly<{
  id: number
  name: string
  description: string
  skill_category: string
  requirements: string
}>

export type TrainableSkillInput = Readonly<{
  id?: number
  name: string
  description: string
  skill_category: string
  requirements: string
}>

export function useWeaponSkills() {
  return useQuery({
    queryKey: ["weapon-skills"],
    queryFn: () => apiGet<TrainableSkill[]>(`${API}/api/abilities?class=passive`),
  });
}

export function useWeaponSkill(id: number | null) {
  return useQuery({
    queryKey: ["weapon-skill", id],
    queryFn: () => (id ? apiGet<TrainableSkill>(`${API}/api/abilities/${id}`) : null),
    enabled: !!id,
  });
}

export function useCreateWeaponSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: TrainableSkillInput) =>
      apiPost<TrainableSkill>(`${API}/api/abilities`, { ...input, ability_class: "passive" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["weapon-skills"] }),
  });
}

export function useUpdateWeaponSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: TrainableSkillInput }) =>
      apiPut<TrainableSkill>(`${API}/api/abilities/${id}`, { ...input, ability_class: "passive" }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["weapon-skills"] });
      qc.invalidateQueries({ queryKey: ["weapon-skill", id] });
    },
  });
}

export function useDeleteWeaponSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/abilities/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["weapon-skills"] }),
  });
}