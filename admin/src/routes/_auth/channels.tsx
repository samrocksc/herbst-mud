/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { useState } from "react";
import { useChannelConfigs, useUpdateChannel, useCreateChannel, useDeleteChannel } from "../../hooks/useChannels";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { Modal } from "../../components/Modal";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { ColorField } from "../../components/fields/ColorField";
import { FormField, NumberField, CheckboxField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import type { ChannelConfig, ChannelInput } from "../../hooks/useChannels";

export const Route = createFileRoute("/_auth/channels")({
  component: ChannelsManagement,
});

const COLUMNS: Column<ChannelConfig>[] = [
  { header: "Name", accessor: "name", render: (v) => <span className="font-semibold text-primary">[{String(v)}]</span> },
  { header: "Description", accessor: "description" },
  { header: "Color", accessor: "color", render: (v) => (
    <span className="inline-flex items-center gap-2">
      <span className="w-4 h-4 rounded border border-border" style={{ backgroundColor: String(v) }} />
      {String(v)}
    </span>
  ) },
  { header: "Default", accessor: "default_enabled", render: (v) => v ? "Enabled" : "Disabled" },
  { header: "Cooldown", accessor: "cooldown_seconds", render: (v) => `${v}s` },
  { header: "Admin Only", accessor: "admin_only", render: (v) => v ? "Yes" : "No" },
];

function ChannelEditForm({ channel, onDone }: { channel: ChannelConfig; onDone: () => void }) {
  const updateChannel = useUpdateChannel();
  const [formData, setFormData] = useState<ChannelInput>({
    name: channel.name,
    description: channel.description,
    color: channel.color,
    default_enabled: channel.default_enabled,
    cooldown_seconds: channel.cooldown_seconds,
    admin_only: channel.admin_only,
  });

  const set = (patch: Partial<ChannelInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateChannel.mutateAsync({ name: channel.name, input: formData });
      showToast("Channel updated", "success");
      onDone();
    } catch { /* toasted globally */ }
  };

  return (
    <PageContainer>
      <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Channel: {channel.name}</h3>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="Description" value={formData.description} onChange={(v) => set({ description: v })} />
          <ColorField label="Color" value={formData.color} onChange={(v) => set({ color: v })} />
        </div>
        <div className="grid grid-cols-3 gap-4">
          <CheckboxField label="Enabled by Default" checked={formData.default_enabled} onChange={(v) => set({ default_enabled: v })} />
          <CheckboxField label="Admin Only" checked={formData.admin_only} onChange={(v) => set({ admin_only: v })} />
          <NumberField label="Cooldown (seconds)" value={formData.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
        </div>
        <div className="flex gap-2">
          <Button type="submit" variant="primary" disabled={updateChannel.isPending}>
            {updateChannel.isPending ? "Saving..." : "Save Changes"}
          </Button>
          <Button variant="secondary" onClick={onDone} type="button">Cancel</Button>
        </div>
      </form>
    </PageContainer>
  );
}

function ChannelsManagement() {
  const [editing, setEditing] = useState<string | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
  const location = useLocation();
  const { data: channels, isLoading, error } = useChannelConfigs();
  const createChannel = useCreateChannel();
  const deleteChannel = useDeleteChannel();
  const [createForm, setCreateForm] = useState<ChannelInput>({
    name: "", description: "", color: "#ffffff", default_enabled: true, cooldown_seconds: 3, admin_only: false,
  });

  if (location.pathname !== "/channels") return <Outlet />;

  if (isLoading) return <div className="loading">Loading channels...</div>;
  if (error) return <div className="error">Failed to load channels: {error.message}</div>;

  const editingChannel = channels?.find((c) => c.name === editing);

  const handleCreate = async () => {
    if (!createForm.name.trim()) return;
    try {
      await createChannel.mutateAsync(createForm);
      showToast("Channel created", "success");
      setShowCreate(false);
      setCreateForm({ name: "", description: "", color: "#ffffff", default_enabled: true, cooldown_seconds: 3, admin_only: false });
    } catch { /* toasted globally */ }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await deleteChannel.mutateAsync(deleteTarget);
      showToast("Channel deleted", "success");
      setDeleteTarget(null);
    } catch { /* toasted globally */ }
  };

  return (
    <PageContainer>
      <PageHeader
        title="Channels"
        backTo="/dashboard"
        actions={
          <Button variant="primary" size="sm" onClick={() => setShowCreate(true)}>
            + Add Channel
          </Button>
        }
      />

      {showCreate && (
        <div className="bg-surface p-6 border border-border rounded mb-4">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Create Channel</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name" value={createForm.name} onChange={(v) => setCreateForm({ ...createForm, name: v })} placeholder="e.g. newbie" />
            <FormField label="Description" value={createForm.description} onChange={(v) => setCreateForm({ ...createForm, description: v })} />
            <ColorField label="Color" value={createForm.color} onChange={(v) => setCreateForm({ ...createForm, color: v })} />
            <NumberField label="Cooldown (seconds)" value={createForm.cooldown_seconds} onChange={(v) => setCreateForm({ ...createForm, cooldown_seconds: v })} />
          </div>
          <div className="flex gap-4 mt-2">
            <CheckboxField label="Enabled by Default" checked={createForm.default_enabled} onChange={(v) => setCreateForm({ ...createForm, default_enabled: v })} />
            <CheckboxField label="Admin Only" checked={createForm.admin_only} onChange={(v) => setCreateForm({ ...createForm, admin_only: v })} />
          </div>
          <div className="flex gap-2 mt-4">
            <Button variant="primary" onClick={handleCreate} disabled={createChannel.isPending}>
              {createChannel.isPending ? "Creating..." : "Create"}
            </Button>
            <Button variant="secondary" onClick={() => setShowCreate(false)}>Cancel</Button>
          </div>
        </div>
      )}

      {editingChannel && (
        <ChannelEditForm channel={editingChannel} onDone={() => setEditing(null)} />
      )}

      <DataTable
        columns={[
          ...COLUMNS,
          {
            header: "",
            accessor: "_actions",
            align: "right",
            render: (_, row) => (
              <div className="flex gap-1 justify-end">
                <Button variant="ghost" size="sm" onClick={() => setEditing(row.name)}>
                  Edit
                </Button>
                <Button variant="danger" size="sm" onClick={() => setDeleteTarget(row.name)}>
                  Delete
                </Button>
              </div>
            ),
          },
        ]}
        data={channels ?? []}
        getKey={(row) => row.name}
        emptyMessage="No channels configured."
      />

      <DeleteConfirmation
        open={deleteTarget !== null}
        title="Delete Channel"
        message={`Are you sure you want to delete channel "${deleteTarget}"? This cannot be undone.`}
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
        isLoading={deleteChannel.isPending}
      />
    </PageContainer>
  );
}