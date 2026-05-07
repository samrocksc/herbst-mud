import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useAbilities,
  useCreateAbility,
  useUpdateAbility,
  useDeleteAbility,
  type Ability,
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

export const Route = createFileRoute('/_auth/abilities')({
  component: AbilitiesManagement,
})

const SKILL_TYPE_OPTS = [
  { value: 'combat', label: 'Combat' },
  { value: 'magic', label: 'Magic' },
  { value: 'utility', label: 'Utility' },
  { value: 'healing', label: 'Healing' },
  { value: 'support', label: 'Support' },
]

const EFFECT_TYPE_OPTS = [
  { value: '', label: '— None —' },
  { value: 'damage', label: 'Damage' },
  { value: 'heal', label: 'Heal' },
  { value: 'buff', label: 'Buff' },
  { value: 'debuff', label: 'Debuff' },
  { value: 'dot', label: 'Damage over Time' },
  { value: 'hot', label: 'Heal over Time' },
  { value: 'concentrate', label: 'Concentrate' },
  { value: 'haymaker', label: 'Haymaker' },
  { value: 'scream', label: 'Scream' },
  { value: 'slap', label: 'Slap' },
  { value: 'backoff', label: 'Back-off' },
]

const SCALING_STAT_OPTS = [
  { value: '', label: '— None —' },
  { value: 'STR', label: 'Strength (STR)' },
  { value: 'DEX', label: 'Dexterity (DEX)' },
  { value: 'INT', label: 'Intelligence (INT)' },
  { value: 'WIS', label: 'Wisdom (WIS)' },
  { value: 'CON', label: 'Constitution (CON)' },
]

const SKILL_CLASS_OPTS = [
  { value: 'active', label: 'Active' },
  { value: 'passive', label: 'Passive' },
  { value: 'toggle', label: 'Toggle' },
]

const EMPTY_ABILITY: AbilityInput = {
  name: '',
  description: '',
  skill_type: 'combat',
  requirements: 1,
  cost: 0,
  cooldown: 0,
  cooldown_seconds: 0,
  mana_cost: 0,
  stamina_cost: 0,
  hp_cost: 0,
  effect_type: '',
  effect_value: 0,
  effect_duration: 0,
  scaling_stat: '',
  scaling_percent_per_point: 0,
  proc_chance: 0,
  proc_event: '',
  skill_class: 'active',
  required_tag: '',
}

