import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

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

const EMPTY_INPUT: EffectDefInput = {
  name: '',
  description: '',
  effect_type: 'hp_change',
  parameters: {},
  stack_mode: 'replace',
  stack_limit: 1,
  is_permanent: false,
  duration_secs: 0,
  messages: {},
}

export { EMPTY_INPUT }

export function useEffectDefs(filters?: { type?: string }) {
  const params = new URLSearchParams()
  if (filters?.type) params.set('type', filters.type)
  const qs = params.toString() ? `?${params.toString()}` : ''
  return useQuery({
    queryKey: ['effectDefs', filters],
    queryFn: () => apiGet<EffectDef[]>(`${API}/api/effects${qs}`),
  })
}

export function useEffectDef(id: number | null) {
  return useQuery({
    queryKey: ['effectDef', id],
    queryFn: () => apiGet<EffectDef>(`${API}/api/effects/${id}`),
    enabled: !!id,
  })
}

export function useCreateEffectDef() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: EffectDefInput) => apiPost<EffectDef>(`${API}/api/effects`, input),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['effectDefs'] }) },
  })
}

export function useUpdateEffectDef() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<EffectDefInput> }) =>
      apiPut<EffectDef>(`${API}/api/effects/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ['effectDefs'] })
      qc.invalidateQueries({ queryKey: ['effectDef', id] })
    },
  })
}

export function useDeleteEffectDef() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/effects/${id}`),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['effectDefs'] }) },
  })
}