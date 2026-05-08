import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTags, useCreateTag, useUpdateTag, useDeleteTag, useTagUsages, type Tag, type TagInput } from '../../hooks/useTags'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { showToast } from '../../components/Toast'
import { TagForm } from './TagForm'
import { TagUsagesPanel, ColorDot } from './TagUsagesPanel'

export const Route = createFileRoute('/_auth/tags')({ component: TagsManagement })

function useTagMutations() {
  const create = useCreateTag()
  const update = useUpdateTag()
  const del = useDeleteTag()
  return {
    create: { mutate: create.mutate, isPending: create.isPending, error: create.error },
    update: { mutate: update.mutate, isPending: update.isPending, error: update.error },
    delete: { mutate: del.mutate, isPending: del.isPending, error: del.error },
  }
}

function TagsManagement() {
  const { data: tags = [], isLoading, error } = useTags()
  const mutations = useTagMutations()
  const [showForm, setShowForm] = useState(false)
  const [editingTag, setEditingTag] = useState<Tag | null>(null)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [viewingTag, setViewingTag] = useState<Tag | null>(null)
  const usagesQuery = useTagUsages(viewingTag?.id ?? null)

  const handleCreate = (input: TagInput) => {
    mutations.create.mutate(input, {
      onSuccess: () => { setShowForm(false); showToast('Tag created', 'success') },
    })
  }
  const handleUpdate = (input: TagInput) => {
    if (!editingTag) return
    mutations.update.mutate({ id: editingTag.id, input }, {
      onSuccess: () => { setEditingTag(null); showToast('Tag updated', 'success'); if (viewingTag?.id === editingTag.id) usagesQuery.refetch() },
    })
  }
  const handleDelete = (id: number) => {
    mutations.delete.mutate(id, {
      onSuccess: () => { setConfirmDelete(null); showToast('Tag deleted', 'success'); if (viewingTag?.id === id) setViewingTag(null) },
    })
  }

  const columns: Column<Tag>[] = [
    { header: 'Name', accessor: 'name', render: (_, r) => <span className="inline-flex items-center gap-2"><ColorDot color={r.color} />{r.name}</span> },
    { header: 'Color', accessor: 'color', render: (v) => v ? <code className="text-xs text-accent">{String(v)}</code> : <span className="text-muted">—</span> },
    { header: 'Usage', accessor: '_usages', render: (_, r) => <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setViewingTag((p) => p?.id === r.id ? null : r) }} aria-label={`View usages for ${r.name}`}>{viewingTag?.id === r.id && usagesQuery.isLoading ? 'Loading…' : 'View Usages'}</Button> },
    { header: '', accessor: '_actions', align: 'right', render: (_, r) => <div className="flex gap-2 justify-end"><Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setEditingTag(r); setShowForm(false) }} aria-label={`Edit ${r.name}`}>Edit</Button><Button variant="danger" size="sm" onClick={(e) => { e.stopPropagation(); setConfirmDelete(r.id) }} aria-label={`Delete ${r.name}`}>Delete</Button></div> },
  ]

  return (
    <div className="management-page">
      <PageHeader title="Tags" backTo="/dashboard" actions={<Button variant="primary" onClick={() => { setShowForm(true); setEditingTag(null) }}>+ Add Tag</Button>} />
      {error && <div className="error-banner">{error instanceof Error ? error.message : 'Failed to load tags'}</div>}
      {showForm && !editingTag && <TagForm tag={null} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={mutations.create.isPending} error={mutations.create.error?.message ?? null} />}
      {editingTag && <TagForm tag={editingTag} onSubmit={handleUpdate} onCancel={() => setEditingTag(null)} isLoading={mutations.update.isPending} error={mutations.update.error?.message ?? null} />}
      <DataTable columns={columns} data={tags} getKey={(r) => r.id} emptyMessage={isLoading ? 'Loading…' : 'No tags yet. Create one above.'} />
      {viewingTag && usagesQuery.data && <TagUsagesPanel tag={viewingTag} report={usagesQuery.data} onClose={() => setViewingTag(null)} />}
      {confirmDelete !== null && (
        <div className="modal-overlay"><div className="modal-card">
          <h3>Delete Tag?</h3><p>Are you sure? This cannot be undone.</p>
          {mutations.delete.error && <div className="error-banner">{mutations.delete.error.message}</div>}
          <div className="form-actions">
            <Button variant="danger" onClick={() => handleDelete(confirmDelete)} disabled={mutations.delete.isPending}>{mutations.delete.isPending ? 'Deleting…' : 'Delete'}</Button>
            <Button variant="secondary" onClick={() => setConfirmDelete(null)}>Cancel</Button>
          </div>
        </div></div>
      )}
    </div>
  )
}