function AbilityForm({
  ability,
  onSubmit,
  onCancel,
  isLoading,
}: Readonly<{
  ability: Ability | null
  onSubmit: (data: AbilityInput) => void
  onCancel: () => void
  isLoading: boolean
}>) {
  const { tags: availableTags } = useTags()

  const [formData, setFormData] = useState<AbilityInput>(() => {
    if (ability) {
      let reqNum = 1
      if (typeof ability.requirements === 'number') {
        reqNum = ability.requirements
      } else if (typeof ability.requirements === 'string') {
        const n = parseInt(ability.requirements, 10)
        if (!isNaN(n)) reqNum = n
      }
      return {
        ...ability,
        requirements: reqNum,
        scaling_percent_per_point: ability.scaling_percent_per_point ?? 0,
        proc_chance: ability.proc_chance ?? 0,
      } as AbilityInput
    }
    return EMPTY_ABILITY
  })

  const selectedTags = formData.required_tag
    ? formData.required_tag.split(',').map((t) => t.trim()).filter(Boolean)
    : []

  const set = (patch: Partial<AbilityInput>) => setFormData((prev) => ({ ...prev, ...patch }))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {ability ? 'Edit Ability' : 'Add New Ability'}
      </h3>
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} tooltip="The command name players type, e.g. 'concentrate'" />

        <TextareaField
          label="Description"
          value={formData.description}
          onChange={(v) => set({ description: v })}
          rows={3}
          tooltip="Flavor text shown to players when they examine or help this ability"
        />

        <SelectField
          label="Skill Type"
          value={formData.skill_type}
          onChange={(v) => set({ skill_type: v })}
          options={SKILL_TYPE_OPTS}
          tooltip="Category: combat=direct attack, magic=spell, utility=non-combat, healing=restore HP, support=buff allies"
        />

        <TagInput
          label="Required Tag (optional)"
          value={selectedTags}
          onChange={(tags) => set({ required_tag: tags.join(', ') })}
          availableTags={availableTags.map((t) => t.name)}
          placeholder="e.g., sword, fire, healing"
          tooltip="Comma-separated item tags. Character must have an item with this tag equipped to use this ability"
        />

        <div className="grid grid-cols-3 gap-3">
          <NumberField
            label="Level Req"
            value={formData.requirements}
            onChange={(v) => set({ requirements: v })}
            tooltip="Minimum character level to learn this ability"
          />
          <NumberField
            label="Cost"
            value={formData.cost}
            onChange={(v) => set({ cost: v })}
            tooltip="Legacy flat energy cost. Prefer mana_cost/stamina_cost for new abilities"
          />
          <NumberField
            label="Cooldown (s)"
            value={formData.cooldown_seconds}
            onChange={(v) => set({ cooldown_seconds: v })}
            tooltip="Seconds before ability can be reused. Back-off=10s, Haymaker=6s"
          />
        </div>

        <SelectField
          label="Effect Type"
          value={formData.effect_type}
          onChange={(v) => set({ effect_type: v })}
          options={EFFECT_TYPE_OPTS}
          tooltip="What happens on use: damage/heal/buff/debuff, or special: concentrate/haymaker/scream/slap/backoff"
        />

        <div className="grid grid-cols-2 gap-3">
          <NumberField
            label="Effect Value"
            value={formData.effect_value}
            onChange={(v) => set({ effect_value: v })}
            tooltip="Base magnitude. Scales with Scaling Stat if set"
          />
          <NumberField
            label="Effect Duration (ticks)"
            value={formData.effect_duration}
            onChange={(v) => set({ effect_duration: v })}
            tooltip="Combat ticks the effect lasts. 1=one round, 4=Concentrate's full duration"
          />
        </div>

        <div className="grid grid-cols-3 gap-3">
          <NumberField
            label="Mana Cost"
            value={formData.mana_cost}
            onChange={(v) => set({ mana_cost: v })}
            tooltip="MP drained on activation. Mages rely on this pool. If mana < cost, ability fails"
          />
          <NumberField
            label="Stamina Cost"
            value={formData.stamina_cost}
            onChange={(v) => set({ stamina_cost: v })}
            tooltip="SP drained on activation. Fighters use this resource"
          />
          <NumberField
            label="HP Cost"
            value={formData.hp_cost}
            onChange={(v) => set({ hp_cost: v })}
            tooltip="Self-damage to cast. Berserker and blood magic abilities"
          />
        </div>

        <div className="grid grid-cols-2 gap-3">
          <SelectField
            label="Scaling Stat"
            value={formData.scaling_stat}
            onChange={(v) => set({ scaling_stat: v })}
            options={SCALING_STAT_OPTS}
            tooltip="Which character stat boosts the effect. STR=damage, WIS=healing, DEX=dodge"
          />
          <NumberField
            label="Scaling %/point"
            value={formData.scaling_percent_per_point}
            onChange={(v) => set({ scaling_percent_per_point: v })}
            tooltip="Percentage of stat value added per point. 0.05=5%. Formula: final = base + (stat × pct × base)"
          />
        </div>

        <div className="grid grid-cols-2 gap-3">
          <NumberField
            label="Proc Chance (0–1)"
            value={formData.proc_chance}
            onChange={(v) => set({ proc_chance: v })}
            tooltip="0.0–1.0 chance the effect triggers on the proc_event. 0.3=30%"
          />
          <FormField
            label="Proc Event"
            value={formData.proc_event}
            onChange={(v) => set({ proc_event: v })}
            placeholder="e.g., on_hit, on_crit"
            tooltip="When to roll proc_chance: on_hit=attack lands, on_crit=critical, on_dodge=dodging, on_kill=enemy death"
          />
        </div>

        <SelectField
          label="Skill Class"
          value={formData.skill_class}
          onChange={(v) => set({ skill_class: v })}
          options={SKILL_CLASS_OPTS}
          tooltip="active=press button, passive=always on, toggle=on/off switch"
        />

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isLoading} fullWidth>
            {isLoading ? 'Saving...' : ability ? 'Update Ability' : 'Create Ability'}
          </Button>
          <Button variant="secondary" onClick={onCancel} fullWidth>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}

