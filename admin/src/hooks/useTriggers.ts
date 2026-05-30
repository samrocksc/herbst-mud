import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";
import { resourceEndpoints } from "../utils/resourceEndpoints";

export type Trigger = Readonly<{
  id: number
  name: string
  world_id: string
  trigger_type: string
  target_type: string
  target_id: number
  room_id: number | null
  equipment_id: number | null
  condition: string
  enabled: boolean
}>

export type TriggerInput = Readonly<{
  name: string
  world_id: string
  trigger_type: string
  target_type: string
  target_id: number
  room_id: number | null
  equipment_id: number | null
  condition: string
  enabled: boolean
}>

export function useTriggers(filters?: { world_id?: string }) {
  return useQuery({
    queryKey: ["triggers", filters],
    queryFn: async (): Promise<Trigger[]> => {
      const params = new URLSearchParams();
      if (filters?.world_id) params.set("world_id", filters.world_id);
      const url = `${resourceEndpoints.triggers}${params.toString() ? "?" + params.toString() : ""}`;
      const data = await apiGet<Trigger[]>(url);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useTrigger(id: number | null) {
  return useQuery({
    queryKey: ["trigger", id],
    queryFn: async (): Promise<Trigger | null> => {
      if (!id) return null;
      const data = await apiGet<Trigger>(`${resourceEndpoints.triggers}/${id}`);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useCreateTrigger() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: TriggerInput) => apiPost<Trigger>(resourceEndpoints.triggers, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["triggers"] }),
  });
}

export function useUpdateTrigger() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { readonly id: number; readonly input: TriggerInput }) =>
      apiPut<Trigger>(`${resourceEndpoints.triggers}/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["triggers"] });
      qc.invalidateQueries({ queryKey: ["trigger", id] });
    },
  });
}

export function useDeleteTrigger() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${resourceEndpoints.triggers}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["triggers"] }),
  });
}
