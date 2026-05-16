/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}`;

export type DialogResponse = Readonly<{
  label: string
  next_node_id: string
  condition?: string
  quest_offer_id?: string
  decline_node_id?: string
  effects?: number[]
}>

export type DialogNode = Readonly<{
  id: string
  npc_text: string
  responses: DialogResponse[]
  is_entry: boolean
  entry_condition?: string
  on_enter_effects: number[]
  npc_template_id: string
}>

export type DialogNodeInput = Readonly<{
  id?: string
  npc_text?: string
  responses?: DialogResponse[]
  is_entry?: boolean
  entry_condition?: string
  on_enter_effects?: number[]
  npc_template_id?: string
}>

export function useDialogNodes(templateId: string) {
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";
  return useQuery({
    queryKey: ["dialog-nodes", templateId, currentWorld],
    queryFn: async (): Promise<DialogNode[]> => {
      const data = await apiGet<{ dialog_nodes: DialogNode[] }>(`${API}/api/npc-templates/${templateId}/dialog-nodes${qs}`);
      return data.dialog_nodes ?? [];
    },
    enabled: !!templateId && !!currentWorld,
  });
}

export function useCreateDialogNode() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";
  return useMutation({
    mutationFn: (input: DialogNodeInput & { npc_template_id: string }) =>
      apiPost<DialogNode>(`${API}/api/dialog-nodes${qs}`, input),
    onSuccess: (_, vars) => qc.invalidateQueries({ queryKey: ["dialog-nodes", vars.npc_template_id] }),
  });
}

export function useUpdateDialogNode() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: DialogNodeInput }) =>
      apiPut<DialogNode>(`${API}/api/dialog-nodes/${encodeURIComponent(id)}${qs}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["dialog-nodes"] }),
  });
}

export function useDeleteDialogNode() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  const params = new URLSearchParams();
  if (currentWorld) params.append("world_id", currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : "";
  return useMutation({
    mutationFn: (id: string) => apiDelete(`${API}/api/dialog-nodes/${encodeURIComponent(id)}${qs}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["dialog-nodes"] }),
  });
}
