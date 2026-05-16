/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { useState } from "react";
import { useChannelConfigs, useUpdateChannel } from "../../hooks/useChannels";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { ColorField } from "../../components/fields/ColorField";
import { FormField, NumberField, CheckboxField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";
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
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
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
    </div>
  );
}

function ChannelsManagement() {
  const [editing, setEditing] = useState<string | null>(null);
  const location = useLocation();
  const { data: channels, isLoading, error } = useChannelConfigs();

  if (location.pathname !== "/channels") return <Outlet />;

  if (isLoading) return <div className="loading">Loading channels...</div>;
  if (error) return <div className="error">Failed to load channels: {error.message}</div>;

  const editingChannel = channels?.find((c) => c.name === editing);

  return (
    <div className="management-page">
      <PageHeader title="Channels" backTo="/dashboard" />

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
              <Button variant="ghost" size="sm" onClick={() => setEditing(row.name)}>
                Edit
              </Button>
            ),
          },
        ]}
        data={channels ?? []}
        getKey={(row) => row.name}
        emptyMessage="No channels configured."
      />
    </div>
  );
}