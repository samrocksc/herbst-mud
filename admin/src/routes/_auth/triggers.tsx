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
import { useWorlds } from "../../hooks/useWorlds";

export const Route = createFileRoute("/_auth/triggers")({
  component: TriggersIndex,
});

export function TriggersIndex() {
  const [searchQuery, setSearchQuery] = useState("");
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const navigate = useNavigate();
  const { currentWorld, setWorld } = useWorldStore();
  const { data: worlds } = useWorlds();

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
    { header: "World", accessor: "world_id", align: "center" },
    { header: "Type", accessor: "trigger_type", align: "center" },
    { header: "Target Type", accessor: "target_type", align: "center" },
    { header: "Target ID", accessor: "target_id", align: "center" },
    {
      header: "Room",
      accessor: "room_id",
      align: "center",
      render: (_, row) => (
        <span className="badge badge-neutral">{row.room_id ?? "-"}</span>
      ),
    },
    {
      header: "Enabled",
      accessor: "enabled",
      align: "center",
      render: (_, row) => (
        <span className={`badge ${row.enabled ? "badge-success" : "badge-error"}`}>
          {row.enabled ? "Yes" : "No"}
        </span>
      ),
    },
    {
      header: "",
      accessor: "_actions",
      align: "right",
      render: (_, row) => (
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setDeleteId(row.id); }}>
            Delete
          </Button>
        </div>
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

      <div className="mb-4 flex items-center gap-4">
        <select
          value={currentWorld}
          onChange={(e) => setWorld(e.target.value)}
          className="px-3 py-2 bg-surface-muted border border-border rounded text-sm focus:outline-none focus:border-primary"
        >
          {worlds?.map((w) => (
            <option key={w.id} value={String(w.id)}>
              {w.name}
            </option>
          ))}
        </select>
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
