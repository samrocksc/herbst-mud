import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

const API = `${window.location.origin}/api/recipes`;

export type CraftingInput = Readonly<{
  equipment_template_id: string
  quantity: number
  consumed: boolean
}>

export type CraftingOutput = Readonly<{
  equipment_template_id: string
  quantity: number
}>

export type Recipe = Readonly<{
  name: string
  display_name: string
  description: string
  required_station_tag: string
  required_class: string
  required_skill_level: number
  required_skill: string
  inputs: CraftingInput[]
  outputs: CraftingOutput[]
  craft_time_secs: number
  world_id: string
}>

export type RecipeInput = Readonly<{
  name: string
  display_name: string
  description: string
  required_station_tag: string
  required_class: string
  required_skill_level: number
  required_skill: string
  inputs: CraftingInput[]
  outputs: CraftingOutput[]
  craft_time_secs: number
  world_id: string
}>

export function useRecipes(filters?: { world_id?: string; station_tag?: string }) {
  return useQuery({
    queryKey: ["recipes", filters],
    queryFn: async (): Promise<Recipe[]> => {
      const params = new URLSearchParams();
      if (filters?.world_id) params.set("world_id", filters.world_id);
      if (filters?.station_tag) params.set("station_tag", filters.station_tag);
      const url = `${API}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<Recipe[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useRecipe(name: string | null) {
  return useQuery({
    queryKey: ["recipe", name],
    queryFn: async (): Promise<Recipe | null> => {
      if (!name) return null;
      const data = await apiGet<Recipe>(`${API}/${name}`);
      return data ?? null;
    },
    enabled: !!name,
  });
}

export function useCreateRecipe() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: RecipeInput) => apiPost<Recipe>(API, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["recipes"] }),
  });
}

export function useUpdateRecipe() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ name, input }: { readonly name: string; readonly input: RecipeInput }) =>
      apiPut<Recipe>(`${API}/${name}`, input),
    onSuccess: (_, { name }) => {
      qc.invalidateQueries({ queryKey: ["recipes"] });
      qc.invalidateQueries({ queryKey: ["recipe", name] });
    },
  });
}

export function useDeleteRecipe() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => apiDelete(`${API}/${name}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["recipes"] }),
  });
}