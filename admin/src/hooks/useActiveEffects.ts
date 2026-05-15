import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}`;

export type ActiveEffect = Readonly<{
  id: number
  character_id: number
  effect_id: number
  effect_name: string
  effect_type: string
  parameters: Record<string, unknown>
  applied_by_id: number
  stack_count: number
  started_at: string
  expires_at: string | null
  is_active: boolean
}>

export function useActiveEffects(characterId: number | null) {
  return useQuery({
    queryKey: ['activeEffects', characterId],
    queryFn: () => apiGet<ActiveEffect[]>(`${API}/api/characters/${characterId}/effects`),
    enabled: !!characterId,
  });
}

export function useRemoveActiveEffect() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ characterId, effectId }: { characterId: number; effectId: number }) =>
      apiDelete(`${API}/api/characters/${characterId}/effects/${effectId}`),
    onSuccess: (_, { characterId }) => {
      qc.invalidateQueries({ queryKey: ['activeEffects', characterId] });
    },
  });
}

export function useApplyEffect() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ characterId, effectId, appliedById }: { characterId: number; effectId: number; appliedById?: number }) =>
      apiPost<ActiveEffect>(`${API}/api/characters/${characterId}/effects/apply`, {
        effect_id: effectId,
        applied_by_id: appliedById ?? 0,
      }),
    onSuccess: (_, { characterId }) => {
      qc.invalidateQueries({ queryKey: ['activeEffects', characterId] });
    },
  });
}