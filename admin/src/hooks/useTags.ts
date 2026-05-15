import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}/api/tags`;

export type Tag = Readonly<{ id: number; name: string; color: string }>
export type TagInput = Readonly<{ name: string; color: string }>

export type TagUsage = Readonly<{ id: number; name: string; type: string }>
export type TagUsageReport = Readonly<{
  tag_name: string
  total_usages: number
  abilities: TagUsage[]
  factions: TagUsage[]
  characters: TagUsage[]
}>

export function useTags() {
  return useQuery({
    queryKey: ['tags'],
    queryFn: async (): Promise<Tag[]> => {
      const data = await apiGet<Tag[]>(API);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useCreateTag() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: TagInput) => apiPost<Tag>(API, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tags'] }),
  });
}

export function useUpdateTag() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<TagInput> }) =>
      apiPut<Tag>(`${API}/${id}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tags'] }),
  });
}

export function useDeleteTag() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tags'] }),
  });
}

export function useTagUsages(id: number | null) {
  return useQuery({
    queryKey: ['tag-usages', id],
    queryFn: async (): Promise<TagUsageReport> =>
      apiGet<TagUsageReport>(`${API}/${id}/usages`),
    enabled: !!id,
  });
}