import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type DialogResponse = Readonly<{
  label: string
  next_node_id: string
  condition?: string
  quest_offer_id?: string
  decline_node_id?: string
  effects?: number[]
}>

export type DialogNode = Readonly<{
  id: string
  npc_text: string
  responses: DialogResponse[]
  is_entry: boolean
  entry_condition?: string
  on_enter_effects: number[]
  npc_template_id: string
}>

export type DialogNodeInput = Readonly<{
  id?: string
  npc_text?: string
  responses?: DialogResponse[]
  is_entry?: boolean
  entry_condition?: string
  on_enter_effects?: number[]
  npc_template_id?: string
}>

export function useDialogNodes(templateId: string) {
  return useQuery({
    queryKey: ['dialog-nodes', templateId],
    queryFn: async (): Promise<DialogNode[]> => {
      const data = await apiGet<{ dialog_nodes: DialogNode[] }>(`${API}/api/npc-templates/${templateId}/dialog-nodes`)
      return data.dialog_nodes ?? []
    },
    enabled: !!templateId,
  })
}

export function useCreateDialogNode() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: DialogNodeInput & { npc_template_id: string }) =>
      apiPost<DialogNode>(`${API}/api/dialog-nodes`, input),
    onSuccess: (_, vars) => qc.invalidateQueries({ queryKey: ['dialog-nodes', vars.npc_template_id] }),
  })
}

export function useUpdateDialogNode() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: DialogNodeInput }) =>
      apiPut<DialogNode>(`${API}/api/dialog-nodes/${encodeURIComponent(id)}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['dialog-nodes'] }),
  })
}

export function useDeleteDialogNode() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => apiDelete(`${API}/api/dialog-nodes/${encodeURIComponent(id)}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['dialog-nodes'] }),
  })
}