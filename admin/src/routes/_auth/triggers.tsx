/* eslint-disable functional/immutable-data, functional/no-loop-statements */
import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTriggers, useDeleteTrigger } from "../../hooks/useTriggers";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import type { Trigger } from "../../hooks/useTriggers";
import { useWorldStore } from "../../contexts/WorldStoreContext";

export const Route = createFileRoute("/_auth/triggers")({
  component: TriggersIndex,
});

const TRIGGER_BADGES: Record<string, string> = {
  use:     "bg-primary/20 text-primary",
  touch:   "bg-accent/20 text-accent",
  press:   "bg-warning/20 text-warning",
  enter:   "bg-success/20 text-success",
  examine: "bg-info/20 text-info",
  talk:    "bg-purple-500/20 text-purple-400",
};

const TARGET_BADGES: Record<string, string> = {
  recipe:      "bg-amber-500/20 text-amber-400",
  effect:      "bg-cyan-500/20 text-cyan-400",
  dialog_node: "bg-pink-500/20 text-pink-400",
};

export function TriggersIndex() {
  const [searchQuery, setSearchQuery] = useState("");
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const navigate = useNavigate();
  const { currentWorld } = useWorldStore();

  const triggersQuery = useTriggers({ world_id: currentWorld });
  const deleteMutation = useDeleteTrigger();

  const filteredTriggers = (triggersQuery.data ?? []).filter((trigger) =>
    trigger.name.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => { setDeleteId(null); showToast("Trigger deleted", "success"); },
    });
  };

  const columns: Column<Trigger>[] = [
    {
      header: "Name",
      accessor: "name",
      render: (_, row) => (
        <Link
          to="/triggers/$triggerId"
          params={{ triggerId: row.id }}
          className="no-underline text-primary hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    {
      header: "Type",
      accessor: "trigger_type",
      render: (_, row) => {
        const badge = TRIGGER_BADGES[row.trigger_type] || "bg-surface-muted text-text-muted";
        return <span className={`px-1.5 py-0.5 rounded text-xs font-medium ${badge}`}>{row.trigger_type}</span>;
      },
    },
    {
      header: "Target",
      accessor: "target_type",
      render: (_, row) => (
        <span className="text-xs text-text-muted">
          <span className={`px-1 py-0.5 rounded text-xs font-medium mr-1 ${TARGET_BADGES[row.target_type] || "bg-surface-muted text-text-muted"}`}>
            {row.target_type}
          </span>
          {row.target_id ? <code className="text-xs text-text">{row.target_id}</code> : <span className="text-xs text-text-muted">—</span>}
        </span>
      ),
    },
    {
      header: "Room",
      accessor: "room_id",
      render: (_, row) => row.room_id != null ? <code className="text-xs text-text">{row.room_id}</code> : <span className="text-xs text-text-muted">—</span>,
    },
    {
      header: "Status",
      accessor: "enabled",
      render: (_, row) => row.enabled
        ? <span className="px-1.5 py-0.5 rounded text-xs font-medium bg-success/20 text-success">Active</span>
        : <span className="px-1.5 py-0.5 rounded text-xs font-medium bg-surface-muted text-text-muted">Off</span>,
    },
    {
      header: "",
      accessor: "_actions",
      render: (_, row) => (
        <Button variant="danger" size="sm" onClick={(e) => { e.stopPropagation(); setDeleteId(row.id); }}>
          Delete
        </Button>
      ),
    },
  ];

  const location = useLocation();
  const isList = location.pathname === "/triggers";

  if (!isList) return <Outlet />;

  return (
    <PageContainer>
      <PageHeader title="Triggers" showBack backTo="/dashboard" actions={
        <Button variant="primary" onClick={() => navigate({ to: "/triggers/new" })}>+ Add Trigger</Button>
      } />

      <div className="mb-4 text-sm text-text-muted">
        Managing triggers for world <span className="text-text font-medium">{currentWorld}</span>. Switch worlds via the dashboard.
      </div>

      <div className="mb-4">
        <input
          type="text"
          placeholder="Search triggers by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      {triggersQuery.isLoading && <div className="p-8 text-text-muted text-center text-xs">Loading triggers...</div>}
      {triggersQuery.isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load triggers: {triggersQuery.error?.message ?? "Unknown error"}
        </div>
      )}
      {triggersQuery.isSuccess && (
        <DataTable<Trigger>
          columns={columns}
          data={filteredTriggers}
          getKey={(row) => row.id}
          emptyMessage="No triggers found."
          variant="dark"
        />
      )}

      {deleteId && (
        <DeleteConfirmation
          open={!!deleteId}
          title="Delete Trigger"
          message="Are you sure? This will permanently delete this trigger."
          onConfirm={() => handleDelete(deleteId)}
          onCancel={() => setDeleteId(null)}
          isLoading={deleteMutation.isPending}
        />
      )}
    </PageContainer>
  );
}
