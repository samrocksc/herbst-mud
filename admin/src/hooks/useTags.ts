/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api/tags`;

export type Tag = Readonly<{ id: number; name: string; color: string; world_id?: string }>
export type TagInput = Readonly<{ name: string; color: string; world_id?: string }>

export type TagUsage = Readonly<{ id: number; name: string; type: string }>
export type TagUsageReport = Readonly<{
  tag_name: string
  total_usages: number
  abilities: TagUsage[]
  factions: TagUsage[]
  characters: TagUsage[]
}>

export function useTags() {
  const { currentWorld } = useWorldStore();
  return useQuery({
    queryKey: ["tags", currentWorld],
    queryFn: async (): Promise<Tag[]> => {
      // Pass world_id when set; backend defaults to "1" when missing or "default".
      const worldId = currentWorld && currentWorld !== "default" ? currentWorld : "1";
      const data = await apiGet<Tag[]>(`${API}?world_id=${encodeURIComponent(worldId)}`);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCreateTag() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  return useMutation({
    mutationFn: (input: TagInput) => {
      const worldId = input.world_id ?? (currentWorld && currentWorld !== "default" ? currentWorld : "1");
      return apiPost<Tag>(`${API}?world_id=${encodeURIComponent(worldId)}`, { ...input, world_id: worldId });
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags"] }),
  });
}

export function useUpdateTag() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<TagInput> }) =>
      apiPut<Tag>(`${API}/${id}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags"] }),
  });
}

export function useDeleteTag() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags"] }),
  });
}

export function useTagUsages(id: number | null) {
  return useQuery({
    queryKey: ["tag-usages", id],
    queryFn: async (): Promise<TagUsageReport> =>
      apiGet<TagUsageReport>(`${API}/${id}/usages`),
    enabled: !!id,
  });
}