import { useState } from 'react'
import { useTemplateHooks, useCreateHook, useUpdateHook, useDeleteHook, type EffectHook, type HookInput } from '../../hooks/useHooks'
import { useEffectDefs } from '../../hooks/useEffectDefs'
import { Button } from '../Button'
import { DeleteConfirmation } from '../DeleteConfirmation'
import { showToast } from '../Toast'
import { FormField, SelectField, CheckboxField } from '../fields'

const HOOK_EVENTS = [
  { value: 'on_death', label: 'On Death' },
  { value: 'on_hit_received', label: 'On Hit Received' },
  { value: 'on_hit_dealt', label: 'On Hit Dealt' },
  { value: 'on_kill', label: 'On Kill' },
  { value: 'on_enter_room', label: 'On Enter Room' },
  { value: 'on_leave_room', label: 'On Leave Room' },
  { value: 'on_equip', label: 'On Equip' },
  { value: 'on_unequip', label: 'On Unequip' },
  { value: 'on_login', label: 'On Login' },
  { value: 'on_effect_start', label: 'On Effect Start' },
  { value: 'on_effect_end', label: 'On Effect End' },
]

const HOOK_TARGETS = [
  { value: 'self', label: 'Self' },
  { value: 'attacker', label: 'Attacker' },
  { value: 'killer', label: 'Killer' },
  { value: 'room', label: 'Room' },
  { value: 'owner', label: 'Owner' },
]

type HookFormProps = {
  hook: EffectHook | null
  npcTemplateId: string
  onSubmit: (input: HookInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string | null
}

function HookForm({ hook, npcTemplateId, onSubmit, onCancel, isLoading, error }: HookFormProps) {
  const { data: effects = [] } = useEffectDefs()
  const isEdit = hook !== null
  const [form, setForm] = useState<HookInput>(() =>
    hook
      ? { name: hook.name, event: hook.event, target: hook.target, condition: hook.condition, enabled: hook.enabled, effect_id: hook.effect_id }
      : { name: '', event: 'on_death', target: 'self', condition: '', enabled: true, effect_id: effects[0]?.id ?? 0 },
  )
  const set = (patch: Partial<HookInput>) => setForm((prev) => ({ ...prev, ...patch }))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim() || !form.effect_id) return
    onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-surface rounded border border-border p-4 space-y-3">
      <h4 className="m-0 text-text font-semibold">{isEdit ? 'Edit Hook' : 'Add Hook'}</h4>
      {error && <div className="error-banner">{error}</div>}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormField label="Name" value={form.name} onChange={(v) => set({ name: v })} required />
        <SelectField label="Event" value={form.event} onChange={(v) => set({ event: v })} options={HOOK_EVENTS} />
        <SelectField label="Target" value={form.target} onChange={(v) => set({ target: v })} options={HOOK_TARGETS} />
        <SelectField
          label="Effect"
          value={String(form.effect_id)}
          onChange={(v) => set({ effect_id: Number(v) })}
          options={effects.map((e) => ({ value: String(e.id), label: e.name }))}
        />
        <FormField label="Condition (optional)" value={form.condition ?? ''} onChange={(v) => set({ condition: v })} />
        <div className="flex items-end">
          <CheckboxField label="Enabled" checked={form.enabled} onChange={(v) => set({ enabled: v })} />
        </div>
      </div>
      <div className="flex gap-2">
        <Button variant="primary" type="submit" disabled={isLoading}>{isLoading ? 'Saving…' : isEdit ? 'Update' : 'Create'}</Button>
        <Button variant="secondary" type="button" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  )
}

export function HooksPanel({ npcTemplateId }: { npcTemplateId: string }) {
  const { data: hooks = [], isLoading, error } = useTemplateHooks(npcTemplateId)
  const create = useCreateHook()
  const update = useUpdateHook()
  const del = useDeleteHook()
  const [showForm, setShowForm] = useState(false)
  const [editingHook, setEditingHook] = useState<EffectHook | null>(null)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)

  const handleCreate = (input: HookInput) => {
    create.mutate({ templateId: npcTemplateId, input }, {
      onSuccess: () => { setShowForm(false); showToast('Hook created', 'success') },
    })
  }
  const handleUpdate = (input: HookInput) => {
    if (!editingHook) return
    update.mutate({ id: editingHook.id, input }, {
      onSuccess: () => { setEditingHook(null); showToast('Hook updated', 'success') },
    })
  }
  const handleDelete = (id: number) => {
    del.mutate(id, {
      onSuccess: () => { setConfirmDelete(null); showToast('Hook deleted', 'success') },
    })
  }

  return (
    <div className="mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="m-0 text-text text-lg font-semibold">Hooks</h2>
        <Button variant="primary" size="sm" onClick={() => { setShowForm(true); setEditingHook(null) }}>+ Add Hook</Button>
      </div>
      {error && <div className="error-banner mb-3">{error instanceof Error ? error.message : 'Failed to load hooks'}</div>}
      {showForm && !editingHook && (
        <HookForm hook={null} npcTemplateId={npcTemplateId} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={create.isPending} error={create.error?.message ?? null} />
      )}
      {editingHook && (
        <HookForm hook={editingHook} npcTemplateId={npcTemplateId} onSubmit={handleUpdate} onCancel={() => setEditingHook(null)} isLoading={update.isPending} error={update.error?.message ?? null} />
      )}
      {isLoading ? (
        <div className="text-text-muted text-sm py-4 text-center">Loading hooks…</div>
      ) : hooks.length === 0 ? (
        <div className="text-text-muted text-sm py-4 text-center">No hooks yet. Add one above.</div>
      ) : (
        <div className="space-y-2">
          {hooks.map((hook) => (
            <div key={hook.id} className="bg-surface-muted rounded border border-border p-3 flex items-center justify-between gap-4">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-text">{hook.name}</span>
                  <code className="text-xs text-accent bg-surface px-1.5 py-0.5 rounded">{hook.event}</code>
                  <span className="text-xs text-text-muted">→ {hook.effect_name}</span>
                  <code className="text-xs text-text-muted">({hook.target})</code>
                  <span className={`text-xs px-1.5 py-0.5 rounded ${hook.enabled ? 'bg-green-900/30 text-green-400' : 'bg-red-900/30 text-red-400'}`}>
                    {hook.enabled ? 'enabled' : 'disabled'}
                  </span>
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="ghost" size="sm" onClick={() => { setEditingHook(hook); setShowForm(false) }}>Edit</Button>
                <Button variant="danger" size="sm" onClick={() => setConfirmDelete(hook.id)}>Delete</Button>
              </div>
            </div>
          ))}
        </div>
      )}
      {confirmDelete !== null && (
        <DeleteConfirmation open={confirmDelete !== null} title="Delete Hook" message="Are you sure? This will remove the event trigger." onConfirm={() => handleDelete(confirmDelete)} onCancel={() => setConfirmDelete(null)} isLoading={del.isPending} />
      )}
    </div>
  )
}