 
import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";
import type { NPC } from "../components/map/types";

const API = `${window.location.origin}`;

export function useNPCs() {
  const { currentWorld } = useWorldStore();

  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";

  return useQuery<NPC[]>({
    queryKey: ["npcs", currentWorld],
    queryFn: async () => {
      const data = await apiGet<NPC[]>(`${API}/npcs${qs}`);
      return Array.isArray(data) ? data : [];
    },
  });
}
