import { useState } from 'react'
import { Button } from '../../components/Button'
import { DeleteConfirmation } from '../../components/DeleteConfirmation'
import { FormError } from '../../components/fields/FormError'
import { showToast } from '../../components/Toast'
import { apiPut, apiDelete } from '../../utils/apiFetch'
import { FactionFormFields } from './FactionFormFields'
import { factionToForm, type Faction, type FactionForm, type FactionCategory } from './factionTypes'

export function FactionDetail({
  faction,
  categories,
  onRefresh,
}: Readonly<{
  faction: Faction
  categories: FactionCategory[]
  onRefresh: () => void
}>) {
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState<FactionForm>(factionToForm(faction))
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState(false)
  const [error, setError] = useState('')

  const handleUpdate = async () => {
    setSaving(true)
    setError('')
    try {
      await apiPut(`/api/factions/${faction.id}`, {
        ...form,
        category_id: form.category_id || null,
      })
      showToast('Faction updated', 'success')
      setEditing(false)
      onRefresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Update failed')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    await apiDelete(`/api/factions/${faction.id}`)
    showToast('Faction deleted', 'success')
    onRefresh()
  }

  if (editing) {
    return (
      <div>
        <h2 className="mt-0 mb-4 text-text">Edit Faction</h2>
        {error && <FormError message={error} />}
        <FactionFormFields form={form} setForm={setForm} categories={categories} />
        <div className="flex gap-2 mt-3">
          <Button variant="primary" size="md" fullWidth onClick={handleUpdate} disabled={saving}>
            {saving ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button variant="secondary" size="md" fullWidth onClick={() => setEditing(false)}>
            Cancel
          </Button>
        </div>
      </div>
    )
  }

  const memberTags = faction.member_tags ?? []

  return (
    <div>
      <h2 className="mt-0 mb-4 text-text">{faction.display_name || faction.name}</h2>
      <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-2">
        <DetailRow label="ID" value={String(faction.id)} />
        <DetailRow label="Name" value={faction.name} />
        <DetailRow label="Display Name" value={faction.display_name || '—'} />
        <DetailRow label="Description" value={faction.description || '—'} />
        <DetailRow label="Standing" value={String(faction.standing ?? 0)} />
        <DetailRow label="Universal" value={faction.is_universal ? 'Yes' : 'No'} />
        <DetailRow label="Members" value={faction.member_count != null ? String(faction.member_count) : '0'} />
        {memberTags.length > 0 && (
          <div className="detail-row">
            <label>Member Tags</label>
            <div className="flex flex-wrap gap-1">
              {memberTags.map((tag) => (
                <span key={tag} className="px-2 py-0.5 bg-primary/20 text-primary text-xs rounded-full">{tag}</span>
              ))}
            </div>
          </div>
        )}
        <div className="flex gap-2 mt-3">
          <Button variant="primary" size="md" fullWidth onClick={() => setEditing(true)}>Edit</Button>
          <Button variant="danger" size="md" fullWidth onClick={() => setConfirmDelete(true)}>Delete</Button>
        </div>
      </div>
      <DeleteConfirmation
        open={confirmDelete}
        title="Delete Faction"
        message={`Are you sure you want to delete "${faction.display_name || faction.name}"? This cannot be undone.`}
        onConfirm={handleDelete}
        onCancel={() => setConfirmDelete(false)}
      />
    </div>
  )
}

function DetailRow({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div className="detail-row">
      <label>{label}</label>
      <span>{value}</span>
    </div>
  )
}