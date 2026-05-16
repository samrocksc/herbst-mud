import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";

type SearchOption = Readonly<{ id: number | string; name: string }>

export function useResourceSearch(
  resourceType: string,
  apiBase: string,
  query: string,
) {
  return useQuery({
    queryKey: ["resource-search", resourceType, query],
    queryFn: () => apiGet<SearchOption[]>(
      `${apiBase}/${resourceType}?search=${encodeURIComponent(query)}`,
    ),
    enabled: query.length >= 2,
  });
}