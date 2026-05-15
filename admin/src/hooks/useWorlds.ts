import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}/api`;

export type World = Readonly<{
  id: number
  name: string
  title: string
  description: string
  active: boolean
}>

export type WorldInput = Readonly<{
  id?: number
  name: string
  title: string
  description: string
  active: boolean
}>

function parseForApi(input: WorldInput): Record<string, unknown> {
  return {
    name: input.name,
    title: input.title,
    description: input.description,
    active: input.active,
  };
}

export function useWorlds() {
  return useQuery({
    queryKey: ['worlds'],
    queryFn: async (): Promise<World[]> => {
      const data = await apiGet<{ worlds: World[]; error?: string }>(`${API}/worlds/db`);
      if (!data) return [];
      if (data.error) {
        throw new Error(data.error);
      }
      return data.worlds ?? [];
    },
  });
}

export function useWorld(id: number | null) {
  return useQuery({
    queryKey: ['world', id],
    queryFn: async (): Promise<World | null> => {
      if (!id) return null;
      const data = await apiGet<World>(`${API}/worlds/${id}`);
      return data ?? null;
    },
    enabled: !!id,
  });
}

export function useCreateWorld() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: WorldInput) =>
      apiPost<World>(`${API}/worlds`, parseForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['worlds'] }),
  });
}

export function useUpdateWorld() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: WorldInput }) =>
      apiPut<World>(`${API}/worlds/${id}`, parseForApi(input)),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['worlds'] });
      qc.invalidateQueries({ queryKey: ['world', id] });
    },
  });
}

export function useSetActiveWorld() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, active }: { id: number; active: boolean }) =>
      apiPut<World>(`${API}/worlds/${id}`, { active }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['worlds'] });
      qc.invalidateQueries({ queryKey: ['world', id] });
    },
  });
}

export function useDeleteWorld() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/worlds/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['worlds'] }),
  });
}