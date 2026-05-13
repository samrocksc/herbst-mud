import { createFileRoute, Link, Outlet, useLocation, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useQuests,
  useCreateQuest,
  useQuestLookups,
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
  type: '', target_id: '', tag_filter: '', count: 1, labels: [], hint: '',
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
  main_type: 'general',
}

function CreateQuestForm({ onSuccess }: { onSuccess: () => void }) {
  const createQuest = useCreateQuest()
  const { data: lookups, isLoading: lookupsLoading } = useQuestLookups()
  const [formData, setFormData] = useState<QuestInput>(EMPTY_QUEST)
  const set = (patch: Partial<QuestInput>) => setFormData((prev) => ({ ...prev, ...patch }))

  if (lookupsLoading) return <div className="loading">Loading options...</div>

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

  // Get targets filtered by objective type
  const getTargetsForType = (type: string) => {
    if (!lookups) return []
    switch (type) {
      case 'kill': return lookups.npcs
      case 'explore': return lookups.rooms
      case 'collect': return lookups.items
      default: return []
    }
  }

  // Get current target options based on selected type
  const currentObjective = formData.objectives?.[0]
  const targetOptions = getTargetsForType(currentObjective?.type ?? '')

  // Convert prerequisite_quest_ids from strings to numbers for API
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createQuest.mutateAsync(formData)
      showToast('Quest created', 'success')
      setFormData(EMPTY_QUEST)
      onSuccess()
    } catch { /* toasted globally */ }
  }

  // Multi-select handlers for rewards
  const togglePrereqQuest = (questId: string) => {
    const current = formData.prerequisite_quest_ids ?? []
    if (current.includes(questId)) {
      set({ prerequisite_quest_ids: current.filter(id => id !== questId) })
    } else {
      set({ prerequisite_quest_ids: [...current, questId] })
    }
  }

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">Add New Quest</h3>
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={formData.name ?? ''} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={formData.description ?? ''} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Quest Type" value={formData.main_type ?? 'general'} onChange={(v) => set({ main_type: v })} options={[
        { value: 'general', label: 'General' },
        { value: 'hunter', label: 'Hunter (Kill NPCs)' },
        { value: 'collector', label: 'Collector (Gather Items)' },
        { value: 'explorer', label: 'Explorer (Visit Rooms)' },
      ]} />
      <SelectField label="Repeat Mode" value={formData.repeat_mode ?? 'none'} onChange={(v) => set({ repeat_mode: v })} options={REPEAT_MODE_OPTS} />
        {(formData.repeat_mode === 'cooldown') && (
          <NumberField label="Cooldown (hours)" value={formData.cooldown_hours ?? 0} onChange={(v) => set({ cooldown_hours: v })} />
        )}
        <div className="flex items-center gap-2">
          <input type="checkbox" checked={formData.is_active ?? true} onChange={(e) => set({ is_active: e.target.checked })} id="quest-active" />
          <label htmlFor="quest-active" className="text-sm text-text">Active</label>
        </div>

        {/* Prerequisite Quests */}
        <div className="border-t border-border pt-3 mt-3">
          <h4 className="text-sm font-semibold text-text mb-2">Prerequisite Quests</h4>
          <div className="flex flex-wrap gap-2">
            {lookups?.prerequisite_quests.map(q => (
              <button
                key={q.id}
                type="button"
                onClick={() => togglePrereqQuest(q.id)}
                className={`px-2 py-1 text-xs rounded border ${
                  (formData.prerequisite_quest_ids ?? []).includes(q.id)
                    ? 'bg-primary/20 border-primary text-text'
                    : 'bg-surface border-border text-muted hover:border-primary'
                }`}
              >
                {q.name}
              </button>
            ))}
          </div>
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-semibold text-text">Objectives</h4>
            <Button variant="ghost" size="sm" onClick={addObjective}>+ Objective</Button>
          </div>
          {(formData.objectives ?? []).map((obj, i) => (
            <div key={i} className="grid grid-cols-7 gap-2 mb-2 items-end">
              <SelectField
                label="Type"
                value={obj.type}
                onChange={(v) => updateObjective(i, { type: v, target_id: '' })}
                options={[
                  { value: 'kill', label: 'Kill NPC' },
                  { value: 'explore', label: 'Explore Room' },
                  { value: 'collect', label: 'Collect Item' },
                ]}
              />
              <SelectField
                label="Target"
                value={obj.target_id}
                onChange={(v) => updateObjective(i, { target_id: v })}
                options={[
                  { value: '', label: 'Select target...' },
                  ...targetOptions.map(t => ({ value: t.id, label: t.name }))
                ]}
              />
              <FormField label="Tag Filter" value={obj.tag_filter} onChange={(v) => updateObjective(i, { tag_filter: v })} placeholder="Optional: filter by tag" />
              <NumberField label="Count" value={obj.count} onChange={(v) => updateObjective(i, { count: v })} />
              <FormField label="Label" value={obj.labels?.[0] ?? ''} onChange={(v) => updateObjective(i, { labels: [v] })} placeholder="Kill Rats" />
              <FormField label="Hint" value={obj.hint} onChange={(v) => updateObjective(i, { hint: v })} placeholder="Optional hint" />
              <Button variant="danger" size="sm" onClick={() => removeObjective(i)}>×</Button>
            </div>
          ))}
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <h4 className="text-sm font-semibold text-text mb-2">Rewards</h4>
          <NumberField label="XP" value={formData.rewards?.xp ?? 0} onChange={(v) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, xp: v } })} />

          {/* Item Rewards */}
          <div className="mt-3">
            <label className="text-sm text-muted mb-1 block">Item Rewards</label>
            <div className="flex flex-wrap gap-2">
              {(lookups?.items ?? []).map(item => (
                <button
                  key={item.id}
                  type="button"
                  onClick={() => {
                    const current = formData.rewards?.item_ids ?? []
                    const newIds = current.includes(item.id)
                      ? current.filter(id => id !== item.id)
                      : [...current, item.id]
                    set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, item_ids: newIds } })
                  }}
                  className={`px-2 py-1 text-xs rounded border ${
                    (formData.rewards?.item_ids ?? []).includes(item.id)
                      ? 'bg-primary/20 border-primary text-text'
                      : 'bg-surface border-border text-muted hover:border-primary'
                  }`}
                >
                  {item.name}
                </button>
              ))}
            </div>
          </div>

          {/* Tag Add Rewards */}
          <div className="mt-3">
            <label className="text-sm text-muted mb-1 block">Tags to Add</label>
            <div className="flex flex-wrap gap-2">
              {(lookups?.tags ?? []).map(tag => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => {
                    const current = formData.rewards?.tag_adds ?? []
                    const newTags = current.includes(tag.id)
                      ? current.filter(t => t !== tag.id)
                      : [...current, tag.id]
                    set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, tag_adds: newTags } })
                  }}
                  className={`px-2 py-1 text-xs rounded border ${
                    (formData.rewards?.tag_adds ?? []).includes(tag.id)
                      ? 'bg-primary/20 border-primary text-text'
                      : 'bg-surface border-border text-muted hover:border-primary'
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          </div>
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