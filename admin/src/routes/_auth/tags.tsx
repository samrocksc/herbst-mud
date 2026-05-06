import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTags, type Tag, type TagInput, fetchTagUsages, type TagUsageReport } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/tags')({
  component: TagsManagement,
})

const DEFAULT_COLOR = 'var(--color-tag-default)'

function ColorDot({ color }: { color: string }) {
  const dotStyle = { '--dot-color': color || DEFAULT_COLOR } as React.CSSProperties
  return (
    <span
      className="inline-block w-3 h-3 rounded-full shrink-0 bg-(--dot-color)"
      style={dotStyle}
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
    <div className="form-card mt-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="m-0">
          <ColorDot color={tag.color} />
          <span className="ml-2">{tag.name}</span>
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
        <div className="mb-4">
          <h4 className="m-0 mb-2">Skills ({report.skills.length})</h4>
          <ul className="list-none p-0 m-0">
            {report.skills.map((s) => (
              <li key={`skill-${s.id}`} className="py-1">
                <span className="badge badge-accent mr-2">skill</span>
                <a href={`/abilities?id=${s.id}`} className="text-accent">
                  {s.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}

      {report.factions.length > 0 && (
        <div className="mb-4">
          <h4 className="m-0 mb-2">Factions ({report.factions.length})</h4>
          <ul className="list-none p-0 m-0">
            {report.factions.map((f) => (
              <li key={`faction-${f.id}`} className="py-1">
                <span className="badge badge-primary mr-2">faction</span>
                <a href={`/factions?id=${f.id}`} className="text-primary">
                  {f.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}

      {report.characters.length > 0 && (
        <div className="mb-4">
          <h4 className="m-0 mb-2">Characters ({report.characters.length})</h4>
          <ul className="list-none p-0 m-0">
            {report.characters.map((ch) => (
              <li key={`char-${ch.id}`} className="py-1">
                <span className="badge badge-success mr-2">character</span>
                <a href={`/characters?id=${ch.id}`} className="text-success">
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
          <div className="flex items-center gap-2">
            <input
              type="color"
              value={form.color || DEFAULT_COLOR}
              onChange={(e) => setForm({ ...form, color: e.target.value })}
              className="w-10 h-8 p-0.5 cursor-pointer"
            />
            <input
              type="text"
              value={form.color}
              onChange={(e) => setForm({ ...form, color: e.target.value })}
              placeholder="CSS color / hex"
              pattern="^#[0-9a-fA-F]{6}$"
              className="w-28"
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
        <span className="inline-flex items-center gap-2">
          <ColorDot color={row.color} />
          {row.name}
        </span>
      ),
    },
    {
      header: 'Color',
      accessor: 'color',
      render: (val: unknown) => val ? (
        <code className="text-xs text-accent">{String(val)}</code>
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
        <div className="flex gap-2 justify-end">
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
