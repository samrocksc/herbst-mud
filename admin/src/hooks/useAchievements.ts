import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}`;

export type Achievement = Readonly<{
  id: number
  name: string
  description: string
  icon: string
  xp_reward: number
  criteria: string
}>

export type AchievementInput = Readonly<{
  name: string
  description: string
  icon: string
  xp_reward: number
  criteria: string
}>

export function useAchievements() {
  return useQuery({
    queryKey: ['achievements'],
    queryFn: () => apiGet<Achievement[]>(`${API}/api/achievements`),
  });
}

export function useCreateAchievement() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: AchievementInput) =>
      apiPost<Achievement>(`${API}/api/achievements`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['achievements'] }),
  });
}

export function useUpdateAchievement() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: AchievementInput }) =>
      apiPut<Achievement>(`${API}/api/achievements/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['achievements'] });
      qc.invalidateQueries({ queryKey: ['achievement', id] });
    },
  });
}

export function useDeleteAchievement() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/achievements/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['achievements'] }),
  });
}