import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import type { ReactNode } from 'react'
import { useTalents, useCreateTalent, useUpdateTalent, useDeleteTalent, type Talent, type TalentInput } from '../../hooks/useTalents'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

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
          <Button type="submit" variant="primary" disabled={isLoading}>
            {isLoading ? 'Saving...' : talent ? 'Update Talent' : 'Create Talent'}
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
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{talent.name}</strong>?</p>
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

const BASE_COLUMNS: Column<Talent>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_, row) => <strong>{row.name}</strong>,
  },
  {
    header: 'Description',
    accessor: 'description',
  },
  {
    header: 'Effect',
    accessor: 'effect_type',
    render: (_: unknown, row: Talent) => {
      const parts: ReactNode[] = [
        <span key="et" className={`talent-effect talent-effect-${row.effect_type}`}>{row.effect_type}</span>,
      ]
      if (row.effect_value > 0) {
        parts.push(<span key="ev" className="talent-effect-value"> {row.effect_value}{row.effect_duration > 0 ? ` (${row.effect_duration}t)` : ''}</span>)
      }
      return parts
    },
  },
  {
    header: 'Requirements',
    accessor: 'requirements',
  },
  {
    header: 'Costs',
    accessor: 'mana_cost',
    render: (_: unknown, row: Talent) => {
      const parts: ReactNode[] = []
      parts.push(<span key="mp" className="cost-badge" title="Mana Cost">MP: {row.mana_cost}</span>)
      parts.push(<span key="sp" className="cost-badge" title="Stamina Cost">SP: {row.stamina_cost}</span>)
      if (row.cooldown > 0) {
        parts.push(<span key="cd" className="cost-badge" title="Cooldown">CD: {row.cooldown}</span>)
      }
      return parts
    },
  },
]

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

  const columns: Column<Talent>[] = [
    ...BASE_COLUMNS,
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: Talent) => (
        <>
          <Button variant="accent" size="sm" onClick={() => { setEditingTalent(row); setShowForm(true) }}>Edit</Button>
          <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeletingTalent(row)}>Delete</Button>
        </>
      ),
    },
  ]

  if (isLoading) return <div className="loading">Loading talents...</div>
  if (error) return <div className="error">Failed to load talents: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Talents Management"
        backTo="/dashboard"
        actions={<Button variant="primary" onClick={() => { setEditingTalent(null); setShowForm(true) }}>+ Add Talent</Button>}
      />

      {showForm && (
        <TalentForm
          talent={editingTalent}
          onSubmit={handleSubmit}
          onCancel={handleCancelForm}
          isLoading={createTalent.isPending || updateTalent.isPending}
        />
      )}

      <DataTable
        columns={columns}
        data={talents ?? []}
        getKey={(row: Talent) => row.id}
        emptyMessage="No talents found. Create your first talent!"
      />

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
