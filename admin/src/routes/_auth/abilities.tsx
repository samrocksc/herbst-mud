import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useAbilities, useCreateAbility, useUpdateAbility, useDeleteAbility, type Ability, type AbilityInput } from '../../hooks/useAbilities'
import { useTags } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { TagInput } from '../../components/TagInput'

export const Route = createFileRoute('/_auth/abilities')({
  component: AbilitiesManagement,
})

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
  isLoading
}: {
  ability: Ability | null
  onSubmit: (data: AbilityInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const { tags: availableTags } = useTags()

  const [formData, setFormData] = useState<AbilityInput>(() => {
    if (ability) {
      return {
        ...ability,
        requirements: ability.requirements ?? 1,
        scaling_percent_per_point: ability.scaling_percent_per_point ?? 0,
        proc_chance: ability.proc_chance ?? 0,
      }
    }
    return EMPTY_ABILITY
  })

  /** Convert between TagInput (string[]) and the form field (string) */
  const selectedTags = formData.required_tag
    ? formData.required_tag.split(',').map((t) => t.trim()).filter(Boolean)
    : []

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="form-card">
      <h3>{ability ? 'Edit Ability' : 'Add New Ability'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label>Name:</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>

        <div className="form-row">
          <label>Description:</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            rows={3}
          />
        </div>

        <div className="form-row">
          <label>Skill Type:</label>
          <select
            value={formData.skill_type}
            onChange={(e) => setFormData({ ...formData, skill_type: e.target.value })}
          >
            <option value="combat">Combat</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
            <option value="healing">Healing</option>
            <option value="support">Support</option>
          </select>
        </div>

        <div className="form-row">
          <TagInput
            label="Required Tag (optional)"
            value={selectedTags}
            onChange={(tags) => setFormData({ ...formData, required_tag: tags.join(', ') })}
            availableTags={availableTags.map((t) => t.name)}
            placeholder="e.g., sword, fire, healing"
          />
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Level Req:</label>
            <input
              type="number"
              min="1"
              max="100"
              value={formData.requirements}
              onChange={(e) => setFormData({ ...formData, requirements: parseInt(e.target.value) || 1 })}
            />
          </div>
          <div className="form-row">
            <label>Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.cost}
              onChange={(e) => setFormData({ ...formData, cost: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Cooldown (s):</label>
            <input
              type="number"
              min="0"
              value={formData.cooldown_seconds}
              onChange={(e) => setFormData({ ...formData, cooldown_seconds: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row">
          <label>Effect Type:</label>
          <select
            value={formData.effect_type}
            onChange={(e) => setFormData({ ...formData, effect_type: e.target.value })}
          >
            <option value="">— None —</option>
            <option value="damage">Damage</option>
            <option value="heal">Heal</option>
            <option value="buff">Buff</option>
            <option value="debuff">Debuff</option>
            <option value="dot">Damage over Time</option>
            <option value="hot">Heal over Time</option>
            <option value="concentrate">Concentrate</option>
            <option value="haymaker">Haymaker</option>
            <option value="scream">Scream</option>
            <option value="slap">Slap</option>
            <option value="backoff">Back-off</option>
          </select>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Effect Value:</label>
            <input
              type="number"
              min="0"
              value={formData.effect_value}
              onChange={(e) => setFormData({ ...formData, effect_value: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Effect Duration (ticks):</label>
            <input
              type="number"
              min="0"
              value={formData.effect_duration}
              onChange={(e) => setFormData({ ...formData, effect_duration: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Mana Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.mana_cost}
              onChange={(e) => setFormData({ ...formData, mana_cost: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Stamina Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.stamina_cost}
              onChange={(e) => setFormData({ ...formData, stamina_cost: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>HP Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.hp_cost}
              onChange={(e) => setFormData({ ...formData, hp_cost: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Scaling Stat:</label>
            <select
              value={formData.scaling_stat}
              onChange={(e) => setFormData({ ...formData, scaling_stat: e.target.value })}
            >
              <option value="">— None —</option>
              <option value="STR">Strength (STR)</option>
              <option value="DEX">Dexterity (DEX)</option>
              <option value="INT">Intelligence (INT)</option>
              <option value="WIS">Wisdom (WIS)</option>
              <option value="CON">Constitution (CON)</option>
            </select>
          </div>
          <div className="form-row">
            <label>Scaling %/point:</label>
            <input
              type="number"
              min="0"
              step="0.1"
              value={formData.scaling_percent_per_point}
              onChange={(e) => setFormData({ ...formData, scaling_percent_per_point: parseFloat(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Proc Chance (0–1):</label>
            <input
              type="number"
              min="0"
              max="1"
              step="0.01"
              value={formData.proc_chance}
              onChange={(e) => setFormData({ ...formData, proc_chance: parseFloat(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Proc Event:</label>
            <input
              type="text"
              value={formData.proc_event}
              onChange={(e) => setFormData({ ...formData, proc_event: e.target.value })}
              placeholder="e.g., on_hit, on_crit"
            />
          </div>
        </div>

        <div className="form-row">
          <label>Skill Class:</label>
          <select
            value={formData.skill_class}
            onChange={(e) => setFormData({ ...formData, skill_class: e.target.value })}
          >
            <option value="active">Active</option>
            <option value="passive">Passive</option>
            <option value="toggle">Toggle</option>
          </select>
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading}>
            {isLoading ? 'Saving...' : ability ? 'Update Ability' : 'Create Ability'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
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
  isLoading
}: {
  ability: Ability
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Ability</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{ability.name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  )
}

// ─── Table column definitions ─────────────────────────────────────────────────

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
    render: (val: unknown) => <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>,
  },
  {
    header: 'Effects',
    accessor: 'effect_type',
    render: (_: unknown, row: Ability) => {
      const parts: React.ReactNode[] = []
      if (row.effect_type) {
        parts.push(<span key="et" className="talent-effect">{row.effect_type}</span>)
      }
      if (row.effect_value > 0) {
        parts.push(<span key="ev" className="talent-effect-value"> {row.effect_value}</span>)
      }
      if (row.effect_duration > 0) {
        parts.push(<span key="ed" className="text-muted"> ({row.effect_duration}t)</span>)
      }
      if (row.proc_chance > 0) {
        parts.push(<span key="pc" className="text-muted"> proc {Math.round(row.proc_chance * 100)}%</span>)
      }
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
    type: filterType || undefined
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
          <Button variant="accent" size="sm" onClick={() => handleEdit(row)}>Edit</Button>
          <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeletingAbility(row)}>Delete</Button>
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
        actions={<Button variant="primary" onClick={() => { setEditingAbility(null); setShowForm(true) }}>+ Add Ability</Button>}
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
        emptyMessage={filterType ? 'No abilities match this filter.' : 'No abilities found. Create your first ability!'}
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
