import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API_BASE = `${window.location.origin}`

// ─── Types ──────────────────────────────────────────────────────────────────

export type ItemInstance = Readonly<{
  id: number
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
  equipment_template_id: string
  ownerId: number | null
  effect_type: string
  effect_value: number
  effect_duration: number
  healing: number
  effect: string
  isContainer: boolean
  containerCapacity: number
  isLocked: boolean
  keyItemID: string
  containedItems: string
  revealCondition: string
}>

export type ItemInstanceCreateInput = Readonly<{
  equipment_template_id?: string
  name?: string
  description?: string
  slot?: string
  level?: number
  weight?: number
  isEquipped?: boolean
  isImmovable?: boolean
  color?: string
  isVisible?: boolean
  itemType?: string
  ownerId?: number | null
  effect_type?: string
  effect_value?: number
  effect_duration?: number
  healing?: number
  effect?: string
  isContainer?: boolean
  containerCapacity?: number
  isLocked?: boolean
  keyItemID?: string
  containedItems?: string
  revealCondition?: string
}>

export type ItemInstanceUpdateInput = Readonly<{
  name?: string
  description?: string
  slot?: string
  level?: number
  weight?: number
  isEquipped?: boolean
  isImmovable?: boolean
  color?: string
  isVisible?: boolean
  itemType?: string
  ownerId?: number | null
  effect_type?: string
  effect_value?: number
  effect_duration?: number
  healing?: number
  effect?: string
  isContainer?: boolean
  containerCapacity?: number
  isLocked?: boolean
  keyItemID?: string
  containedItems?: string
  revealCondition?: string
}>

// ─── Query hooks ────────────────────────────────────────────────────────────

export function useItemInstances() {
  return useQuery({
    queryKey: ['item-instances'],
    queryFn: async (): Promise<ItemInstance[]> => {
      const data = await apiGet<ItemInstance[]>(`${API_BASE}/api/item-instances`)
      return Array.isArray(data) ? data : []
    },
  })
}

// ─── Mutation hooks ─────────────────────────────────────────────────────────

export function useCreateItemInstance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (input: ItemInstanceCreateInput): Promise<ItemInstance> => {
      return apiPost<ItemInstance>(`${API_BASE}/api/item-instances`, input)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })
}

export function useUpdateItemInstance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, update }: { id: number; update: ItemInstanceUpdateInput }): Promise<ItemInstance> => {
      return apiPut<ItemInstance>(`${API_BASE}/api/item-instances/${id}`, update)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })
}

export function useDeleteItemInstance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      await apiDelete(`${API_BASE}/api/item-instances/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })
}