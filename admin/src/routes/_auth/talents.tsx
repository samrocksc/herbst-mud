import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTalents, useCreateTalent, useUpdateTalent, useDeleteTalent, type Talent, type TalentInput } from '../../hooks/useTalents'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/_auth/talents')({
  component: TalentsManagement,
})

const EMPTY_TALENT: TalentInput = {
  name: '',
  description: '',
  requirements: '',
  effect_type: 'heal',
  effect_value: 0,
  effect_duration: 0,
  cooldown: 0,
  mana_cost: 0,
  stamina_cost: 0,
}

const EFFECT_TYPES = ['heal', 'damage', 'dot', 'buff_armor', 'buff_dodge', 'buff_crit', 'debuff']

function TalentForm({
  talent,
  onSubmit,
  onCancel,
  isLoading
}: {
  talent: Talent | null
  onSubmit: (data: TalentInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [formData, setFormData] = useState<TalentInput>(() => {
    if (talent) {
      return {
        id: talent.id,
        name: talent.name,
        description: talent.description,
        requirements: talent.requirements,
        effect_type: talent.effect_type,
        effect_value: talent.effect_value,
        effect_duration: talent.effect_duration,
        cooldown: talent.cooldown,
        mana_cost: talent.mana_cost,
        stamina_cost: talent.stamina_cost,
      }
    }
    return EMPTY_TALENT
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="form-card">
      <h3>{talent ? 'Edit Talent' : 'Add New Talent'}</h3>
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
          <label>Requirements (comma-separated skill IDs or level requirements):</label>
          <input
            type="text"
            value={formData.requirements}
            onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
            placeholder="e.g., skill:1, level:5"
          />
        </div>

        <div className="form-row">
          <label>Effect Type:</label>
          <select
            value={formData.effect_type}
            onChange={(e) => setFormData({ ...formData, effect_type: e.target.value })}
          >
            {EFFECT_TYPES.map(type => (
              <option key={type} value={type}>{type}</option>
            ))}
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
            <label>Cooldown (ticks):</label>
            <input
              type="number"
              min="0"
              value={formData.cooldown}
              onChange={(e) => setFormData({ ...formData, cooldown: parseInt(e.target.value) || 0 })}
            />
          </div>
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
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : talent ? 'Update Talent' : 'Create Talent'}
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
  talent,
  onConfirm,
  onCancel,
  isLoading
}: {
  talent: Talent
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Talent</h3>
          <button className="modal-close" onClick={onCancel}>×</button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{talent.name}</strong>?</p>
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

function TalentsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editingTalent, setEditingTalent] = useState<Talent | null>(null)
  const [deletingTalent, setDeletingTalent] = useState<Talent | null>(null)

  const createTalent = useCreateTalent()
  const updateTalent = useUpdateTalent()
  const deleteTalent = useDeleteTalent()

  const { data: talents, isLoading, error } = useTalents()

  const handleSubmit = async (formData: TalentInput) => {
    if (editingTalent) {
      await updateTalent.mutateAsync({ id: editingTalent.id, input: formData })
    } else {
      await createTalent.mutateAsync(formData)
    }
    setShowForm(false)
    setEditingTalent(null)
  }

  const handleEdit = (talent: Talent) => {
    setEditingTalent(talent)
    setShowForm(true)
  }

  const handleDelete = async () => {
    if (deletingTalent) {
      await deleteTalent.mutateAsync(deletingTalent.id)
      setDeletingTalent(null)
    }
  }

  const handleCancelForm = () => {
    setShowForm(false)
    setEditingTalent(null)
  }

  if (isLoading) return <div className="loading">Loading talents...</div>
  if (error) return <div className="error">Failed to load talents: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Talents Management"
        backTo="/dashboard"
        actions={<button className="btn-primary" onClick={() => { setEditingTalent(null); setShowForm(true) }}>+ Add Talent</button>}
      />

      {showForm && (
        <TalentForm
          talent={editingTalent}
          onSubmit={handleSubmit}
          onCancel={handleCancelForm}
          isLoading={createTalent.isPending || updateTalent.isPending}
        />
      )}

      <div className="table-container">
        <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Effect</th>
              <th>Requirements</th>
              <th>Costs</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {talents?.map((talent) => (
              <tr key={talent.id}>
                <td>
                  <strong>{talent.name}</strong>
                </td>
                <td className="text-muted">{talent.description || '-'}</td>
                <td>
                  <span className={`talent-effect talent-effect-${talent.effect_type}`}>
                    {talent.effect_type}
                  </span>
                  {talent.effect_value > 0 && (
                    <span className="talent-effect-value">
                      {talent.effect_value}
                      {talent.effect_duration > 0 && ` (${talent.effect_duration}t)`}
                    </span>
                  )}
                </td>
                <td className="text-muted">{talent.requirements || '-'}</td>
                <td>
                  <span className="cost-badge" title="Mana Cost">MP: {talent.mana_cost}</span>
                  <span className="cost-badge" title="Stamina Cost">SP: {talent.stamina_cost}</span>
                  {talent.cooldown > 0 && (
                    <span className="cost-badge" title="Cooldown">CD: {talent.cooldown}</span>
                  )}
                </td>
                <td>
                  <button onClick={() => handleEdit(talent)}>Edit</button>
                  <button className="danger" onClick={() => setDeletingTalent(talent)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {(!talents || talents.length === 0) && (
          <div className="empty-state">
            <p>No talents found. Create your first talent!</p>
          </div>
        )}
      </div>

      {deletingTalent && (
        <DeleteConfirmation
          talent={deletingTalent}
          onConfirm={handleDelete}
          onCancel={() => setDeletingTalent(null)}
          isLoading={deleteTalent.isPending}
        />
      )}
    </div>
  )
}
