import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTags, type Tag, type TagInput, fetchTagUsages, type TagUsageReport } from '../../hooks/useTags'
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

function TagUsagesPanel({
  tag,
  report,
  onClose,
}: {
  tag: Tag
  report: TagUsageReport
  onClose: () => void
}) {
  const hasUsages =
    report.skills.length > 0 || report.factions.length > 0 || report.characters.length > 0

  return (
    <div className="form-card" style={{ marginTop: '1rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h3 style={{ margin: 0 }}>
          <ColorDot color={tag.color} />
          <span style={{ marginLeft: '0.5rem' }}>{tag.name}</span>
        </h3>
        <Button variant="ghost" size="sm" onClick={onClose} aria-label="Close usages panel">
          ×
        </Button>
      </div>

      {!hasUsages && (
        <div className="empty-state">
          <p className="text-muted">This tag is orphaned — no entities reference it.</p>
        </div>
      )}

      {report.skills.length > 0 && (
        <div style={{ marginBottom: '1rem' }}>
          <h4 style={{ margin: '0 0 0.5rem 0' }}>Skills ({report.skills.length})</h4>
          <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {report.skills.map((s) => (
              <li key={`skill-${s.id}`} style={{ padding: '0.25rem 0' }}>
                <span className="badge badge-accent" style={{ marginRight: '0.5rem' }}>skill</span>
                <a href={`/abilities?id=${s.id}`} style={{ color: 'var(--color-accent)' }}>
                  {s.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}

      {report.factions.length > 0 && (
        <div style={{ marginBottom: '1rem' }}>
          <h4 style={{ margin: '0 0 0.5rem 0' }}>Factions ({report.factions.length})</h4>
          <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {report.factions.map((f) => (
              <li key={`faction-${f.id}`} style={{ padding: '0.25rem 0' }}>
                <span className="badge badge-primary" style={{ marginRight: '0.5rem' }}>faction</span>
                <a href={`/factions?id=${f.id}`} style={{ color: 'var(--color-primary)' }}>
                  {f.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}

      {report.characters.length > 0 && (
        <div style={{ marginBottom: '1rem' }}>
          <h4 style={{ margin: '0 0 0.5rem 0' }}>Characters ({report.characters.length})</h4>
          <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {report.characters.map((ch) => (
              <li key={`char-${ch.id}`} style={{ padding: '0.25rem 0' }}>
                <span className="badge badge-success" style={{ marginRight: '0.5rem' }}>character</span>
                <a href={`/characters?id=${ch.id}`} style={{ color: 'var(--color-success)' }}>
                  {ch.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
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
  const [selectedTag, setSelectedTag] = useState<Tag | null>(null)
  const [usageReport, setUsageReport] = useState<TagUsageReport | null>(null)
  const [loadingUsages, setLoadingUsages] = useState(false)

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
      // Refresh usage panel if open
      if (selectedTag?.id === editingTag.id) {
        await handleViewUsages(editingTag)
      }
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    setDeleting(true)
    try {
      await deleteTag(id)
      setConfirmDelete(null)
      if (selectedTag?.id === id) {
        setSelectedTag(null)
        setUsageReport(null)
      }
    } finally {
      setDeleting(false)
    }
  }

  const handleViewUsages = async (tag: Tag) => {
    if (selectedTag?.id === tag.id) {
      // Toggle off if already open
      setSelectedTag(null)
      setUsageReport(null)
      return
    }
    setLoadingUsages(true)
    setSelectedTag(tag)
    try {
      const report = await fetchTagUsages(tag.id)
      setUsageReport(report)
    } catch (e) {
      setUsageReport({
        tag_name: tag.name,
        total_usages: 0,
        skills: [],
        factions: [],
        characters: [],
      })
    } finally {
      setLoadingUsages(false)
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
      header: 'Usage',
      accessor: 'id',
      render: (_, row) => (
        <Button
          variant="ghost"
          size="sm"
          onClick={(e) => { e.stopPropagation(); void handleViewUsages(row) }}
          aria-label={`View usages for ${row.name}`}
        >
          {selectedTag?.id === row.id && loadingUsages
            ? 'Loading…'
            : 'View Usages'
          }
        </Button>
      ),
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

      {selectedTag && usageReport && (
        <TagUsagesPanel
          tag={selectedTag}
          report={usageReport}
          onClose={() => { setSelectedTag(null); setUsageReport(null) }}
        />
      )}

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
