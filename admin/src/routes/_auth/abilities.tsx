import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useAbilities,
  useCreateAbility,
  useDeleteAbility,
  type AbilityInput,
} from '../../hooks/useAbilities'
import { useTags } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { TagInput } from '../../components/TagInput'
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from '../../components/FormFields'
import { showToast } from '../../components/Toast'
import type { Ability } from '../../hooks/useAbilities'

export const Route = createFileRoute('/_auth/abilities')({
  component: AbilitiesManagement,
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

const EMPTY_ABILITY: AbilityInput = {
  name: '',
  description: '',
  ability_type: 'combat',
  requirements: '1',
  cost: 0,
  cooldown: 0,
  cooldown_seconds: 0,
  mana_cost: 0,
  stamina_cost: 0,
  hp_cost: 0,
  proc_chance: 0,
  proc_event: '',
  ability_class: 'active',
  required_tag: '',
}

function CreateAbilityForm({ onSuccess }: { onSuccess: () => void }) {
  const createAbility = useCreateAbility()
  const { data: availableTags } = useTags()
  const [formData, setFormData] = useState<AbilityInput>(EMPTY_ABILITY)

  const selectedTags = formData.required_tag
    ? formData.required_tag.split(',').map((t) => t.trim()).filter(Boolean)
    : []

  const set = (patch: Partial<AbilityInput>) => setFormData((prev) => ({ ...prev, ...patch }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createAbility.mutateAsync(formData)
      showToast('Ability created', 'success')
      setFormData(EMPTY_ABILITY)
      onSuccess()
    } catch {
      // Error is toasted by global onError handler
    }
  }

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">Add New Ability</h3>
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={formData.description} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Ability Type" value={formData.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} />
        <TagInput
          label="Required Tag (optional)"
          value={selectedTags}
          onChange={(tags) => set({ required_tag: tags.join(', ') })}
          availableTags={(availableTags ?? []).map((t) => t.name)}
          placeholder="e.g., sword, fire, healing"
        />
        <div className="grid grid-cols-3 gap-3">
          <FormField label="Level Req" value={formData.requirements} onChange={(v) => set({ requirements: v })} />
          <NumberField label="Cost" value={formData.cost} onChange={(v) => set({ cost: v })} />
          <NumberField label="Cooldown (s)" value={formData.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
        </div>
        <div className="grid grid-cols-3 gap-3">
          <NumberField label="Mana Cost" value={formData.mana_cost} onChange={(v) => set({ mana_cost: v })} />
          <NumberField label="Stamina Cost" value={formData.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
          <NumberField label="HP Cost" value={formData.hp_cost} onChange={(v) => set({ hp_cost: v })} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <NumberField label="Proc Chance (0–1)" value={formData.proc_chance} onChange={(v) => set({ proc_chance: v })} />
          <FormField label="Proc Event" value={formData.proc_event} onChange={(v) => set({ proc_event: v })} placeholder="e.g., on_hit, on_crit" />
        </div>
        <SelectField label="Ability Class" value={formData.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} />
        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={createAbility.isPending} fullWidth>
            {createAbility.isPending ? 'Creating...' : 'Create Ability'}
          </Button>
        </div>
      </form>
    </div>
  )
}

const COLUMNS: Column<Ability>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_, row) => (
      <Link
        to="/abilities/$abilityId"
        params={{ abilityId: String(row.id) }}
        className="no-underline text-primary hover:underline font-bold"
      >
        {row.name}
      </Link>
    ),
  },
  { header: 'Description', accessor: 'description' },
  {
    header: 'Type',
    accessor: 'ability_type',
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: 'Class',
    accessor: 'ability_class',
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: 'Costs',
    accessor: 'mana_cost',
    render: (_, row: Ability) => {
      const parts: React.ReactNode[] = []
      if (row.mana_cost > 0) parts.push(<span key="mp" className="cost-badge" title="Mana Cost">MP: {row.mana_cost}</span>)
      if (row.stamina_cost > 0) parts.push(<span key="sp" className="cost-badge" title="Stamina Cost">SP: {row.stamina_cost}</span>)
      if (row.hp_cost > 0) parts.push(<span key="hp" className="cost-badge" title="HP Cost">HP: {row.hp_cost}</span>)
      if (row.cooldown_seconds > 0) parts.push(<span key="cd" className="cost-badge" title="Cooldown">CD: {row.cooldown_seconds}s</span>)
      return parts.length > 0 ? parts : <span className="text-muted">—</span>
    },
  },
]

function AbilitiesManagement() {
  const [filterType, setFilterType] = useState<string>('')
  const [filterClass, setFilterClass] = useState<string>('')
  const [showCreate, setShowCreate] = useState(false)

  const { data: abilities, isLoading, error } = useAbilities({
    type: filterType || undefined,
    abilityClass: filterClass || undefined,
  })

  if (isLoading) return <div className="loading">Loading abilities...</div>
  if (error) return <div className="error">Failed to load abilities: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Abilities"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => setShowCreate(!showCreate)}>
            {showCreate ? 'Cancel' : '+ Add Ability'}
          </Button>
        }
      />

      <div className="filters-bar">
        <div className="filter-group">
          <label>Type:</label>
          <select value={filterType} onChange={(e) => setFilterType(e.target.value)}>
            <option value="">All Types</option>
            <option value="combat">Combat</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
            <option value="healing">Healing</option>
            <option value="support">Support</option>
            <option value="defensive">Defensive</option>
          </select>
        </div>
        <div className="filter-group">
          <label>Class:</label>
          <select value={filterClass} onChange={(e) => setFilterClass(e.target.value)}>
            <option value="">All Classes</option>
            <option value="active">Active</option>
            <option value="passive">Passive</option>
            <option value="toggle">Toggle</option>
          </select>
        </div>
        {(filterType || filterClass) && (
          <Button variant="ghost" size="sm" onClick={() => { setFilterType(''); setFilterClass('') }}>
            Clear Filters
          </Button>
        )}
      </div>

      {showCreate && <CreateAbilityForm onSuccess={() => setShowCreate(false)} />}

      <DataTable
        columns={COLUMNS}
        data={abilities ?? []}
        getKey={(row) => row.id}
        emptyMessage={
          filterType
            ? 'No abilities match this filter.'
            : 'No abilities found. Create your first ability!'
        }
      />
    </div>
  )
}