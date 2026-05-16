/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import {
  useEffectDefs,
  useCreateEffectDef,
  useUpdateEffectDef,
  useDeleteEffectDef,
  createEmptyInput,
  type EffectDef,
  type EffectDefInput,
  EFFECT_TYPES,
  STACK_MODES,
  MESSAGE_TYPES,
} from "../../hooks/useEffectDefs";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { showToast } from "../../components/Toast";
import { FormField, NumberField, TextareaField, SelectField, CheckboxField, FieldLabel } from "../../components/fields";

// Group effect types for better organization
const EFFECT_TYPE_GROUPS = Array.from(
  new Set(EFFECT_TYPES.map((t) => t.group)),
).map((group) => ({
  label: group,
  types: EFFECT_TYPES.filter((t) => t.group === group),
}));

// Get parameter fields based on effect type
function getParamConfig(effectType: string) {
  switch (effectType) {
    case "hp_change":
    case "stamina_change":
    case "mana_change":
      return [
        { key: "amount", label: "Amount (+/-)", type: "number" as const },
      ];
    case "xp_drain":
    case "xp_gain":
    case "xp_set":
      return [
        { key: "amount", label: "Amount", type: "number" as const },
      ];
    case "bind_point_set":
    case "teleport":
      return [
        { key: "room_id", label: "Room ID", type: "number" as const },
      ];
    case "message":
      return [
        { key: "text", label: "Message Text", type: "text" as const },
        { key: "message_type", label: "Message Type", type: "select" as const, options: MESSAGE_TYPES },
      ];
    case "room_message":
      return [
        { key: "text", label: "Message Text", type: "text" as const },
        { key: "message_type", label: "Message Type", type: "select" as const, options: MESSAGE_TYPES },
      ];
    case "whisper":
      return [
        { key: "text", label: "Message Text", type: "text" as const },
        { key: "target", label: "Target Name", type: "text" as const },
      ];
    case "apply_effect":
      return [
        { key: "effect_id", label: "Effect ID", type: "number" as const },
      ];
    case "tag_add":
    case "tag_remove":
      return [
        { key: "tag_name", label: "Tag Name", type: "text" as const },
      ];
    default:
      return [];
  }
}

function useEffectDefMutations() {
  const create = useCreateEffectDef();
  const update = useUpdateEffectDef();
  const del = useDeleteEffectDef();
  return {
    create: { mutate: create.mutate, isPending: create.isPending, error: create.error },
    update: { mutate: update.mutate, isPending: update.isPending, error: update.error },
    delete: { mutate: del.mutate, isPending: del.isPending, error: del.error },
  };
}

