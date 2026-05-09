import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type EffectHook = Readonly<{
  id: number
  name: string
  event: string
  target: string
  condition: string
  enabled: boolean
  effect_id: number
  effect_name: string
  npc_template_id: string
  npc_template_name: string
}>

export type HookInput = {
  name: string
  event: string
  target: string
  condition?: string
  enabled: boolean
  effect_id: number
}

export function useHooks() {
  return useQuery({
    queryKey: ['hooks'],
    queryFn: () => apiGet<EffectHook[]>(`${API}/api/hooks`),
  })
}

export function useTemplateHooks(npcTemplateId: string | null) {
  return useQuery({
    queryKey: ['hooks', 'template', npcTemplateId],
    queryFn: () => apiGet<EffectHook[]>(`${API}/api/npc-templates/${npcTemplateId}/hooks`),
    enabled: !!npcTemplateId,
  })
}

export function useCreateHook() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ templateId, input }: { templateId: string; input: HookInput }) =>
      apiPost<EffectHook>(`${API}/api/npc-templates/${templateId}/hooks`, input),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['hooks'] }) },
  })
}

export function useUpdateHook() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<HookInput> }) =>
      apiPut<EffectHook>(`${API}/api/hooks/${id}`, input),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['hooks'] }) },
  })
}

export function useDeleteHook() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/api/hooks/${id}`),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['hooks'] }) },
  })
}