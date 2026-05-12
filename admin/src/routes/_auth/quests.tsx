import { createFileRoute, Link, Outlet, useLocation, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useQuests,
  useCreateQuest,
  type QuestInput,
  EMPTY_REWARDS,
  type QuestObjective,
} from '../../hooks/useQuests'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { FormField, TextareaField, NumberField, SelectField } from '../../components/FormFields'
import { showToast } from '../../components/Toast'
import type { Quest } from '../../hooks/useQuests'

export const Route = createFileRoute('/_auth/quests')({
  component: QuestsManagement,
})

const REPEAT_MODE_OPTS = [
  { value: 'none', label: 'None (one-time)' },
  { value: 'cooldown', label: 'Cooldown' },
  { value: 'always', label: 'Always repeatable' },
]

const EMPTY_OBJECTIVE: QuestObjective = {
  type: '', target_id: '', count: 1, label: '', hint: '',
}

const EMPTY_QUEST: QuestInput = {
  name: '',
  description: '',
  prerequisite_quest_ids: [],
  objectives: [{ ...EMPTY_OBJECTIVE }],
  rewards: { ...EMPTY_REWARDS },
  repeat_mode: 'none',
  cooldown_hours: 0,
  is_active: true,
}

function CreateQuestForm({ onSuccess }: { onSuccess: () => void }) {
  const createQuest = useCreateQuest()
  const [formData, setFormData] = useState<QuestInput>(EMPTY_QUEST)
  const set = (patch: Partial<QuestInput>) => setFormData((prev) => ({ ...prev, ...patch }))

  const addObjective = () => {
    const objs = [...(formData.objectives ?? []), { ...EMPTY_OBJECTIVE }]
    set({ objectives: objs })
  }
  const updateObjective = (i: number, patch: Partial<QuestObjective>) => {
    const objs = formData.objectives?.map((o, idx) => idx === i ? { ...o, ...patch } : o) ?? []
    set({ objectives: objs })
  }
  const removeObjective = (i: number) => {
    const objs = formData.objectives?.filter((_, idx) => idx !== i) ?? []
    set({ objectives: objs })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createQuest.mutateAsync(formData)
      showToast('Quest created', 'success')
      setFormData(EMPTY_QUEST)
      onSuccess()
    } catch { /* toasted globally */ }
  }

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">Add New Quest</h3>
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={formData.name ?? ''} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={formData.description ?? ''} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Repeat Mode" value={formData.repeat_mode ?? 'none'} onChange={(v) => set({ repeat_mode: v })} options={REPEAT_MODE_OPTS} />
        {(formData.repeat_mode === 'cooldown') && (
          <NumberField label="Cooldown (hours)" value={formData.cooldown_hours ?? 0} onChange={(v) => set({ cooldown_hours: v })} />
        )}
        <div className="flex items-center gap-2">
          <input type="checkbox" checked={formData.is_active ?? true} onChange={(e) => set({ is_active: e.target.checked })} id="quest-active" />
          <label htmlFor="quest-active" className="text-sm text-text">Active</label>
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-semibold text-text">Objectives</h4>
            <Button variant="ghost" size="sm" onClick={addObjective}>+ Objective</Button>
          </div>
          {(formData.objectives ?? []).map((obj, i) => (
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
          <NumberField label="XP" value={formData.rewards?.xp ?? 0} onChange={(v) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, xp: v } })} />
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={createQuest.isPending} fullWidth>
            {createQuest.isPending ? 'Creating...' : 'Create Quest'}
          </Button>
        </div>
      </form>
    </div>
  )
}

const COLUMNS: Column<Quest>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_, row) => (
      <Link to="/quests/$questId" params={{ questId: String(row.id) }} className="no-underline text-primary hover:underline font-bold">
        {row.name}
      </Link>
    ),
  },
  { header: 'Active', accessor: 'is_active', render: (val) => val ? '✓' : '✗' },
  { header: 'Repeat', accessor: 'repeat_mode' },
  { header: 'Objectives', accessor: 'objectives', render: (val) => (val as unknown as unknown[])?.length ?? 0 },
  { header: 'XP', accessor: 'rewards', render: (val) => (val as { xp?: number })?.xp ?? 0 },
]

function QuestsManagement() {
  const [showCreate, setShowCreate] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { data: quests, isLoading, error } = useQuests()

  if (location.pathname !== '/quests') return <Outlet />

  if (isLoading) return <div className="loading">Loading quests...</div>
  if (error) return <div className="error">Failed to load quests: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Quests"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => setShowCreate(!showCreate)}>
            {showCreate ? 'Cancel' : '+ Add Quest'}
          </Button>
        }
      />
      {showCreate && <CreateQuestForm onSuccess={() => setShowCreate(false)} />}
      <DataTable
        columns={COLUMNS}
        data={quests ?? []}
        getKey={(row) => row.id}
        onRowClick={(row) => navigate({ to: '/quests/$questId', params: { questId: String(row.id) } })}
        emptyMessage="No quests found. Create your first quest!"
      />
    </div>
  )
}