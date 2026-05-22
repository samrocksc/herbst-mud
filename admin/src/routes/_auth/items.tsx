/* eslint-disable functional/immutable-data, functional/no-loop-statements */
import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { useEquipmentTemplates, useDeleteTemplate } from "../../hooks/useEquipmentTemplates";
import { useItemInstances } from "../../hooks/useItemInstances";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { showToast } from "../../components/Toast";
import type { EquipmentTemplate } from "../../hooks/useEquipmentTemplates";

export const Route = createFileRoute("/_auth/items")({
  component: ItemsIndex,
});

function ItemsIndex() {
  const [searchQuery, setSearchQuery] = useState("");
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const navigate = useNavigate();

  const templatesQuery = useEquipmentTemplates();
  const instancesQuery = useItemInstances();
  const deleteMutation = useDeleteTemplate();

  const instanceCounts = useMemo(() => {
    const counts: Record<number, number> = {};
    for (const inst of instancesQuery.data ?? []) {
      const tid = inst.equipment_template_id;
      if (tid) counts[tid] = (counts[tid] ?? 0) + 1;
    }
    return counts;
  }, [instancesQuery.data]);

  const filteredItems = (templatesQuery.data ?? []).filter((item) =>
    item.name.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => { setDeleteId(null); showToast("Item template deleted", "success"); },
    });
  };

  const columns: Column<EquipmentTemplate>[] = [
    {
      header: "Name",
      accessor: "name",
      render: (_, row) => (
        <Link
          to="/items/$itemId"
          params={{ itemId: row.id }}
          className="no-underline text-primary hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    { header: "Slot", accessor: "slot" },
    { header: "Level", accessor: "level", align: "center" },
    { header: "Type", accessor: "item_type" },
    { header: "Weight", accessor: "weight", align: "center" },
    {
      header: "Instances",
      accessor: "instances",
      align: "center",
      render: (_, row) => (
        <span className="badge badge-neutral">{instanceCounts[row.id] ?? 0}</span>
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
  const isList = location.pathname === "/items";

  if (!isList) return <Outlet />;

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title="Items" showBack backTo="/dashboard" actions={
        <Button variant="primary" onClick={() => navigate({ to: "/items/new" })}>+ Add Item</Button>
      } />

      <div className="mb-4">
        <input
          type="text"
          placeholder="Search items by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      {templatesQuery.isLoading && <div className="p-8 text-text-muted text-center text-xs">Loading items...</div>}
      {templatesQuery.isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load items: {templatesQuery.error?.message ?? "Unknown error"}
        </div>
      )}
      {templatesQuery.isSuccess && (
        <DataTable<EquipmentTemplate>
          columns={columns}
          data={filteredItems}
          getKey={(row) => row.id}
          emptyMessage="No items found."
          variant="dark"
        />
      )}

      {deleteId && (
        <DeleteConfirmation
          open={!!deleteId}
          title="Delete Item Template"
          message="Are you sure? This will permanently delete this item template. Instances based on this template will not be deleted."
          onConfirm={() => handleDelete(deleteId)}
          onCancel={() => setDeleteId(null)}
          isLoading={deleteMutation.isPending}
        />
      )}
    </div>
  );
}

