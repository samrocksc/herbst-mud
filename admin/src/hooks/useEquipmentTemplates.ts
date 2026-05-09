import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'

const API = `${window.location.origin}/api/equipment-templates`

export type EquipmentTemplate = Readonly<{
  id: string
  name: string
  description: string
  slot: string
  level: number
  weight: number
  item_type: string
  stats: Record<string, number> | string
  color: string
  is_visible: boolean
  is_immovable: boolean
  effect_type: string
  effect_value: number
  effect_duration: number
  is_container: boolean
  container_capacity: number
  is_locked: boolean
  key_item_id: string
  reveal_condition: string
  armor_rating: number
  armor_type: string
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  damage_type: string
  weapon_type: string
  is_two_handed: boolean
}>

export type EquipmentTemplateInput = Partial<{
  name: string
  description: string
  slot: string
  level: number
  weight: number
  item_type: string
  stats: Record<string, number>
  color: string
  is_visible: boolean
  is_immovable: boolean
  effect_type: string
  effect_value: number
  effect_duration: number
  is_container: boolean
  container_capacity: number
  is_locked: boolean
  key_item_id: string
  reveal_condition: string
  armor_rating: number
  armor_type: string
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  damage_type: string
  weapon_type: string
  is_two_handed: boolean
}>

export function useEquipmentTemplates() {
  return useQuery({
    queryKey: ['equipment-templates'],
    queryFn: async (): Promise<EquipmentTemplate[]> => {
      const data = await apiGet<unknown[]>(API)
      return Array.isArray(data) ? data : []
    },
  })
}

export function useEquipmentTemplate(id: string | null) {
  return useQuery({
    queryKey: ['equipment-template', id],
    queryFn: () => apiGet<EquipmentTemplate>(`${API}/${id}`),
    enabled: !!id,
  })
}

export function useCreateTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: EquipmentTemplateInput) => apiPost<EquipmentTemplate>(API, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['equipment-templates'] }),
  })
}

export function useUpdateTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: EquipmentTemplateInput }) =>
      apiPut<EquipmentTemplate>(`${API}/${id}`, input),
    onSuccess: (_data, { id }) => {
      qc.invalidateQueries({ queryKey: ['equipment-templates'] })
      qc.invalidateQueries({ queryKey: ['equipment-template', id] })
    },
  })
}

export function useDeleteTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => apiDelete(`${API}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['equipment-templates'] }),
  })
}