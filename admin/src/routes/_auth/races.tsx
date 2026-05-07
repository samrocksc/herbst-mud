import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useRaces, type Race, type RaceInput } from '../../hooks/useRaces'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { RaceForm } from '../../components/RaceForm'

export const Route = createFileRoute('/_auth/races')({
  component: RacesManagement,
})

const COLUMNS: Column<Race>[] = [
  { header: 'Name', accessor: 'name', render: (_: unknown, row: Race) => <strong>{row.display_name || row.name}</strong> },
  { header: 'Description', accessor: 'description' },
  { header: 'Slots', accessor: 'equipment_slots', render: (_: unknown, row: Race) => <span className="text-xs">{(row.equipment_slots ?? []).join(', ') || '—'}</span> },
  { header: 'Playable', accessor: 'is_playable' },
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
    try { await createRace(input); setShowForm(false) } finally { setSaving(false) }
  }

  const handleUpdate = async (input: RaceInput) => {
    if (!editingRace) return
    setSaving(true)
    try { await updateRace(editingRace.id, input); setEditingRace(null) } finally { setSaving(false) }
  }

  const handleDelete = async () => {
    if (!deletingRace) return
    setDeleting(true)
    try { await deleteRace(deletingRace.id); setDeletingRace(null) } finally { setDeleting(false) }
  }

  const columns: Column<Race>[] = [
    ...COLUMNS,
    { header: 'Actions', accessor: '_actions', render: (_: unknown, row: Race) => (
      <span className="inline-flex gap-2">
        <Button variant="accent" size="sm" onClick={() => { setEditingRace(row); setShowForm(false) }}>Edit</Button>
        <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeletingRace(row)}>Delete</Button>
      </span>
    )},
  ]

  return (
    <div className="management-page">
      <PageHeader title="Races" backTo="/dashboard" actions={<Button variant="primary" onClick={() => { setShowForm(true); setEditingRace(null) }}>+ Add Race</Button>} />
      {error && <div className="error-banner">{error}</div>}
      {showForm && !editingRace && <RaceForm race={null} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={saving} />}
      {editingRace && <RaceForm race={editingRace} onSubmit={handleUpdate} onCancel={() => setEditingRace(null)} isLoading={saving} />}
      {loading ? <div className="loading">Loading races...</div> : (
        <DataTable columns={columns} data={races} getKey={(row: Race) => row.id} emptyMessage="No races found. Add your first race!" />
      )}
      {deletingRace && <DeleteConfirmation race={deletingRace} onConfirm={handleDelete} onCancel={() => setDeletingRace(null)} isLoading={deleting} />}
    </div>
  )
}

function DeleteConfirmation({ race, onConfirm, onCancel, isLoading }: Readonly<{ race: Race; onConfirm: () => void; onCancel: () => void; isLoading: boolean }>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header"><h3>Delete Race</h3><Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button></div>
        <div className="modal-body"><p>Are you sure you want to delete <strong>{race.display_name || race.name}</strong>?</p><p className="text-muted">This action cannot be undone.</p></div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>{isLoading ? 'Deleting...' : 'Delete'}</Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  )
}