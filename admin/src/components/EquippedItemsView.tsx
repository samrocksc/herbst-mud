import { useQuery } from '@tanstack/react-query'
import { apiGet } from '../utils/apiFetch'
import { DEFAULT_HUMANOID_SLOTS } from './equipConstants'
import { EquippedSlotGrid, UnequippedSection } from './EquippedItemsSlots'

type EquippedItem = Readonly<{
  id: number; name: string; slot: string; level: number; itemType: string
  isEquipped: boolean; damage_dice_count: number; damage_dice_sides: number
  damage_bonus: number; damage_type: string; weapon_type: string
  armor_rating: number; armor_type: string; rarity: string; equipment_template_id: string
}>

type EquippedItemsViewProps = Readonly<{ characterId: number; characterRace: string }>

/** Display equipped items organized by slot. */
export function EquippedItemsView({ characterId, characterRace }: EquippedItemsViewProps) {
  const { data: raceData, isLoading: raceLoading } = useQuery({
    queryKey: ['races'],
    queryFn: () => apiGet<Readonly<{ name: string; equipment_slots: string[] }[]>>(`${window.location.origin}/api/races`),
  })

  const { data: items, isLoading: itemsLoading, error } = useQuery<EquippedItem[]>({
    queryKey: ['item-instances', 'owner', characterId],
    queryFn: () => apiGet<EquippedItem[]>(`${window.location.origin}/api/item-instances?ownerId=${characterId}`),
  })

  const races = Array.isArray(raceData) ? raceData : []
  const raceObj = races.find((r) => r.name === characterRace)
  const slots = raceObj?.equipment_slots ?? DEFAULT_HUMANOID_SLOTS
  const equipped = (items ?? []).filter((i) => i.isEquipped)
  const unequipped = (items ?? []).filter((i) => !i.isEquipped)

  if (itemsLoading || raceLoading) return <div className="bg-surface-muted rounded-lg p-6 border border-border"><div className="text-text-muted text-sm">Loading equipment...</div></div>
  if (error) return <div className="bg-surface-muted rounded-lg p-6 border border-border"><div className="text-danger text-sm">Failed to load equipment</div></div>

  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Equipped Items</h2>
      <EquippedSlotGrid slots={slots} items={equipped} />
      {unequipped.length > 0 && <UnequippedSection items={unequipped} />}
    </div>
  )
}