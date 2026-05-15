import { useState } from 'react';
import { useTemplateHooks, useCreateHook, useUpdateHook, useDeleteHook, type EffectHook, type HookInput } from '../hooks/useHooks';
import { useEffectDefs } from '../hooks/useEffectDefs';
import { Button } from './Button';
import { showToast } from './Toast';
import { FormField, SelectField, CheckboxField } from './fields';

const HOOK_EVENTS = [
  { value: 'on_death', label: 'On Death', group: 'Combat' },
  { value: 'on_hit_received', label: 'On Hit Received', group: 'Combat' },
  { value: 'on_hit_dealt', label: 'On Hit Dealt', group: 'Combat' },
  { value: 'on_kill', label: 'On Kill', group: 'Combat' },
  { value: 'on_enter_room', label: 'On Enter Room', group: 'Location' },
  { value: 'on_leave_room', label: 'On Leave Room', group: 'Location' },
  { value: 'on_equip', label: 'On Equip', group: 'Inventory' },
  { value: 'on_unequip', label: 'On Unequip', group: 'Inventory' },
  { value: 'on_login', label: 'On Login', group: 'Session' },
  { value: 'on_effect_start', label: 'On Effect Start', group: 'Effects' },
  { value: 'on_effect_end', label: 'On Effect End', group: 'Effects' },
];

const HOOK_TARGETS = [
  { value: 'self', label: 'Self (the character)' },
  { value: 'attacker', label: 'Attacker (hit dealer)' },
  { value: 'killer', label: 'Killer (death dealer)' },
  { value: 'room', label: 'Room (all in room)' },
  { value: 'owner', label: 'Owner (item/NPC owner)' },
];

// Group events by category
const EVENT_GROUPS = Array.from(
  new Set(HOOK_EVENTS.map((e) => e.group)),
).map((group) => ({
  label: group,
  events: HOOK_EVENTS.filter((e) => e.group === group),
}));