function DeleteConfirmation({
  ability,
  onConfirm,
  onCancel,
  isLoading,
}: Readonly<{
  ability: Ability
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Ability</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">
            ×
          </Button>
        </div>
        <div className="modal-body">
          <p>
            Are you sure you want to delete <strong>{ability.name}</strong>?
          </p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </div>
    </div>
  )
}

const BASE_COLUMNS: Column<Ability>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_: unknown, row: Ability) => <strong>{row.name}</strong>,
  },
  { header: 'Description', accessor: 'description' },
  {
    header: 'Type',
    accessor: 'skill_type',
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: 'Effects',
    accessor: 'effect_type',
    render: (_: unknown, row: Ability) => {
      const parts: React.ReactNode[] = []
      if (row.effect_type) parts.push(<span key="et" className="talent-effect">{row.effect_type}</span>)
      if (row.effect_value > 0) parts.push(<span key="ev" className="talent-effect-value"> {row.effect_value}</span>)
      if (row.effect_duration > 0) parts.push(<span key="ed" className="text-muted"> ({row.effect_duration}t)</span>)
      if (row.proc_chance > 0) parts.push(<span key="pc" className="text-muted"> proc {Math.round(row.proc_chance * 100)}%</span>)
      return parts.length > 0 ? parts : <span className="text-muted">—</span>
    },
  },
  {
    header: 'Costs',
    accessor: 'mana_cost',
    render: (_: unknown, row: Ability) => {
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
  const [showForm, setShowForm] = useState(false)
  const [editingAbility, setEditingAbility] = useState<Ability | null>(null)
  const [deletingAbility, setDeletingAbility] = useState<Ability | null>(null)

  const createAbility = useCreateAbility()
  const updateAbility = useUpdateAbility()
  const deleteAbility = useDeleteAbility()

  const { data: abilities, isLoading, error } = useAbilities({
    type: filterType || undefined,
  })

  const handleSubmit = async (formData: AbilityInput) => {
    if (editingAbility) {
      await updateAbility.mutateAsync({ id: editingAbility.id, input: formData })
    } else {
      await createAbility.mutateAsync(formData)
    }
    setShowForm(false)
    setEditingAbility(null)
  }

  const handleEdit = (ability: Ability) => {
    setEditingAbility(ability)
    setShowForm(true)
  }

  const handleDelete = async () => {
    if (deletingAbility) {
      await deleteAbility.mutateAsync(deletingAbility.id)
      setDeletingAbility(null)
    }
  }

  const handleCancelForm = () => {
    setShowForm(false)
    setEditingAbility(null)
  }

  const columns: Column<Ability>[] = [
    ...BASE_COLUMNS,
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: Ability) => (
        <>
          <Button variant="accent" size="sm" onClick={() => handleEdit(row)}>
            Edit
          </Button>
          <Button
            variant="danger"
            size="sm"
            className="ml-2"
            onClick={() => setDeletingAbility(row)}
          >
            Delete
          </Button>
        </>
      ),
    },
  ]

  if (isLoading) return <div className="loading">Loading abilities...</div>
  if (error) return <div className="error">Failed to load abilities: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Abilities"
        backTo="/dashboard"
        actions={
          <Button
            variant="primary"
            onClick={() => {
              setEditingAbility(null)
              setShowForm(true)
            }}
          >
            + Add Ability
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
          </select>
        </div>
        {filterType && (
          <Button variant="ghost" size="sm" onClick={() => setFilterType('')}>
            Clear Filters
          </Button>
        )}
      </div>

      {showForm && (
        <AbilityForm
          ability={editingAbility}
          onSubmit={handleSubmit}
          onCancel={handleCancelForm}
          isLoading={createAbility.isPending || updateAbility.isPending}
        />
      )}

      <DataTable
        columns={columns}
        data={abilities ?? []}
        getKey={(row) => row.id}
        emptyMessage={
          filterType
            ? 'No abilities match this filter.'
            : 'No abilities found. Create your first ability!'
        }
      />

      {deletingAbility && (
        <DeleteConfirmation
          ability={deletingAbility}
          onConfirm={handleDelete}
          onCancel={() => setDeletingAbility(null)}
          isLoading={deleteAbility.isPending}
        />
      )}
    </div>
  )
}
