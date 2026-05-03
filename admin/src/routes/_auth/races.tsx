import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useRaces, type Race, type RaceInput } from '../../hooks/useRaces'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/races')({
  component: RacesManagement,
})

const EMPTY_FORM: RaceInput = {
  name: '',
  display_name: '',
  description: '',
  stat_modifiers: '',
  is_playable: true,
  color: '',
}

function raceToForm(r: Race): RaceInput {
  return {
    name: r.name,
    display_name: r.display_name,
    description: r.description ?? '',
    stat_modifiers: r.stat_modifiers ? JSON.stringify(r.stat_modifiers, null, 2) : '',
    is_playable: r.is_playable,
    color: r.color ?? '',
  }
}

function RaceForm({
  race,
  onSubmit,
  onCancel,
  isLoading,
}: {
  race: Race | null
  onSubmit: (data: RaceInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [form, setForm] = useState<RaceInput>(() => race ? raceToForm(race) : { ...EMPTY_FORM })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) return
    onSubmit(form)
  }

  const set = <K extends keyof RaceInput>(key: K, value: RaceInput[K]) => {
    setForm(prev => ({ ...prev, [key]: value }))
  }

  return (
    <div className="form-card">
      <h3>{race ? 'Edit Race' : 'Add New Race'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label>Name:</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => set('name', e.target.value)}
            placeholder="e.g. elf, dwarf, human"
            required
          />
        </div>

        <div className="form-row">
          <label>Display Name:</label>
          <input
            type="text"
            value={form.display_name}
            onChange={(e) => set('display_name', e.target.value)}
            placeholder="Defaults to name if blank"
          />
        </div>

        <div className="form-row">
          <label>Description:</label>
          <textarea
            value={form.description}
            onChange={(e) => set('description', e.target.value)}
            rows={3}
          />
        </div>

        <div className="form-row">
          <label>Stat Modifiers (JSON):</label>
          <textarea
            value={form.stat_modifiers}
            onChange={(e) => set('stat_modifiers', e.target.value)}
            rows={4}
            placeholder='e.g. {"str": 2, "dex": -1}'
          />
        </div>

        <div className="form-row">
          <label>Playable:</label>
          <input
            type="checkbox"
            checked={form.is_playable}
            onChange={(e) => set('is_playable', e.target.checked)}
          />
        </div>

        <div className="form-row">
          <label>Color:</label>
          <input
            type="text"
            value={form.color}
            onChange={(e) => set('color', e.target.value)}
            placeholder="e.g. #8b5cf6"
          />
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading || !form.name.trim()}>
            {isLoading ? 'Saving…' : race ? 'Update Race' : 'Create Race'}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}

function DeleteConfirmation({
  race,
  onConfirm,
  onCancel,
  isLoading,
}: {
  race: Race
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Race</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">
            ×
          </Button>
        </div>
        <div className="modal-body">
          <p>
            Are you sure you want to delete <strong>{race.display_name || race.name}</strong>?
          </p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting…' : 'Delete'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </div>
    </div>
  )
}

const COLUMNS: Column<Race>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_: unknown, row: Race) => <strong>{row.display_name || row.name}</strong>,
  },
  {
    header: 'Description',
    accessor: 'description',
  },
  {
    header: 'Playable',
    accessor: 'is_playable',
  },
]

function RacesManagement() {
  const { races, loading, error, createRace, updateRace, deleteRace } = useRaces()
  const [showForm, setShowForm] = useState(false)
  const [editingRace, setEditingRace] = useState<Race | null>(null)
  const [saving, setSaving] = useState(false)
  const [deletingRace, setDeletingRace] = useState<Race | null>(null)
  const [deleting, setDeleting] = useState(false)

  const handleCreate = async (input: RaceInput) => {
    setSaving(true)
    try {
      await createRace(input)
      setShowForm(false)
    } finally {
      setSaving(false)
    }
  }

  const handleUpdate = async (input: RaceInput) => {
    if (!editingRace) return
    setSaving(true)
    try {
      await updateRace(editingRace.id, input)
      setEditingRace(null)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!deletingRace) return
    setDeleting(true)
    try {
      await deleteRace(deletingRace.id)
      setDeletingRace(null)
    } finally {
      setDeleting(false)
    }
  }

  const columns: Column<Race>[] = [
    ...COLUMNS,
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: Race) => (
        <span className="inline-flex gap-2">
          <Button
            variant="accent"
            size="sm"
            onClick={() => { setEditingRace(row); setShowForm(false) }}
          >
            Edit
          </Button>
          <Button
            variant="danger"
            size="sm"
            className="ml-2"
            onClick={() => setDeletingRace(row)}
          >
            Delete
          </Button>
        </span>
      ),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader
        title="Races"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => { setShowForm(true); setEditingRace(null) }}>
            + Add Race
          </Button>
        }
      />

      {error && (
        <div className="error-banner">
          {error}
        </div>
      )}

      {showForm && !editingRace && (
        <RaceForm
          race={null}
          onSubmit={handleCreate}
          onCancel={() => setShowForm(false)}
          isLoading={saving}
        />
      )}

      {editingRace && (
        <RaceForm
          race={editingRace}
          onSubmit={handleUpdate}
          onCancel={() => setEditingRace(null)}
          isLoading={saving}
        />
      )}

      {loading ? (
        <div className="loading">Loading races…</div>
      ) : (
        <DataTable
          columns={columns}
          data={races}
          getKey={(row: Race) => row.id}
          emptyMessage="No races found. Add your first race!"
        />
      )}

      {deletingRace && (
        <DeleteConfirmation
          race={deletingRace}
          onConfirm={handleDelete}
          onCancel={() => setDeletingRace(null)}
          isLoading={deleting}
        />
      )}
    </div>
  )
}