type HookFormProps = {
  hook: EffectHook | null
  onSubmit: (input: HookInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string | null
}

function HookForm({ hook, onSubmit, onCancel, isLoading, error }: HookFormProps) {
  const { data: effects = [] } = useEffectDefs();
  const isEdit = hook !== null;
  const [form, setForm] = useState<HookInput>(() =>
    hook
      ? { name: hook.name, event: hook.event, target: hook.target, condition: hook.condition, enabled: hook.enabled, effect_id: hook.effect_id }
      : { name: '', event: 'on_death', target: 'self', condition: '', enabled: true, effect_id: effects[0]?.id ?? 0 },
  );
  const [showAdvanced, setShowAdvanced] = useState(false);

  const set = (patch: Partial<HookInput>) => setForm((prev) => ({ ...prev, ...patch }));

  // Find effect details for display
  const selectedEffect = effects.find((e) => e.id === form.effect_id);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim() || !form.effect_id) return;
    onSubmit(form);
  };

  return (
    <form onSubmit={handleSubmit} className="bg-surface rounded-lg border border-border p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h4 className="m-0 text-lg font-semibold text-text">{isEdit ? 'Edit Hook' : 'Add Hook'}</h4>
        <button
          type="button"
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="text-sm text-primary hover:underline"
        >
          {showAdvanced ? 'Hide Advanced' : 'Show Advanced'}
        </button>
      </div>

      {error && (
        <div className="bg-danger/10 border border-danger rounded p-3 text-danger text-sm">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="md:col-span-2">
          <FormField label="Hook Name" value={form.name} onChange={(v) => set({ name: v })} required />
        </div>

        <SelectField
          label="Trigger Event"
          value={form.event}
          onChange={(v) => set({ event: v })}
          options={HOOK_EVENTS.map((e) => ({ value: e.value, label: e.label }))}
          required
        />

        <SelectField
          label="Target"
          value={form.target}
          onChange={(v) => set({ target: v })}
          options={HOOK_TARGETS.map((t) => ({ value: t.value, label: t.label }))}
        />

        <div className="md:col-span-2">
          <SelectField
            label="Effect to Apply"
            value={String(form.effect_id)}
            onChange={(v) => set({ effect_id: Number(v) })}
            options={effects.map((e) => ({ value: String(e.id), label: `${e.name} (${e.effect_type})` }))}
            required
          />
          {selectedEffect && (
            <p className="mt-1 text-xs text-text-muted">
              Type: {selectedEffect.effect_type} • {selectedEffect.description || 'No description'}
            </p>
          )}
        </div>

        {showAdvanced && (
          <>
            <div className="md:col-span-2">
              <FormField
                label="Condition (optional)"
                value={form.condition ?? ''}
                onChange={(v) => set({ condition: v })}
                placeholder="character.hp < character.max_hp * 0.3"
              />
            </div>
            <div className="md:col-span-2">
              <CheckboxField label="Enabled" checked={form.enabled} onChange={(v) => set({ enabled: v })} />
            </div>
          </>
        )}

        {!showAdvanced && <div className="md:col-span-2">
          <CheckboxField label="Enabled" checked={form.enabled} onChange={(v) => set({ enabled: v })} />
        </div>}
      </div>

      <div className="flex gap-3 pt-4 border-t border-border">
        <Button variant="primary" type="submit" disabled={isLoading}>
          {isLoading ? 'Saving…' : isEdit ? 'Update Hook' : 'Create Hook'}
        </Button>
        <Button variant="secondary" type="button" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </form>
  );
}

export function HooksPanel({ npcTemplateId }: { npcTemplateId: string }) {
  const { data: hooks = [] } = useTemplateHooks(npcTemplateId);
  const create = useCreateHook();
  const update = useUpdateHook();
  const del = useDeleteHook();
  const [showForm, setShowForm] = useState(false);
  const [editingHook, setEditingHook] = useState<EffectHook | null>(null);

  const handleCreate = (input: HookInput) => {
    create.mutate({ templateId: npcTemplateId, input }, {
      onSuccess: () => { setShowForm(false); showToast('Hook created', 'success'); },
    });
  };
  const handleUpdate = (input: HookInput) => {
    if (!editingHook) return;
    update.mutate({ id: editingHook.id, input }, {
      onSuccess: () => { setEditingHook(null); showToast('Hook updated', 'success'); },
    });
  };

  // Group hooks by event type for display
  const hooksByEvent = EVENT_GROUPS.map((group) => ({
    group: group.label,
    hooks: hooks.filter((h) => {
      const event = HOOK_EVENTS.find((e) => e.value === h.event);
      return event?.group === group.label;
    }),
  }));

  return (
    <div className="mt-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="m-0 text-text text-lg font-semibold">NPC Hooks</h2>
          <p className="text-xs text-text-muted mt-1">
            Define what happens when game events trigger on this NPC
          </p>
        </div>
        <Button variant="primary" size="sm" onClick={() => { setShowForm(true); setEditingHook(null); }}>
          + Add Hook
        </Button>
      </div>

      {showForm && !editingHook && (
        <HookForm
          hook={null}
          onSubmit={handleCreate}
          onCancel={() => setShowForm(false)}
          isLoading={create.isPending}
          error={create.error?.message ?? null}
        />
      )}
      {editingHook && (
        <HookForm
          hook={editingHook}
          onSubmit={handleUpdate}
          onCancel={() => setEditingHook(null)}
          isLoading={update.isPending}
          error={update.error?.message ?? null}
        />
      )}

      {/* Hooks list grouped by event type */}
      {!showForm && !editingHook && (
        <div className="space-y-8">
          {hooksByEvent.map((group) => (
            group.hooks.length > 0 && (
              <section key={group.group}>
                <h3 className="text-sm font-semibold text-text uppercase tracking-wider mb-4">
                  {group.group} Hooks ({group.hooks.length})
                </h3>
                <div className="space-y-3">
                  {group.hooks.map((hook) => (
                    <div
                      key={hook.id}
                      className="bg-surface-muted rounded border border-border p-4 flex items-start justify-between gap-4 hover:border-primary/50 transition-colors"
                    >
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <h4 className="font-medium text-text">{hook.name}</h4>
                          <code className="text-xs text-accent bg-surface px-2 py-0.5 rounded">
                            {hook.event}
                          </code>
                          <span className="text-xs text-text-muted">→</span>
                          <span className="text-xs font-medium text-primary">{hook.effect_name}</span>
                          <code className="text-xs text-text-muted bg-surface px-1 py-0.5 rounded">
                            {hook.target}
                          </code>
                          <span className={`text-xs px-2 py-0.5 rounded font-medium ${
                            hook.enabled
                              ? 'bg-success/10 text-success'
                              : 'bg-danger/10 text-danger'
                          }`}>
                            {hook.enabled ? 'enabled' : 'disabled'}
                          </span>
                        </div>
                        <p className="text-xs text-text-muted truncate">
                          NPC: {hook.npc_template_name}
                        </p>
                      </div>
                      <div className="flex gap-2 shrink-0">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => { setEditingHook(hook); setShowForm(false); }}
                        >
                          Edit
                        </Button>
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => del.mutate(hook.id)}
                        >
                          Delete
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </section>
            )
          ))}
        </div>
      )}

      {/* Empty state */}
      {hooks.length === 0 && !showForm && !editingHook && (
        <div className="text-center py-12 bg-surface-muted rounded-lg border border-border border-dashed">
          <p className="text-text-muted mb-4">No hooks configured for this NPC</p>
          <Button variant="primary" size="sm" onClick={() => setShowForm(true)}>
            + Configure First Hook
          </Button>
        </div>
      )}
    </div>
  );
}
