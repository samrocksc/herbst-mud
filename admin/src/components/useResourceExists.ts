import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";

type ExistsResult = Readonly<{
  exists: boolean | null
  name: string | null
  isValidating: boolean
}>

export function useResourceExists(
  resourceType: string,
  apiBase: string,
  id: number | string | null | undefined,
): ExistsResult {
  const numericId = typeof id === "string" ? parseInt(id, 10) : id;
  const enabled = id != null && id !== "" && numericId !== 0 && !isNaN(numericId as number);

  const query = useQuery({
    queryKey: ["resource-exists", resourceType, numericId],
    queryFn: () => apiGet<{ id: number; name: string }>(
      `${apiBase}/${resourceType}/${numericId}`,
    ),
    enabled,
    retry: false,
  });

  if (!enabled) {
    return { exists: null, name: null, isValidating: false };
  }

  if (query.isLoading) {
    return { exists: null, name: null, isValidating: true };
  }

  if (query.isError) {
    return { exists: false, name: null, isValidating: false };
  }

  const data = query.data as { id: number; name: string } | undefined;
  return {
    exists: !!data,
    name: data?.name ?? null,
    isValidating: false,
  };
}