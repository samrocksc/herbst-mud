/* eslint-disable functional/prefer-immutable-types */
import { useQuery } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { apiGet } from "../utils/apiFetch";

type Recipe = Readonly<{
  name: string;
  displayName: string;
  description: string;
  requiredStationTag: string;
  requiredClass: string;
  requiredSkillLevel: number;
  requiredSkill: string;
  inputs: any[];
  outputs: any[];
  craftTimeSecs: number;
  world_id: string;
}>;

export function useSearchRecipes(query: string) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["search-recipes", query, currentWorld],
    queryFn: async (): Promise<Recipe[]> => {
      if (!query || query.length < 2) return [];
      const params = new URLSearchParams({
        search: query,
        ...(currentWorld ? { world_id: currentWorld } : {}),
      });
      const url = `${window.location.origin}/api/recipes/search?${params.toString()}`;
      const data = await apiGet<{ recipes: Recipe[] }>(url);
      return Array.isArray(data?.recipes) ? data.recipes : [];
    },
    enabled: query.length >= 2,
  });
}
