import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTags, type Tag, type TagInput } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/tags')({
  component: TagsManagement,
})

const DEFAULT_COLOR = '#8b5cf6'

function ColorDot({ color }: { color: string }) {
  return (
    <span
      className="inline-block w-3 h-3 rounded-full shrink-0"
      style={{ backgroundColor: color || DEFAULT_COLOR }}
    />
  )
}

function TagForm({
  tag,
  onSubmit,
  onCancel,
  isLoading,
}: {
  tag: Tag | null
  onSubmit: (data: TagInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [form, setForm] = useState<TagInput>(() => ({
    name: tag?.name ?? '',
    color: tag?.color ?? DEFAULT_COLOR,
  }))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) return
    onSubmit({ name: form.name.trim(), color: form.color })
  }

  return (
    <div className="form-card">
      <h3>{tag ? 'Edit Tag' : 'Add New Tag'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label>Name:</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="e.g. fire, magic, warrior"
            required
          />
        </div>

        <div className="form-row">
          <label>Color:</label>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <input
              type="color"
              value={form.color || DEFAULT_COLOR}
              onChange={(e) => setForm({ ...form, color: e.target.value })}
              style={{ width: '2.5rem', height: '2rem', padding: '2px', cursor: 'pointer' }}
            />
            <input
              type="text"
              value={form.color}
              onChange={(e) => setForm({ ...form, color: e.target.value })}
              placeholder="#8b5cf6"
              pattern="^#[0-9a-fA-F]{6}$"
              style={{ width: '7rem' }}
            />
            <ColorDot color={form.color} />
          </div>
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading || !form.name.trim()}>
            {isLoading ? 'Saving…' : tag ? 'Update Tag' : 'Create Tag'}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}

function TagsManagement() {
  const { tags, error, createTag, updateTag, deleteTag } = useTags()
  const [showForm, setShowForm] = useState(false)
  const [editingTag, setEditingTag] = useState<Tag | null>(null)
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [deleting, setDeleting] = useState(false)

  const handleCreate = async (input: TagInput) => {
    setSaving(true)
    try {
      await createTag(input)
      setShowForm(false)
    } finally {
      setSaving(false)
    }
  }

  const handleUpdate = async (input: TagInput) => {
    if (!editingTag) return
    setSaving(true)
    try {
      await updateTag(editingTag.id, input)
      setEditingTag(null)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    setDeleting(true)
    try {
      await deleteTag(id)
      setConfirmDelete(null)
    } finally {
      setDeleting(false)
    }
  }

  const columns: Column<Tag>[] = [
    {
      header: 'Name',
      accessor: 'name',
      render: (_, row) => (
        <span style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <ColorDot color={row.color} />
          {row.name}
        </span>
      ),
    },
    {
      header: 'Color',
      accessor: 'color',
      render: (val: unknown) => val ? (
        <code style={{ fontSize: '0.75rem', color: 'var(--color-accent)' }}>{String(val)}</code>
      ) : <span className="text-muted">—</span>,
    },
    {
      header: '',
      accessor: 'id',
      align: 'right',
      render: (_, row) => (
        <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => { e.stopPropagation(); setEditingTag(row); setShowForm(false) }}
            aria-label={`Edit ${row.name}`}
          >
            Edit
          </Button>
          <Button
            variant="danger"
            size="sm"
            onClick={(e) => { e.stopPropagation(); setConfirmDelete(row.id) }}
            aria-label={`Delete ${row.name}`}
          >
            Delete
          </Button>
        </div>
      ),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader
        title="Tags"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => { setShowForm(true); setEditingTag(null) }}>
            + Add Tag
          </Button>
        }
      />

      {error && (
        <div className="error-banner">
          {error}
        </div>
      )}

      {showForm && !editingTag && (
        <TagForm
          tag={null}
          onSubmit={handleCreate}
          onCancel={() => setShowForm(false)}
          isLoading={saving}
        />
      )}

      {editingTag && (
        <TagForm
          tag={editingTag}
          onSubmit={handleUpdate}
          onCancel={() => setEditingTag(null)}
          isLoading={saving}
        />
      )}

      <DataTable
        columns={columns}
        data={tags}
        getKey={(row) => row.id}
        emptyMessage="No tags yet. Create one above."
      />

      {confirmDelete !== null && (
        <div className="modal-overlay">
          <div className="modal-card">
            <h3>Delete Tag?</h3>
            <p>
              Are you sure you want to delete this tag? This cannot be undone.
            </p>
            <div className="form-actions">
              <Button
                variant="danger"
                onClick={() => void handleDelete(confirmDelete)}
                disabled={deleting}
              >
                {deleting ? 'Deleting…' : 'Delete'}
              </Button>
              <Button variant="secondary" onClick={() => setConfirmDelete(null)}>
                Cancel
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
