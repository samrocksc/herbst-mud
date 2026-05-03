import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useAchievements,
  useCreateAchievement,
  useUpdateAchievement,
  useDeleteAchievement,
  type Achievement,
  type AchievementInput,
} from '../../hooks/useAchievements'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/achievements')({
  component: AchievementsManagement,
})

const EMPTY_INPUT: AchievementInput = {
  name: '',
  description: '',
  icon: '',
  xp_reward: 0,
  criteria: '',
}

function AchievementForm({
  achievement,
  onSubmit,
  onCancel,
  isLoading,
}: {
  achievement: Achievement | null
  onSubmit: (data: AchievementInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [formData, setFormData] = useState<AchievementInput>(() => {
    if (achievement) {
      return {
        name: achievement.name,
        description: achievement.description,
        icon: achievement.icon,
        xp_reward: achievement.xp_reward,
        criteria: achievement.criteria,
      }
    }
    return EMPTY_INPUT
  })

  return (
    <div className="form-card">
      <h3>{achievement ? 'Edit Achievement' : 'Add New Achievement'}</h3>
      <form
        onSubmit={(e) => {
          e.preventDefault()
          onSubmit(formData)
        }}
      >
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
          <label>Icon (emoji):</label>
          <input
            type="text"
            value={formData.icon}
            onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
            placeholder="e.g. 🏆"
          />
        </div>

        <div className="form-row">
          <label>XP Reward:</label>
          <input
            type="number"
            value={formData.xp_reward}
            onChange={(e) => setFormData({ ...formData, xp_reward: Number(e.target.value) })}
            min={0}
          />
        </div>

        <div className="form-row">
          <label>Criteria (JSON):</label>
          <textarea
            value={formData.criteria}
            onChange={(e) => setFormData({ ...formData, criteria: e.target.value })}
            rows={3}
            placeholder='e.g. {"type":"kill_count","target":10}'
          />
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading}>
            {isLoading ? 'Saving...' : achievement ? 'Update Achievement' : 'Create Achievement'}
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
  achievement,
  onConfirm,
  onCancel,
  isLoading,
}: {
  achievement: Achievement
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Achievement</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">
            ×
          </Button>
        </div>
        <div className="modal-body">
          <p>
            Are you sure you want to delete <strong>{achievement.name}</strong>?
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

const COLUMNS: Column<Achievement>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_: unknown, row: Achievement) => <strong>{row.icon ? `${row.icon} ` : ''}{row.name}</strong>,
  },
  {
    header: 'Description',
    accessor: 'description',
  },
  {
    header: 'XP',
    accessor: 'xp_reward',
    render: (val: unknown) =>
      val ? <span className="talent-effect">{String(val)} XP</span> : <span className="text-muted">—</span>,
  },
  {
    header: 'Criteria',
    accessor: 'criteria',
    render: (val: unknown) =>
      val ? <span className="text-sm">{String(val)}</span> : <span className="text-muted">—</span>,
  },
  {
    header: 'Actions',
    accessor: '_actions',
    render: () => null,
  },
]

function AchievementsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Achievement | null>(null)
  const [deleting, setDeleting] = useState<Achievement | null>(null)

  const create = useCreateAchievement()
  const update = useUpdateAchievement()
  const remove = useDeleteAchievement()
  const { data: achievements, isLoading, error } = useAchievements()

  const handleSubmit = async (formData: AchievementInput) => {
    if (editing) {
      await update.mutateAsync({ id: editing.id, input: formData })
    } else {
      await create.mutateAsync(formData)
    }
    setShowForm(false)
    setEditing(null)
  }

  if (isLoading) return <div className="loading">Loading achievements...</div>
  if (error) return <div className="error">Failed to load achievements: {error.message}</div>

  const columns: Column<Achievement>[] = [
    ...COLUMNS.slice(0, 4),
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: Achievement) => (
        <span className="inline-flex gap-2">
          <Button variant="accent" size="sm" onClick={() => { setEditing(row); setShowForm(true) }}>
            Edit
          </Button>
          <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeleting(row)}>
            Delete
          </Button>
        </span>
      ),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader
        title="Achievements"
        backTo="/dashboard"
        actions={
          <Button
            variant="primary"
            onClick={() => {
              setEditing(null)
              setShowForm(true)
            }}
          >
            + Add Achievement
          </Button>
        }
      />

      {showForm && (
        <AchievementForm
          achievement={editing}
          onSubmit={handleSubmit}
          onCancel={() => {
            setShowForm(false)
            setEditing(null)
          }}
          isLoading={create.isPending || update.isPending}
        />
      )}

      <DataTable
        columns={columns}
        data={achievements ?? []}
        getKey={(row: Achievement) => row.id}
        emptyMessage="No achievements found. Add your first achievement!"
      />

      {deleting && (
        <DeleteConfirmation
          achievement={deleting}
          onConfirm={async () => {
            await remove.mutateAsync(deleting.id)
            setDeleting(null)
          }}
          onCancel={() => setDeleting(null)}
          isLoading={remove.isPending}
        />
      )}
    </div>
  )
}