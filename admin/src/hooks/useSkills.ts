import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch';
import { useWorldStore } from './useWorldStore';

const API = `${window.location.origin}/api/abilities`;

export type Skill = Readonly<{
  id: number
  name: string
  description: string
  type: 'weapon' | 'magic' | 'utility'
  tags: string[]
  primary_stat: 'STR' | 'DEX' | 'INT' | 'WIS'
  level_req: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  classless: boolean
  effects?: Record<string, unknown>
}>

export type SkillInput = Readonly<{
  id?: number
  name: string
  description: string
  type: 'weapon' | 'magic' | 'utility'
  tags: string
  primary_stat: 'STR' | 'DEX' | 'INT' | 'WIS'
  level_req: number
  cooldown: number
  mana_cost: number
  stamina_cost: number
  classless: boolean
  effects?: string
}>

function parseSkillForApi(input: SkillInput): Omit<Skill, 'id'> {
  return {
    name: input.name,
    description: input.description,
    type: input.type,
    tags: input.tags.split(',').map(t => t.trim()).filter(Boolean),
    primary_stat: input.primary_stat,
    level_req: input.level_req,
    cooldown: input.cooldown,
    mana_cost: input.mana_cost,
    stamina_cost: input.stamina_cost,
    classless: input.classless,
    effects: input.effects ? JSON.parse(input.effects) : undefined,
  };
}

export function useSkills(filters?: { type?: string; tag?: string }) {
  const { currentWorld } = useWorldStore();

  const params = new URLSearchParams();
  if (filters?.type) params.append('type', filters.type);
  if (filters?.tag) params.append('tag', filters.tag);
  if (currentWorld) params.append('world_id', currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : '';

  return useQuery({
    queryKey: ['skills', filters, currentWorld],
    queryFn: async () => {
      const data = await apiGet<Skill[]>(`${API}${qs}`);
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useSkill(id: number | null) {
  const { currentWorld } = useWorldStore();

  const params = new URLSearchParams();
  if (currentWorld) params.append('world_id', currentWorld);
  const qs = params.toString() ? `?${params.toString()}` : '';

  return useQuery({
    queryKey: ['skill', id, currentWorld],
    queryFn: () => apiGet<Skill>(`${API}/${id}${qs}`),
    enabled: !!id,
  });
}

export function useCreateSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: SkillInput) => apiPost<Skill>(API, parseSkillForApi(input)),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['skills'] }),
  });
}

export function useUpdateSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: SkillInput }) =>
      apiPut<Skill>(`${API}/${id}`, parseSkillForApi(input)),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['skills'] });
      qc.invalidateQueries({ queryKey: ['skill', id] });
    },
  });
}

export function useDeleteSkill() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['skills'] }),
  });
}