function EffectDefForm({ effect, onSubmit, onCancel, isLoading, error }: {
  effect: EffectDef | null
  onSubmit: (input: EffectDefInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string | null
}) {
  const isEdit = effect !== null;
  const [form, setForm] = useState<EffectDefInput>(() =>
    effect
      ? { ...effect, parameters: { ...effect.parameters }, messages: { ...effect.messages } }
      : createEmptyInput("hp_change"),
  );
  const [showAdvanced, setShowAdvanced] = useState(false);

  const set = (patch: Partial<EffectDefInput>) =>
    setForm((prev) => ({ ...prev, ...patch }));

  // Get parameter config for current effect type
  const paramConfig = getParamConfig(form.effect_type);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim() || !form.effect_type) return;
    onSubmit(form);
  };

  return (
    <form onSubmit={handleSubmit} className="bg-surface rounded-lg border border-border p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-text">{isEdit ? "Edit Effect" : "Create Effect"}</h3>
        <button
          type="button"
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="text-sm text-primary hover:underline"
        >
          {showAdvanced ? "Hide Advanced" : "Show Advanced"}
        </button>
      </div>

      {error && (
        <div className="bg-danger/10 border border-danger rounded p-3 text-danger text-sm">
          {error}
        </div>
      )}

      {/* Basic Info Section */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="md:col-span-2">
          <FormField label="Name" value={form.name} onChange={(v) => set({ name: v })} required />
        </div>
        <div className="md:col-span-2">
          <SelectField
            label="Effect Type"
            value={form.effect_type}
            onChange={(v) => set({ effect_type: v, parameters: {} })}
            options={EFFECT_TYPES.map((t) => ({ value: t.value, label: t.label }))}
            required
          />
        </div>
        <div className="md:col-span-2">
          <TextareaField
            label="Description"
            value={form.description}
            onChange={(v) => set({ description: v })}
          />
        </div>
      </div>

      {/* Parameters Section */}
      <div className="border-t border-border pt-4">
        <h4 className="text-sm font-semibold text-text mb-3">Parameters</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {paramConfig.map((pc) => {
            if (pc.type === "number") {
              return (
                <div key={pc.key}>
                  <FieldLabel htmlFor={pc.key}>{pc.label}</FieldLabel>
                  <NumberField
                    id={pc.key}
                    label=""
                    value={Number(form.parameters[pc.key] ?? 0)}
                    onChange={(v) => set({ parameters: { ...form.parameters, [pc.key]: v } })}
                  />
                </div>
              );
            }
            if (pc.type === "select") {
              return (
                <SelectField
                  key={pc.key}
                  label={pc.label}
                  value={String(form.parameters[pc.key] ?? "")}
                  onChange={(v) => set({ parameters: { ...form.parameters, [pc.key]: v } })}
                  options={pc.options ?? []}
                />
              );
            }
            return (
              <FormField
                key={pc.key}
                label={pc.label}
                value={String(form.parameters[pc.key] ?? "")}
                onChange={(v) => set({ parameters: { ...form.parameters, [pc.key]: v } })}
              />
            );
          })}
        </div>
      </div>

      {/* Advanced Section */}
      {showAdvanced && (
        <div className="space-y-4 border-t border-border pt-4">
          {/* Stack Settings */}
          <div>
            <h4 className="text-sm font-semibold text-text mb-3">Stack Settings</h4>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <SelectField
                label="Stack Mode"
                value={form.stack_mode}
                onChange={(v) => set({ stack_mode: v })}
                options={STACK_MODES}
              />
              <NumberField
                label="Stack Limit"
                value={form.stack_limit}
                onChange={(v) => set({ stack_limit: v })}
              />
              <CheckboxField
                label="Permanent"
                checked={form.is_permanent}
                onChange={(v) => set({ is_permanent: v })}
              />
            </div>
          </div>

          {/* Duration */}
          {!form.is_permanent && (
            <NumberField
              label="Duration (seconds)"
              value={form.duration_secs}
              onChange={(v) => set({ duration_secs: v })}
            />
          )}

          {/* Messages */}
          <div>
            <h4 className="text-sm font-semibold text-text mb-3">Messages (Optional)</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField
                label="On Start"
                value={form.messages.on_start ?? ""}
                onChange={(v) => set({ messages: { ...form.messages, on_start: v } })}
              />
              <FormField
                label="On End"
                value={form.messages.on_end ?? ""}
                onChange={(v) => set({ messages: { ...form.messages, on_end: v } })}
              />
            </div>
          </div>
        </div>
      )}

      <div className="flex gap-3 pt-4 border-t border-border">
        <Button variant="primary" type="submit" disabled={isLoading}>
          {isLoading ? "Saving…" : isEdit ? "Update Effect" : "Create Effect"}
        </Button>
        <Button variant="secondary" type="button" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </form>
  );
}

