/* eslint-disable functional/prefer-immutable-types */
import { useQuery } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { apiGet } from "../utils/apiFetch";

type DialogNode = Readonly<{
  id: string;
  npcTemplateID: string;
  npcText: string;
  responses: any[];
  isEntry: boolean;
  entryCondition: string;
  onEnterEffects: number[];
  world_id: string;
  npcTemplate?: {
    id: string;
    name: string;
    slug: string;
  };
}>;

export function useSearchDialogNodes(query: string) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["search-dialog-nodes", query, currentWorld],
    queryFn: async (): Promise<DialogNode[]> => {
      if (!query || query.length < 2) return [];
      const params = new URLSearchParams({
        search: query,
        ...(currentWorld ? { world_id: currentWorld } : {}),
      });
      const url = `${window.location.origin}/api/dialog-nodes/search?${params.toString()}`;
      const data = await apiGet<{ dialog_nodes: DialogNode[] }>(url);
      return Array.isArray(data?.dialog_nodes) ? data.dialog_nodes : [];
    },
    enabled: query.length >= 2,
  });
}
