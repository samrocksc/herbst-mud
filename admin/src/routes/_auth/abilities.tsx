import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useAbilities, useCreateAbility, useUpdateAbility, useDeleteAbility, type Ability, type AbilityInput } from '../../hooks/useAbilities'
import { PageHeader } from '../../components/PageHeader'

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
          <label>Required Tag (optional):</label>
          <input
            type="text"
            value={formData.required_tag}
            onChange={(e) => setFormData({ ...formData, required_tag: e.target.value })}
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
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : ability ? 'Update Ability' : 'Create Ability'}
          </button>
          <button type="button" className="btn-cancel" onClick={onCancel}>
            Cancel
          </button>
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
          <button className="modal-close" onClick={onCancel}>×</button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{ability.name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <button className="btn-danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </button>
          <button className="btn-cancel" onClick={onCancel}>Cancel</button>
        </div>
      </div>
    </div>
  )
}

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

  if (isLoading) return <div className="loading">Loading abilities...</div>
  if (error) return <div className="error">Failed to load abilities: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Abilities"
        backTo="/dashboard"
        actions={<button className="btn-primary" onClick={() => { setEditingAbility(null); setShowForm(true) }}>+ Add Ability</button>}
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
          <button className="btn-clear-filters" onClick={() => setFilterType('')}>
            Clear Filters
          </button>
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

      <div className="table-container">
        <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Type</th>
              <th>Effects</th>
              <th>Costs</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {abilities?.map((ability) => (
              <tr key={ability.id}>
                <td><strong>{ability.name}</strong></td>
                <td className="text-muted">{ability.description || '-'}</td>
                <td>
                  <span className={`talent-effect talent-effect-${ability.skill_type}`}>
                    {ability.skill_type}
                  </span>
                </td>
                <td>
                  {ability.effect_type && (
                    <span className="talent-effect">{ability.effect_type}</span>
                  )}
                  {ability.effect_value > 0 && (
                    <span className="talent-effect-value"> {ability.effect_value}</span>
                  )}
                  {ability.effect_duration > 0 && (
                    <span className="text-muted"> ({ability.effect_duration}t)</span>
                  )}
                  {ability.proc_chance > 0 && (
                    <span className="text-muted"> proc {Math.round(ability.proc_chance * 100)}%</span>
                  )}
                </td>
                <td>
                  {ability.mana_cost > 0 && <span className="cost-badge" title="Mana Cost">MP: {ability.mana_cost}</span>}
                  {ability.stamina_cost > 0 && <span className="cost-badge" title="Stamina Cost">SP: {ability.stamina_cost}</span>}
                  {ability.hp_cost > 0 && <span className="cost-badge" title="HP Cost">HP: {ability.hp_cost}</span>}
                  {ability.cooldown_seconds > 0 && <span className="cost-badge" title="Cooldown">CD: {ability.cooldown_seconds}s</span>}
                </td>
                <td>
                  <button onClick={() => handleEdit(ability)}>Edit</button>
                  <button className="danger" onClick={() => setDeletingAbility(ability)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {(!abilities || abilities.length === 0) && (
          <div className="empty-state">
            <p>No abilities found. {filterType ? 'Try clearing filters.' : 'Create your first ability!'}</p>
          </div>
        )}
      </div>

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
