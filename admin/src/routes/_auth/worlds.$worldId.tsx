import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useWorld, useUpdateWorld, useDeleteWorld, type WorldInput } from '../../hooks/useWorlds'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'
import { FormField, TextareaField } from '../../components/FormFields'
import { showToast } from '../../components/Toast'

export const Route = createFileRoute('/_auth/worlds/$worldId')({
  component: EditWorldPage,
})

function EditWorldPage() {
  const { worldId } = Route.useParams()
  const id = Number(worldId)
  const navigate = useNavigate()
  const { data: world, isLoading, error } = useWorld(id)
  const updateWorld = useUpdateWorld()
  const deleteWorld = useDeleteWorld()
  const [formData, setFormData] = useState<WorldInput | null>(null)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  const set = (patch: Partial<WorldInput>) => setFormData((prev) => prev ? { ...prev, ...patch } : prev)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData) return
    try {
      await updateWorld.mutateAsync({ id, input: formData })
      showToast('World updated', 'success')
      navigate({ to: '/worlds' })
    } catch {
      // Error is toasted by global onError handler
    }
  }

  const handleDelete = async () => {
    try {
      await deleteWorld.mutateAsync(id)
      showToast('World deleted', 'success')
      navigate({ to: '/worlds' })
    } catch {
      // Error is toasted by global onError handler
    }
  }

  if (isLoading) return <div className="p-8 text-text-muted">Loading world...</div>
  if (error || !world) return <div className="p-8 text-danger">Error: {error?.message ?? 'World not found'}</div>

  const currentData = formData ?? {
    id: world.id,
    name: world.name,
    title: world.title,
    description: world.description,
    active: world.active,
  }

  return (
    <div className="p-6 max-w-[800px] mx-auto">
      <PageHeader
        title={`Edit: ${world.name}`}
        showBack
        backTo="/worlds"
        actions={
          <Button variant="danger" onClick={() => setShowDeleteConfirm(true)}>
            Delete
          </Button>
        }
      />

      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              label="Name *"
              value={currentData.name}
              onChange={(v) => set({ name: v })}
              required
              placeholder="e.g., My Fantasy World"
            />
            <FormField
              label="Title"
              value={currentData.title}
              onChange={(v) => set({ title: v })}
              placeholder="Human-readable title"
            />
          </div>
          <TextareaField
            label="Description"
            value={currentData.description}
            onChange={(v) => set({ description: v })}
            rows={3}
            placeholder="Describe this world..."
          />

          <div className="flex items-center gap-2 mt-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={currentData.active}
                onChange={(e) => set({ active: e.target.checked })}
                className="accent-primary"
              />
              <span className="text-text">Active</span>
            </label>
            <span className="text-text-muted text-sm">(Enable this world for players)</span>
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: '/worlds' })}>
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={updateWorld.isPending}>
              {updateWorld.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </form>
      </div>

      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-surface p-6 rounded border border-border max-w-sm">
            <h3 className="text-lg font-semibold text-text mb-4">Delete World?</h3>
            <p className="text-text-muted mb-6">
              Are you sure you want to delete "{world.name}"? This action cannot be undone.
            </p>
            <div className="flex gap-2 justify-end">
              <Button variant="secondary" onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </Button>
              <Button variant="danger" onClick={handleDelete} disabled={deleteWorld.isPending}>
                {deleteWorld.isPending ? 'Deleting...' : 'Delete'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}