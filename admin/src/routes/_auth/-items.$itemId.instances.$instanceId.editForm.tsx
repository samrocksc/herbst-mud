import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { apiPut } from '../../utils/apiFetch'
import { Button } from '../../components/Button'
import { CombatFieldsEditor, type CombatFields } from '../../components/CombatFieldsEditor'
import { NumberField, SelectField, CheckboxField } from '../../components/FormFields'
import { SLOT_OPTIONS, ITEM_TYPE_OPTIONS } from '../../components/itemConstants'
import type { ItemInstance } from '../../hooks/useItemInstances'

type InstanceEditFormState = Readonly<{
  name: string; description: string; slot: string; itemType: string
  level: number; weight: number; color: string; ownerId: number | null
  roomId: number | null; isVisible: boolean; isImmovable: boolean; isEquipped: boolean
}> & CombatFields

export function InstanceEditForm({ instance, instanceId, onDone }: Readonly<{
  instance: ItemInstance; instanceId: string; onDone: () => void
}>) {
  const queryClient = useQueryClient()
  const [form, setForm] = useState<InstanceEditFormState>(() => ({
    name: instance.name, description: instance.description, slot: instance.slot,
    itemType: instance.itemType, level: instance.level, weight: instance.weight,
    color: instance.color, ownerId: instance.ownerId, roomId: instance.roomId,
    isVisible: instance.isVisible, isImmovable: instance.isImmovable, isEquipped: instance.isEquipped,
    armor_rating: instance.armor_rating, armor_type: instance.armor_type, rarity: instance.rarity,
    skill_requirement: instance.skill_requirement, skill_requirement_level: instance.skill_requirement_level,
    damage_dice_count: instance.damage_dice_count, damage_dice_sides: instance.damage_dice_sides,
    damage_bonus: instance.damage_bonus, damage_type: instance.damage_type,
    weapon_type: instance.weapon_type, is_two_handed: instance.is_two_handed,
    stats: instance.stats ? JSON.stringify(instance.stats) : '{}',
  }))

  const updateMutation = useMutation({
    mutationFn: (body: Record<string, unknown>) => apiPut(`${window.location.origin}/api/item-instances/${instanceId}`, body),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances', instanceId] }); queryClient.invalidateQueries({ queryKey: ['item-instances'] }); onDone() },
  })

  const set = <K extends keyof InstanceEditFormState>(key: K, value: InstanceEditFormState[K]) =>
    setForm(prev => ({ ...prev, [key]: value }))

  const handleSave = () => {
    const body: Record<string, unknown> = { ...form }
    try { body.stats = JSON.parse(form.stats) } catch { body.stats = {} }
    updateMutation.mutate(body)
  }

  return (
    <div className="max-w-2xl">
      <div className="bg-surface-muted rounded-lg p-6 border border-border">
        <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Instance</h2>
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div><label className="text-text-muted text-xs block mb-1">Name</label><input type="text" value={form.name} onChange={(e) => set('name', e.target.value)} className="w-full p-2 bg-surface border border-border rounded text-text text-sm" /></div>
          <SelectField label="Slot" value={form.slot} onChange={(v) => set('slot', v)} options={[...SLOT_OPTIONS]} />
          <SelectField label="Item Type" value={form.itemType} onChange={(v) => set('itemType', v)} options={[...ITEM_TYPE_OPTIONS]} />
          <NumberField label="Level" value={form.level} onChange={(v) => set('level', v)} />
          <NumberField label="Weight" value={form.weight} onChange={(v) => set('weight', v)} />
          <div><label className="text-text-muted text-xs block mb-1">Color</label><input type="text" value={form.color} onChange={(e) => set('color', e.target.value)} className="w-full p-2 bg-surface border border-border rounded text-text text-sm" /></div>
          <div><label className="text-text-muted text-xs block mb-1">Owner ID</label><input type="number" value={form.ownerId ?? ''} onChange={(e) => set('ownerId', e.target.value === '' ? null : parseInt(e.target.value) || null)} className="w-full p-2 bg-surface border border-border rounded text-text text-sm" /></div>
          <div className="flex items-center gap-4 col-span-2 pt-1">
            <CheckboxField label="Visible" checked={form.isVisible} onChange={(v) => set('isVisible', v)} />
            <CheckboxField label="Immovable" checked={form.isImmovable} onChange={(v) => set('isImmovable', v)} />
            <CheckboxField label="Equipped" checked={form.isEquipped} onChange={(v) => set('isEquipped', v)} />
          </div>
        </div>
        <div className="pt-4 border-t border-border">
          <CombatFieldsEditor form={form} onChange={(u) => setForm(prev => ({ ...prev, ...u }))} slot={form.slot} />
        </div>
        <div className="flex gap-2 mt-4">
          <Button variant="primary" onClick={handleSave} disabled={updateMutation.isPending}>{updateMutation.isPending ? 'Saving...' : 'Save'}</Button>
          <Button variant="secondary" onClick={onDone}>Cancel</Button>
        </div>
        {updateMutation.isError && <div className="mt-3 text-danger text-sm">Failed to save: {(updateMutation.error as Error)?.message}</div>}
      </div>
    </div>
  )
}