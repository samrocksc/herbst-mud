import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}`;

export const EFFECT_TYPES = [
  { value: 'hp_change', label: 'HP Change', group: 'Combat' },
  { value: 'stamina_change', label: 'Stamina Change', group: 'Combat' },
  { value: 'mana_change', label: 'Mana Change', group: 'Combat' },
  { value: 'xp_drain', label: 'XP Drain', group: 'Progression' },
  { value: 'xp_gain', label: 'XP Gain', group: 'Progression' },
  { value: 'xp_set', label: 'XP Set', group: 'Progression' },
  { value: 'bind_point_set', label: 'Bind Point Set', group: 'Location' },
  { value: 'teleport', label: 'Teleport', group: 'Location' },
  { value: 'message', label: 'Message', group: 'Messaging' },
  { value: 'room_message', label: 'Room Message', group: 'Messaging' },
  { value: 'whisper', label: 'Whisper', group: 'Messaging' },
  { value: 'tag_add', label: 'Tag Add', group: 'Tagging' },
  { value: 'tag_remove', label: 'Tag Remove', group: 'Tagging' },
  { value: 'apply_effect', label: 'Apply Effect', group: 'Chaining' },
] as const;

export const STACK_MODES = [
  { value: 'replace', label: 'Replace (single stack)' },
  { value: 'refresh', label: 'Refresh (extend duration)' },
  { value: 'stack', label: 'Stack (multiple stacks)' },
];

export const MESSAGE_TYPES = [
  { value: 'info', label: 'Info (neutral)' },
  { value: 'success', label: 'Success (green)' },
  { value: 'error', label: 'Error (red)' },
  { value: 'warn', label: 'Warn (yellow)' },
  { value: 'combat', label: 'Combat (orange)' },
  { value: 'heal', label: 'Heal (green)' },
  { value: 'system', label: 'System (bold)' },
  { value: 'say', label: 'Say (room wide)' },
  { value: 'yell', label: 'Yell (zone wide)' },
  { value: 'shout', label: 'Shout (world wide)' },
  { value: 'whisper', label: 'Whisper (private)' },
];

export type EffectDef = Readonly<{
  id: number
  name: string
  description: string
  effect_type: string
  parameters: Record<string, unknown>
  stack_mode: string
  stack_limit: number
  is_permanent: boolean
  duration_secs: number
  messages: Record<string, string>
  hook_count: number
}>

export type EffectDefInput = {
  name: string
  description: string
  effect_type: string
  parameters: Record<string, unknown>
  stack_mode: string
  stack_limit: number
  is_permanent: boolean
  duration_secs: number
  messages: Record<string, string>
}

// Empty input with sensible defaults based on effect type
export const createEmptyInput = (effectType: string): EffectDefInput => ({
  name: '',
  description: '',
  effect_type: effectType,
  parameters: {},
  stack_mode: 'replace',
  stack_limit: 1,
  is_permanent: false,
  duration_secs: 0,
  messages: {},
});

export function useEffectDefs(filters?: { type?: string }) {
  const params = new URLSearchParams();
  if (filters?.type) params.set('type', filters.type);
  const qs = params.toString() ? `?${params.toString()}` : '';
  return useQuery({
    queryKey: ['effectDefs', filters],
    queryFn: () => apiGet<EffectDef[]>(`${API}/api/effects${qs}`),
  });
}

export function useEffectDef(id: number | null) {
  return useQuery({
    queryKey: ['effectDef', id],
    queryFn: () => apiGet<EffectDef>(`${API}/api/effects/${id}`),
    enabled: !!id,
  });
}

export function useCreateEffectDef() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: EffectDefInput) => apiPost<EffectDef>(`${API}/api/effects`, input),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['effectDefs'] }); },
  });
}

export function useUpdateEffectDef() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<EffectDefInput> }) =>
      apiPut<EffectDef>(`${API}/api/effects/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['effectDefs'] });
      qc.invalidateQueries({ queryKey: ['effectDef', id] });
    },
  });
}

export function useDeleteEffectDef() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/effects/${id}`),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['effectDefs'] }); },
  });
}
