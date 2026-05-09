import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useQuest,
  useUpdateQuest,
  useDeleteQuest,
  type QuestInput,
  EMPTY_REWARDS,
  type QuestObjective,
} from '../../hooks/useQuests'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'
import { FormField, TextareaField, NumberField, SelectField } from '../../components/FormFields'
import { showToast } from '../../components/Toast'

export const Route = createFileRoute('/_auth/quests/$questId')({
  component: QuestDetailPage,
})

const REPEAT_MODE_OPTS = [
  { value: 'none', label: 'None (one-time)' },
  { value: 'cooldown', label: 'Cooldown' },
  { value: 'always', label: 'Always repeatable' },
]

const EMPTY_OBJECTIVE: QuestObjective = {
  type: '', target_id: '', count: 1, label: '', hint: '',
}

function QuestDetailPage() {
  const questId = Route.useParams().questId
  const navigate = useNavigate()
  const { data: quest, isLoading, error } = useQuest(Number(questId))
  const updateQuest = useUpdateQuest()
  const deleteQuest = useDeleteQuest()

  const [formData, setFormData] = useState<QuestInput | null>(null)
  const [confirmDelete, setConfirmDelete] = useState(false)

  if (isLoading) return <div className="loading">Loading quest...</div>
  if (error) return <div className="error">Failed to load quest: {error.message}</div>
  if (!quest) return <div className="error">Quest not found</div>

  const current = formData ?? {
    name: quest.name,
    description: quest.description,
    prerequisite_quest_ids: quest.prerequisite_quest_ids ?? [],
    objectives: quest.objectives ?? [],
    rewards: quest.rewards ?? EMPTY_REWARDS,
    repeat_mode: quest.repeat_mode,
    cooldown_hours: quest.cooldown_hours,
    is_active: quest.is_active,
  }

  const set = (patch: Partial<QuestInput>) => setFormData({ ...current, ...patch })

  const addObjective = () => {
    const objs = [...(current.objectives ?? []), { ...EMPTY_OBJECTIVE }]
    set({ objectives: objs })
  }
  const updateObjective = (i: number, patch: Partial<QuestObjective>) => {
    const objs = current.objectives?.map((o, idx) => idx === i ? { ...o, ...patch } : o) ?? []
    set({ objectives: objs })
  }
  const removeObjective = (i: number) => {
    const objs = current.objectives?.filter((_, idx) => idx !== i) ?? []
    set({ objectives: objs })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await updateQuest.mutateAsync({ id: Number(questId), input: current })
      showToast('Quest updated', 'success')
      setFormData(null)
    } catch { /* toasted globally */ }
  }

  const handleDelete = async () => {
    try {
      await deleteQuest.mutateAsync(Number(questId))
      showToast('Quest deleted', 'success')
      navigate({ to: '/quests' })
    } catch { /* toasted globally */ }
  }

  return (
    <div className="management-page">
      <PageHeader title={quest.name} backTo="/quests" />
      <form onSubmit={handleSubmit} className="form-card space-y-3">
        <FormField label="Name" value={current.name ?? ''} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={current.description ?? ''} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Repeat Mode" value={current.repeat_mode ?? 'none'} onChange={(v) => set({ repeat_mode: v })} options={REPEAT_MODE_OPTS} />
        {(current.repeat_mode === 'cooldown') && (
          <NumberField label="Cooldown (hours)" value={current.cooldown_hours ?? 0} onChange={(v) => set({ cooldown_hours: v })} />
        )}
        <div className="flex items-center gap-2">
          <input type="checkbox" checked={current.is_active ?? true} onChange={(e) => set({ is_active: e.target.checked })} id="quest-active-edit" />
          <label htmlFor="quest-active-edit" className="text-sm text-text">Active</label>
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-semibold text-text">Objectives</h4>
            <Button variant="ghost" size="sm" onClick={addObjective}>+ Objective</Button>
          </div>
          {(current.objectives ?? []).map((obj, i) => (
            <div key={i} className="grid grid-cols-5 gap-2 mb-2 items-end">
              <FormField label="Type" value={obj.type} onChange={(v) => updateObjective(i, { type: v })} placeholder="kill" />
              <FormField label="Target" value={obj.target_id} onChange={(v) => updateObjective(i, { target_id: v })} placeholder="rat" />
              <NumberField label="Count" value={obj.count} onChange={(v) => updateObjective(i, { count: v })} />
              <FormField label="Label" value={obj.label} onChange={(v) => updateObjective(i, { label: v })} placeholder="Kill Rats" />
              <Button variant="danger" size="sm" onClick={() => removeObjective(i)}>×</Button>
            </div>
          ))}
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <h4 className="text-sm font-semibold text-text mb-2">Rewards</h4>
          <NumberField label="XP" value={current.rewards?.xp ?? 0} onChange={(v) => set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, xp: v } })} />
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={updateQuest.isPending}>
            {updateQuest.isPending ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button variant="danger" onClick={() => setConfirmDelete(true)}>Delete Quest</Button>
        </div>
      </form>

      {confirmDelete && (
        <div className="modal-overlay" onClick={() => setConfirmDelete(false)}>
          <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Delete Quest</h3>
              <Button variant="ghost" size="sm" onClick={() => setConfirmDelete(false)} aria-label="Close">×</Button>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to delete <strong>{quest.name}</strong>?</p>
              <p className="text-muted">This action cannot be undone.</p>
            </div>
            <div className="modal-footer">
              <Button variant="danger" onClick={handleDelete} disabled={deleteQuest.isPending}>
                {deleteQuest.isPending ? 'Deleting...' : 'Delete'}
              </Button>
              <Button variant="secondary" onClick={() => setConfirmDelete(false)}>Cancel</Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}