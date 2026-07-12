/* eslint-disable react-refresh/only-export-components */
import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

export type ClassOption = Readonly<{
  id: number
  name: string
  display_name?: string
  description?: string
  world_id: string
}>

async function fetchClasses(worldId?: string): Promise<ClassOption[]> {
  const url = worldId ? `/classes?world_id=${worldId}` : "/classes";
  const data = await apiGet<{ classes: ClassOption[]; count: number }>(url);
  // The API returns { classes: [...], count: N }; unwrap to the list.
  if (Array.isArray(data)) return data as unknown as ClassOption[];
  return data.classes ?? [];
}

export function useClasses() {
  const { currentWorld } = useWorldStore();
  return useQuery({
    queryKey: ["classes", currentWorld],
    queryFn: () => fetchClasses(currentWorld),
  });
}