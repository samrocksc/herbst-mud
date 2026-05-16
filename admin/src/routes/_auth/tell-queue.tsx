import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { useState } from "react";
import { useTellQueue, useDeleteTell } from "../../hooks/useTellQueue";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { showToast } from "../../components/Toast";
import type { TellQueueEntry } from "../../hooks/useTellQueue";

export const Route = createFileRoute("/_auth/tell-queue")({
  component: TellQueuePage,
});

const COLUMNS: Column<TellQueueEntry>[] = [
  { header: "ID", accessor: "id", align: "center" },
  { header: "Sender", accessor: "sender_name" },
  { header: "Recipient", accessor: "recipient_name" },
  { header: "Message", accessor: "message" },
  { header: "Sent", accessor: "sent_at", render: (v) => formatDate(String(v)) },
  { header: "Delivered", accessor: "delivered_at", render: (v) => v ? formatDate(String(v)) : "—" },
  {
    header: "Status",
    accessor: "delivered_at",
    render: (v) => v
      ? <span className="badge badge-success">Delivered</span>
      : <span className="badge badge-warning">Pending</span>,
  },
];

function formatDate(t: string): string {
  if (!t) return "—";
  try {
    const d = new Date(t);
    return d.toLocaleString([], { month: "short", day: "numeric", hour: "2-digit", minute: "2-digit" });
  } catch {
    return t;
  }
}

function TellQueuePage() {
  const [showUndelivered, setShowUndelivered] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const location = useLocation();
  const { data: entries, isLoading, error } = useTellQueue({
    undelivered: showUndelivered || undefined,
    limit: 200,
  });
  const deleteMutation = useDeleteTell();

  const filtered = (entries ?? []).filter((e) => showUndelivered ? !e.delivered_at : true);

  const handleDelete = () => {
    if (deleteId == null) return;
    deleteMutation.mutate(deleteId, {
      onSuccess: () => { setDeleteId(null); showToast("Tell deleted", "success"); },
    });
  };

  if (location.pathname !== "/tell-queue") return <Outlet />;

  if (isLoading) return <div className="loading">Loading tell queue...</div>;
  if (error) return <div className="error">Failed to load tell queue: {error.message}</div>;

  return (
    <div className="management-page">
      <PageHeader title="Offline Tell Queue" backTo="/dashboard" />

      <div className="flex items-center gap-4 mb-4">
        <label className="flex items-center gap-2 text-sm text-text">
          <input
            type="checkbox"
            checked={showUndelivered}
            onChange={(e) => setShowUndelivered(e.target.checked)}
          />
          Show only pending
        </label>
        <span className="text-xs text-text-muted">
          {filtered.length} of {entries?.length ?? 0} entries
        </span>
      </div>

      <DataTable
        columns={[
          ...COLUMNS,
          {
            header: "",
            accessor: "_actions",
            align: "right",
            render: (_, row) => (
              <div className="flex gap-2 justify-end">
                {!row.delivered_at && (
                  <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setDeleteId(row.id); }}>
                    Delete
                  </Button>
                )}
              </div>
            ),
          },
        ]}
        data={filtered}
        getKey={(row) => row.id}
        emptyMessage={showUndelivered ? "No pending tells." : "Tell queue is empty."}
      />

      <DeleteConfirmation
        open={deleteId != null}
        title="Delete Pending Tell"
        message="Are you sure you want to delete this pending tell? The message will not be delivered."
        onConfirm={handleDelete}
        onCancel={() => setDeleteId(null)}
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}