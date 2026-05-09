import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useEffectDefs,
  useCreateEffectDef,
  useUpdateEffectDef,
  useDeleteEffectDef,
  EMPTY_INPUT,
  type EffectDef,
  type EffectDefInput,
} from '../../hooks/useEffectDefs'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { DeleteConfirmation } from '../../components/DeleteConfirmation'
import { showToast } from '../../components/Toast'
import { FormField, NumberField, TextareaField, SelectField, CheckboxField } from '../../components/fields'

export const Route = createFileRoute('/_auth/effects')({ component: EffectsManagement })

const EFFECT_TYPES = [
  { value: 'xp_drain', label: 'XP Drain' },
  { value: 'xp_gain', label: 'XP Gain' },
  { value: 'xp_set', label: 'XP Set' },
  { value: 'bind_point_set', label: 'Bind Point Set' },
  { value: 'hp_change', label: 'HP Change' },
  { value: 'stamina_change', label: 'Stamina Change' },
  { value: 'mana_change', label: 'Mana Change' },
  { value: 'message', label: 'Message' },
  { value: 'teleport', label: 'Teleport' },
  { value: 'apply_effect', label: 'Apply Effect' },
  { value: 'tag_add', label: 'Tag Add' },
  { value: 'tag_remove', label: 'Tag Remove' },
]

const STACK_MODES = [
  { value: 'replace', label: 'Replace' },
  { value: 'refresh', label: 'Refresh' },
  { value: 'stack', label: 'Stack' },
]

type ParamConfig = { key: string; label: string; type: 'number' | 'text' }

function getParamsForType(type: string): ParamConfig[] {
  switch (type) {
    case 'xp_drain': case 'xp_gain': case 'xp_set':
      return [{ key: 'amount', label: 'Amount', type: 'number' }]
    case 'bind_point_set': case 'teleport':
      return [{ key: 'room_id', label: 'Room ID', type: 'number' }]
    case 'hp_change': case 'stamina_change': case 'mana_change':
      return [{ key: 'amount', label: 'Amount (+/-)', type: 'number' }]
    case 'message':
      return [{ key: 'text', label: 'Message Text', type: 'text' }, { key: 'message_type', label: 'Message Type', type: 'text' }]
    case 'apply_effect':
      return [{ key: 'effect_id', label: 'Effect ID', type: 'number' }]
    case 'tag_add': case 'tag_remove':
      return [{ key: 'tag_name', label: 'Tag Name', type: 'text' }]
    default:
      return []
  }
}

function useEffectDefMutations() {
  const create = useCreateEffectDef()
  const update = useUpdateEffectDef()
  const del = useDeleteEffectDef()
  return {
    create: { mutate: create.mutate, isPending: create.isPending, error: create.error },
    update: { mutate: update.mutate, isPending: update.isPending, error: update.error },
    delete: { mutate: del.mutate, isPending: del.isPending, error: del.error },
  }
}

function EffectDefForm({ effect, onSubmit, onCancel, isLoading, error }: {
  effect: EffectDef | null
  onSubmit: (input: EffectDefInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string | null
}) {
  const isEdit = effect !== null
  const [form, setForm] = useState<EffectDefInput>(() =>
    effect ? { ...effect, parameters: { ...effect.parameters }, messages: { ...effect.messages } } : { ...EMPTY_INPUT },
  )
  const set = (patch: Partial<EffectDefInput>) => setForm((prev) => ({ ...prev, ...patch }))
  const paramConfigs = getParamsForType(form.effect_type)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim() || !form.effect_type) return
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-surface rounded-lg border border-border p-4 mb-4 space-y-4">
      <h3 className="text-lg font-semibold text-text">{isEdit ? 'Edit Effect' : 'Create Effect'}</h3>
      {error && <div className="error-banner">{error}</div>}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField label="Name" value={form.name} onChange={(v) => set({ name: v })} required />
        <SelectField label="Type" value={form.effect_type} onChange={(v) => set({ effect_type: v, parameters: {} })} options={EFFECT_TYPES} required />
        <div className="md:col-span-2">
          <TextareaField label="Description" value={form.description} onChange={(v) => set({ description: v })} />
        </div>
      </div>
      {paramConfigs.length > 0 && (
        <div className="space-y-3">
          <h4 className="text-sm font-medium text-text-muted">Parameters</h4>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {paramConfigs.map((pc) => (
              pc.type === 'number'
                ? <NumberField key={pc.key} label={pc.label} value={Number(form.parameters[pc.key] ?? 0)} onChange={(v) => set({ parameters: { ...form.parameters, [pc.key]: v } })} />
                : <FormField key={pc.key} label={pc.label} value={String(form.parameters[pc.key] ?? '')} onChange={(v) => set({ parameters: { ...form.parameters, [pc.key]: v } })} />
            ))}
          </div>
        </div>
      )}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SelectField label="Stack Mode" value={form.stack_mode} onChange={(v) => set({ stack_mode: v })} options={STACK_MODES} />
        <NumberField label="Stack Limit" value={form.stack_limit} onChange={(v) => set({ stack_limit: v })} />
        <CheckboxField label="Permanent" checked={form.is_permanent} onChange={(v) => set({ is_permanent: v })} />
      </div>
      {!form.is_permanent && (
        <NumberField label="Duration (seconds, 0=instant)" value={form.duration_secs} onChange={(v) => set({ duration_secs: v })} />
      )}
      <div className="space-y-3">
        <h4 className="text-sm font-medium text-text-muted">Messages</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="On Start" value={form.messages.on_start ?? ''} onChange={(v) => set({ messages: { ...form.messages, on_start: v } })} />
          <FormField label="On End" value={form.messages.on_end ?? ''} onChange={(v) => set({ messages: { ...form.messages, on_end: v } })} />
        </div>
      </div>
      <div className="flex gap-2">
        <Button variant="primary" type="submit" disabled={isLoading}>{isLoading ? 'Saving…' : isEdit ? 'Update' : 'Create'}</Button>
        <Button variant="secondary" type="button" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  )
}

