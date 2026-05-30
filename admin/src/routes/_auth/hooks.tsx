/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useHooks, useUpdateHook, useDeleteHook, type EffectHook, type HookInput } from "../../hooks/useHooks";
import { useEffectDefs } from "../../hooks/useEffectDefs";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { showToast } from "../../components/Toast";
import { FormField, SelectField, CheckboxField } from "../../components/fields";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/hooks")({
  component: HooksManagement,
});

const HOOK_EVENTS = [
  { value: "on_death", label: "On Death" },
  { value: "on_hit_received", label: "On Hit Received" },
  { value: "on_hit_dealt", label: "On Hit Dealt" },
  { value: "on_kill", label: "On Kill" },
  { value: "on_enter_room", label: "On Enter Room" },
  { value: "on_leave_room", label: "On Leave Room" },
  { value: "on_equip", label: "On Equip" },
  { value: "on_unequip", label: "On Unequip" },
  { value: "on_login", label: "On Login" },
  { value: "on_effect_start", label: "On Effect Start" },
  { value: "on_effect_end", label: "On Effect End" },
];

const HOOK_TARGETS = [
  { value: "self", label: "Self (the character)" },
  { value: "attacker", label: "Attacker (hit dealer)" },
  { value: "killer", label: "Killer (death dealer)" },
  { value: "room", label: "Room (all in room)" },
  { value: "owner", label: "Owner (item/NPC owner)" },
];

function HooksManagement() {
  const { data: hooks, isLoading, error } = useHooks();
  const updateHook = useUpdateHook();
  const deleteHook = useDeleteHook();
  const { data: effectDefs } = useEffectDefs();
  const [editing, setEditing] = useState<number | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<HookInput | null>(null);

  const editingHook = hooks?.find((h) => h.id === editing);

  const startEdit = (hook: EffectHook) => {
    setEditing(hook.id);
    setEditForm({
      name: hook.name,
      event: hook.event,
      target: hook.target,
      condition: hook.condition ?? "",
      enabled: hook.enabled,
      effect_id: hook.effect_id,
    });
  };

  const handleSave = async () => {
    if (editing == null || !editForm) return;
    try {
      await updateHook.mutateAsync({ id: editing, input: editForm });
      showToast("Hook updated", "success");
      setEditing(null);
      setEditForm(null);
    } catch { /* toasted globally */ }
  };

  const handleDelete = async () => {
    if (deleting == null) return;
    try {
      await deleteHook.mutateAsync(deleting);
      showToast("Hook deleted", "success");
      setDeleting(null);
    } catch { /* toasted globally */ }
  };

  const effectOpts = [
    { value: "", label: "— None —" },
    ...(effectDefs ?? []).map((e) => ({ value: String(e.id), label: e.name })),
  ];

  const columns: Column<EffectHook>[] = [
    { header: "ID", accessor: "id", align: "center" },
    { header: "Name", accessor: "name" },
    { header: "Event", accessor: "event" },
    { header: "Target", accessor: "target" },
    { header: "Condition", accessor: "condition", render: (v) => String(v || "—") },
    { header: "Active", accessor: "enabled", render: (v) => v ? "Yes" : "No" },
    {
      header: "Effect",
      accessor: "effect_name",
      render: (v, row) => row.effect_id ? `${v || `#${row.effect_id}`}` : "—",
    },
    {
      header: "NPC Template",
      accessor: "npc_template_name",
      render: (v, row) => row.npc_template_id ? (v || `#${row.npc_template_id}`) : "—",
    },
    {
      header: "",
      accessor: "_actions",
      align: "right",
      render: (_, row) => (
        <div className="flex gap-1 justify-end">
          <Button variant="ghost" size="sm" onClick={() => startEdit(row)}>Edit</Button>
          <Button variant="danger" size="sm" onClick={() => setDeleting(row.id)}>Delete</Button>
        </div>
      ),
    },
  ];

  if (isLoading) return <div className="loading">Loading hooks...</div>;
  if (error) return <div className="error">Failed to load hooks: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader title="Hooks" backTo="/dashboard" />

      {editingHook && editForm && (
        <div className="bg-surface p-6 border border-border rounded mb-4">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Hook: {editingHook.name}</h3>
          <p className="text-text-muted text-xs mb-4">
            NPC Template: {editingHook.npc_template_name || editingHook.npc_template_id}
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name" value={editForm.name} onChange={(v) => setEditForm({ ...editForm, name: v })} />
            <SelectField label="Event" value={editForm.event} onChange={(v) => setEditForm({ ...editForm, event: v })} options={HOOK_EVENTS} />
            <SelectField label="Target" value={editForm.target} onChange={(v) => setEditForm({ ...editForm, target: v })} options={HOOK_TARGETS} />
            <FormField label="Condition" value={editForm.condition ?? ""} onChange={(v) => setEditForm({ ...editForm, condition: v })} placeholder="Optional condition expression" />
            <SelectField label="Effect" value={String(editForm.effect_id)} onChange={(v) => setEditForm({ ...editForm, effect_id: Number(v) || 0 })} options={effectOpts} />
            <CheckboxField label="Enabled" checked={editForm.enabled} onChange={(v) => setEditForm({ ...editForm, enabled: v })} />
          </div>
          <div className="flex gap-2 mt-4">
            <Button variant="primary" onClick={handleSave} disabled={updateHook.isPending}>
              {updateHook.isPending ? "Saving..." : "Save Changes"}
            </Button>
            <Button variant="secondary" onClick={() => { setEditing(null); setEditForm(null); }}>Cancel</Button>
          </div>
        </div>
      )}

      <DataTable
        columns={columns}
        data={hooks ?? []}
        getKey={(row) => row.id}
        emptyMessage="No hooks configured."
        variant="dark"
      />

      <DeleteConfirmation
        open={deleting !== null}
        title="Delete Hook"
        message={`Are you sure you want to delete this hook? This cannot be undone.`}
        onConfirm={handleDelete}
        onCancel={() => setDeleting(null)}
        isLoading={deleteHook.isPending}
      />
    </PageContainer>
  );
}
