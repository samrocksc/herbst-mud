import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useAchievements, useCreateAchievement, useUpdateAchievement,
  useDeleteAchievement, type Achievement, type AchievementInput,
} from '../../hooks/useAchievements'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { showToast } from '../../components/Toast'
import { COLUMNS } from './AchievementColumns'
import { AchievementForm, DeleteConfirmation } from './AchievementForm'

export const Route = createFileRoute('/_auth/achievements')({ component: AchievementsManagement })

function AchievementsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Achievement | null>(null)
  const [deleting, setDeleting] = useState<Achievement | null>(null)
  const [formError, setFormError] = useState('')
  const create = useCreateAchievement()
  const update = useUpdateAchievement()
  const remove = useDeleteAchievement()
  const { data: achievements, isLoading, error } = useAchievements()

  const handleSubmit = async (formData: AchievementInput) => {
    setFormError('')
    try {
      if (editing) { await update.mutateAsync({ id: editing.id, input: formData }); showToast('Achievement updated', 'success') }
      else { await create.mutateAsync(formData); showToast('Achievement created', 'success') }
      setShowForm(false); setEditing(null)
    } catch (err) { setFormError(err instanceof Error ? err.message : 'Failed to save achievement') }
  }

  const handleDelete = async () => {
    if (!deleting) return
    try { await remove.mutateAsync(deleting.id); setDeleting(null) } catch { /* global onError toasts */ }
  }

  if (isLoading) return <div className="loading">Loading achievements...</div>
  if (error) return <div className="error">Failed to load achievements: {error.message}</div>

  const columns: Column<Achievement>[] = [
    ...COLUMNS.slice(0, 4),
    { header: 'Actions', accessor: '_actions',
      render: (_: unknown, row: Achievement) => (
        <span className="inline-flex gap-2">
          <Button variant="accent" size="sm" onClick={() => { setEditing(row); setShowForm(true) }}>Edit</Button>
          <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeleting(row)}>Delete</Button>
        </span>),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader title="Achievements" backTo="/dashboard"
        actions={<Button variant="primary" onClick={() => { setEditing(null); setShowForm(true) }}>+ Add Achievement</Button>} />
      {showForm && <AchievementForm achievement={editing} onSubmit={handleSubmit}
        onCancel={() => { setShowForm(false); setEditing(null); setFormError('') }}
        isLoading={create.isPending || update.isPending} error={formError} />}
      <DataTable columns={columns} data={achievements ?? []} getKey={(row: Achievement) => row.id}
        emptyMessage="No achievements found. Add your first achievement!" />
      {deleting && <DeleteConfirmation achievement={deleting} onConfirm={handleDelete}
        onCancel={() => setDeleting(null)} isLoading={remove.isPending} />}
    </div>
  )
}