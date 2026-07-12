/* eslint-disable functional/prefer-immutable-types */
import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";

type Effect = Readonly<{
  id: number;
  name: string;
  description: string;
  effect_type: string;
  parameters: Record<string, any>;
  stack_mode: string;
  stack_limit: number;
  is_permanent: boolean;
  duration_secs: number;
}>;

export function useSearchEffects(query: string) {
  return useQuery({
    queryKey: ["search-effects", query],
    queryFn: async (): Promise<Effect[]> => {
      if (!query || query.length < 2) return [];
      const params = new URLSearchParams({
        search: query,
      });
      const url = `${window.location.origin}/api/effects/search?${params.toString()}`;
      const data = await apiGet<{ effects: Effect[] }>(url);
      return Array.isArray(data?.effects) ? data.effects : [];
    },
    enabled: query.length >= 2,
  });
}
