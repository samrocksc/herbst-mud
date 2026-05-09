import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useAbility,
  useUpdateAbility,
  useDeleteAbility,
  type AbilityInput,
} from '../../hooks/useAbilities'
import { useTags } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'
import { TagInput } from '../../components/TagInput'
import { EffectsSubForm } from '../../components/EffectsSubForm'
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from '../../components/FormFields'
import { showToast } from '../../components/Toast'

export const Route = createFileRoute('/_auth/abilities/$abilityId')({
  component: AbilityDetailPage,
})

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

function AbilityDetailPage() {
  const abilityId = Route.useParams().abilityId
  const navigate = useNavigate()
  const { data: ability, isLoading, error } = useAbility(Number(abilityId))
  const updateAbility = useUpdateAbility()
  const deleteAbility = useDeleteAbility()
  const { data: availableTags } = useTags()

  const [formData, setFormData] = useState<AbilityInput | null>(null)
  const [confirmDelete, setConfirmDelete] = useState(false)

  if (isLoading) return <div className="loading">Loading ability...</div>
  if (error) return <div className="error">Failed to load ability: {error.message}</div>
  if (!ability) return <div className="error">Ability not found</div>

  const current = formData ?? {
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
  }

  const selectedTags = current.required_tag
    ? current.required_tag.split(',').map((t) => t.trim()).filter(Boolean)
    : []

  const set = (patch: Partial<AbilityInput>) => setFormData({ ...current, ...patch })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await updateAbility.mutateAsync({ id: Number(abilityId), input: current })
      showToast('Ability updated', 'success')
      setFormData(null)
    } catch {
      // Error is toasted by global onError handler
    }
  }

  const handleDelete = async () => {
    try {
      await deleteAbility.mutateAsync(Number(abilityId))
      showToast('Ability deleted', 'success')
      navigate({ to: '/abilities' })
    } catch {
      // Error is toasted by global onError handler
    }
  }

  return (
    <div className="management-page">
      <PageHeader title={ability.name} backTo="/abilities" />
      <form onSubmit={handleSubmit} className="form-card space-y-3">
        <FormField label="Name" value={current.name} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={current.description} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Ability Type" value={current.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} />
        <TagInput
          label="Required Tag (optional)"
          value={selectedTags}
          onChange={(tags) => set({ required_tag: tags.join(', ') })}
          availableTags={(availableTags ?? []).map((t) => t.name)}
          placeholder="e.g., sword, fire, healing"
        />
        <div className="grid grid-cols-3 gap-3">
          <FormField label="Level Req" value={current.requirements} onChange={(v) => set({ requirements: v })} />
          <NumberField label="Cost" value={current.cost} onChange={(v) => set({ cost: v })} />
          <NumberField label="Cooldown (s)" value={current.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
        </div>
        <div className="grid grid-cols-3 gap-3">
          <NumberField label="Mana Cost" value={current.mana_cost} onChange={(v) => set({ mana_cost: v })} />
          <NumberField label="Stamina Cost" value={current.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
          <NumberField label="HP Cost" value={current.hp_cost} onChange={(v) => set({ hp_cost: v })} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <NumberField label="Proc Chance (0–1)" value={current.proc_chance} onChange={(v) => set({ proc_chance: v })} />
          <FormField label="Proc Event" value={current.proc_event} onChange={(v) => set({ proc_event: v })} placeholder="e.g., on_hit, on_crit" />
        </div>
        <SelectField label="Ability Class" value={current.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} />

        <EffectsSubForm abilityId={Number(abilityId)} />

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={updateAbility.isPending}>
            {updateAbility.isPending ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button variant="danger" onClick={() => setConfirmDelete(true)}>Delete Ability</Button>
        </div>
      </form>

      {confirmDelete && (
        <div className="modal-overlay" onClick={() => setConfirmDelete(false)}>
          <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Delete Ability</h3>
              <Button variant="ghost" size="sm" onClick={() => setConfirmDelete(false)} aria-label="Close">×</Button>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to delete <strong>{ability.name}</strong>?</p>
              <p className="text-muted">This action cannot be undone.</p>
            </div>
            <div className="modal-footer">
              <Button variant="danger" onClick={handleDelete} disabled={deleteAbility.isPending}>
                {deleteAbility.isPending ? 'Deleting...' : 'Delete'}
              </Button>
              <Button variant="secondary" onClick={() => setConfirmDelete(false)}>Cancel</Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}