function EffectsManagement() {
  const { data: effects = [], isLoading, error } = useEffectDefs()
  const mutations = useEffectDefMutations()
  const [showForm, setShowForm] = useState(false)
  const [editingEffect, setEditingEffect] = useState<EffectDef | null>(null)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)

  const handleCreate = (input: EffectDefInput) => {
    mutations.create.mutate(input, {
      onSuccess: () => { setShowForm(false); showToast('Effect created', 'success') },
    })
  }
  const handleUpdate = (input: EffectDefInput) => {
    if (!editingEffect) return
    mutations.update.mutate({ id: editingEffect.id, input }, {
      onSuccess: () => { setEditingEffect(null); showToast('Effect updated', 'success') },
    })
  }
  const handleDelete = (id: number) => {
    mutations.delete.mutate(id, {
      onSuccess: () => { setConfirmDelete(null); showToast('Effect deleted', 'success') },
    })
  }

  const columns: Column<EffectDef>[] = [
    { header: 'Name', accessor: 'name' },
    { header: 'Type', accessor: 'effect_type', render: (v) => <code className="text-xs text-accent">{String(v)}</code> },
    { header: 'Stack', accessor: 'stack_mode', render: (v, r) => <span>{String(v)}{r.stack_limit > 1 ? ` (×${r.stack_limit})` : ''}</span> },
    { header: 'Duration', accessor: 'duration_secs', render: (v, r) => r.is_permanent ? 'Permanent' : v === 0 ? 'Instant' : `${v}s` },
    { header: 'Hooks', accessor: 'hook_count', render: (v) => <span className="text-text-muted">{Number(v)}</span> },
    { header: '', accessor: '_actions', align: 'right', render: (_, r) => (
      <div className="flex gap-2 justify-end">
        <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setEditingEffect(r); setShowForm(false) }}>Edit</Button>
        <Button variant="danger" size="sm" onClick={(e) => { e.stopPropagation(); setConfirmDelete(r.id) }}>Delete</Button>
      </div>
    )},
  ]

  return (
    <div className="management-page">
      <PageHeader title="Effects" backTo="/dashboard" actions={<Button variant="primary" onClick={() => { setShowForm(true); setEditingEffect(null) }}>+ Create Effect</Button>} />
      {error && <div className="error-banner">{error instanceof Error ? error.message : 'Failed to load effects'}</div>}
      {showForm && !editingEffect && <EffectDefForm effect={null} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={mutations.create.isPending} error={mutations.create.error?.message ?? null} />}
      {editingEffect && <EffectDefForm effect={editingEffect} onSubmit={handleUpdate} onCancel={() => setEditingEffect(null)} isLoading={mutations.update.isPending} error={mutations.update.error?.message ?? null} />}
      <DataTable columns={columns} data={effects} getKey={(r) => r.id} emptyMessage={isLoading ? 'Loading…' : 'No effects yet. Create one above.'} />
      {confirmDelete !== null && (
        <DeleteConfirmation open={confirmDelete !== null} title="Delete Effect" message="Are you sure? Any hooks referencing this effect will break." onConfirm={() => handleDelete(confirmDelete)} onCancel={() => setConfirmDelete(null)} isLoading={mutations.delete.isPending} />
      )}
    </div>
  )
}