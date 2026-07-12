import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

export type NPCTemplate = Readonly<{
  id: string;
  name: string;
  world_id: string;
}>;

export function useNPCTemplates() {
  const { currentWorld } = useWorldStore();
  return useQuery({
    queryKey: ["npc-templates", currentWorld],
    queryFn: async (): Promise<NPCTemplate[]> => {
      const qs = currentWorld ? `?world_id=${currentWorld}` : "";
      return apiGet<NPCTemplate[]>(`${window.location.origin}/api/npc-templates${qs}`);
    },
    enabled: !!currentWorld,
  });
}