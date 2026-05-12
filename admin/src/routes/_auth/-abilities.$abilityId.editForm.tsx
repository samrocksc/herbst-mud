import { useState } from 'react'
import {
  useUpdateAbility,
  useDeleteAbility,
  type Ability,
  type AbilityInput,
} from '../../hooks/useAbilities'
import { useTags } from '../../hooks/useTags'
import { useNavigate } from '@tanstack/react-router'
import { Button } from '../../components/Button'
import { TagInput } from '../../components/TagInput'
import { EffectsSubForm } from '../../components/EffectsSubForm'
import { DeleteConfirmation } from '../../components/DeleteConfirmation'
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from '../../components/FormFields'
import { showToast } from '../../components/Toast'

const ABILITY_TYPE_OPTS = [
  { value: 'combat', label: 'Combat' },
  { value: 'magic', label: 'Magic' },
  { value: 'utility', label: 'Utility' },
  { value: 'healing', label: 'Healing' },
  { value: 'support', label: 'Support' },
  { value: 'defensive', label: 'Defensive' },
]

const ABILITY_CLASS_OPTS = [
  { value: 'active', label: 'Active' },
  { value: 'passive', label: 'Passive' },
  { value: 'toggle', label: 'Toggle' },
]

export function AbilityEditForm({
  ability,
  abilityId,
  onDone,
}: Readonly<{
  ability: Ability
  abilityId: number
  onDone: () => void
}>) {
  const navigate = useNavigate()
  const updateAbility = useUpdateAbility()
  const deleteAbility = useDeleteAbility()
  const { data: availableTags } = useTags()
  const [showDeleteModal, setShowDeleteModal] = useState(false)

  const [formData, setFormData] = useState<AbilityInput>({
    name: ability.name,
    description: ability.description,
    ability_type: ability.ability_type,
    requirements: ability.requirements ?? '',
    cost: ability.cost,
    cooldown: ability.cooldown,
    cooldown_seconds: ability.cooldown_seconds,
    mana_cost: ability.mana_cost,
    stamina_cost: ability.stamina_cost,
    hp_cost: ability.hp_cost,
    proc_chance: ability.proc_chance,
    proc_event: ability.proc_event ?? '',
    ability_class: ability.ability_class,
    required_tag: ability.required_tag ?? '',
  })

  const selectedTags = formData.required_tag
    ? formData.required_tag.split(',').map((t) => t.trim()).filter(Boolean)
    : []

  const set = (patch: Partial<AbilityInput>) => setFormData((prev) => ({ ...prev, ...patch }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await updateAbility.mutateAsync({ id: abilityId, input: formData })
      showToast('Ability updated', 'success')
      onDone()
    } catch {
      // Error is toasted by global onError handler
    }
  }

  const handleDelete = async () => {
    try {
      await deleteAbility.mutateAsync(abilityId)
      showToast('Ability deleted', 'success')
      navigate({ to: '/abilities' })
    } catch {
      // Error is toasted by global onError handler
    }
  }

  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Ability</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} />
          <SelectField label="Ability Type" value={formData.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} />
          <TagInput
            label="Required Tag (optional)"
            value={selectedTags}
            onChange={(tags) => set({ required_tag: tags.join(', ') })}
            availableTags={(availableTags ?? []).map((t) => t.name)}
            placeholder="e.g., sword, fire, healing"
          />
          <SelectField label="Ability Class" value={formData.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} />
        </div>
        <TextareaField label="Description" value={formData.description} onChange={(v) => set({ description: v })} rows={3} />

        <div className="grid grid-cols-3 gap-4">
          <FormField label="Level Req" value={formData.requirements} onChange={(v) => set({ requirements: v })} />
          <NumberField label="Cost" value={formData.cost} onChange={(v) => set({ cost: v })} />
          <NumberField label="Cooldown (s)" value={formData.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
        </div>
        <div className="grid grid-cols-3 gap-4">
          <NumberField label="Mana Cost" value={formData.mana_cost} onChange={(v) => set({ mana_cost: v })} />
          <NumberField label="Stamina Cost" value={formData.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
          <NumberField label="HP Cost" value={formData.hp_cost} onChange={(v) => set({ hp_cost: v })} />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <NumberField label="Proc Chance (0–1)" value={formData.proc_chance} onChange={(v) => set({ proc_chance: v })} step={0.01} />
          <FormField label="Proc Event" value={formData.proc_event} onChange={(v) => set({ proc_event: v })} placeholder="e.g., on_hit, on_crit" />
        </div>

        <EffectsSubForm abilityId={abilityId} />

        <div className="flex gap-2 pt-2">
          <Button type="submit" variant="primary" disabled={updateAbility.isPending}>
            {updateAbility.isPending ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button variant="secondary" onClick={onDone} type="button">Cancel</Button>
          <Button variant="danger" onClick={() => setShowDeleteModal(true)} type="button">Delete Ability</Button>
        </div>
      </form>

      <DeleteConfirmation
        open={showDeleteModal}
        title="Delete Ability"
        message="Are you sure you want to delete this ability? This action cannot be undone."
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteAbility.isPending}
      />
    </div>
  )
}