function EffectsManagement() {
  const { data: effects = [] } = useEffectDefs();
  const mutations = useEffectDefMutations();
  const [showForm, setShowForm] = useState(false);
  const [editingEffect, setEditingEffect] = useState<EffectDef | null>(null);

  const handleCreate = (input: EffectDefInput) => {
    mutations.create.mutate(input, {
      onSuccess: () => { setShowForm(false); showToast("Effect created", "success"); },
    });
  };
  const handleUpdate = (input: EffectDefInput) => {
    if (!editingEffect) return;
    mutations.update.mutate({ id: editingEffect.id, input }, {
      onSuccess: () => { setEditingEffect(null); showToast("Effect updated", "success"); },
    });
  };
  const handleDelete = (id: number) => {
    mutations.delete.mutate(id, {
      onSuccess: () => { showToast("Effect deleted", "success"); },
    });
  };

  const columns: Column<EffectDef>[] = [
    { header: "Name", accessor: "name" },
    {
      header: "Type",
      accessor: "effect_type",
      render: (v) => (
        <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-accent/10 text-accent">
          {String(v)}
        </span>
      ),
    },
    {
      header: "Stack",
      accessor: "stack_mode",
      render: (v, r) => (
        <div className="flex items-center gap-2">
          <span className="text-xs">{String(v)}</span>
          {r.stack_limit > 1 && <span className="text-xs text-text-muted">×{r.stack_limit}</span>}
        </div>
      ),
    },
    {
      header: "Duration",
      accessor: "duration_secs",
      render: (v, r) => {
        const duration = Number(v);
        return r.is_permanent ? (
          <span className="text-success font-medium">Permanent</span>
        ) : duration === 0 ? (
          <span className="text-text-muted">Instant</span>
        ) : (
          <span className="text-xs">{duration}s</span>
        );
      },
    },
    {
      header: "Hooks",
      accessor: "hook_count",
      render: (v) => <span className="text-text-muted">{Number(v)}</span>,
    },
    {
      header: "",
      accessor: "_actions",
      align: "right",
      render: (_, r) => (
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setEditingEffect(r); }}>
            Edit
          </Button>
          <Button variant="danger" size="sm" onClick={(e) => { e.stopPropagation(); handleDelete(r.id); }}>
            Delete
          </Button>
        </div>
      ),
    },
  ];

  // Filter effects by group for display
  const effectsByGroup = EFFECT_TYPE_GROUPS.map((group) => ({
    group: group.label,
    effects: effects.filter((e) => {
      const type = EFFECT_TYPES.find((t) => t.value === e.effect_type);
      return type?.group === group.label;
    }),
  }));

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <PageHeader
        title="Effects"
        backTo="/dashboard"
        actions={
          <Button
            variant="primary"
            onClick={() => { setShowForm(true); setEditingEffect(null); }}
          >
            + Create Effect
          </Button>
        }
      />

      {showForm && !editingEffect && (
        <EffectDefForm
          effect={null}
          onSubmit={handleCreate}
          onCancel={() => setShowForm(false)}
          isLoading={mutations.create.isPending}
          error={mutations.create.error?.message ?? null}
        />
      )}
      {editingEffect && (
        <EffectDefForm
          effect={editingEffect}
          onSubmit={handleUpdate}
          onCancel={() => setEditingEffect(null)}
          isLoading={mutations.update.isPending}
          error={mutations.update.error?.message ?? null}
        />
      )}

      {/* Effects list grouped by type */}
      {!showForm && !editingEffect && (
        <div className="space-y-8 mt-6">
          {effectsByGroup.map((group) => (
            group.effects.length > 0 && (
              <section key={group.group}>
                <h2 className="text-xl font-semibold text-text mb-4 pb-2 border-b border-border">
                  {group.group} Effects ({group.effects.length})
                </h2>
                <DataTable<EffectDef>
                  columns={columns}
                  data={group.effects}
                  getKey={(r) => r.id}
                  emptyMessage={`No ${group.group.toLowerCase()} effects yet.`}
                />
              </section>
            )
          ))}
        </div>
      )}

      {/* Empty state */}
      {effects.length === 0 && !showForm && !editingEffect && (
        <div className="text-center py-20 bg-surface-muted rounded-lg border border-border">
          <h3 className="text-lg font-semibold text-text mb-2">No Effects Yet</h3>
          <p className="text-text-muted mb-6 max-w-md mx-auto">
            Effects are data-driven game mechanics that can be applied to NPCs, items, and abilities.
            Create your first effect to get started.
          </p>
          <Button variant="primary" onClick={() => setShowForm(true)}>
            + Create Your First Effect
          </Button>
        </div>
      )}
    </div>
  );
}

export const Route = createFileRoute("/_auth/effects")({ component: EffectsManagement });
