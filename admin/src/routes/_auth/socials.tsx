import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useSocials, useDeleteSocial } from "../../hooks/useSocials";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import type { SocialCommand } from "../../hooks/useSocials";

export const Route = createFileRoute("/_auth/socials")({
  component: SocialsManagement,
});

const COLUMNS: Column<SocialCommand>[] = [
  {
    header: "Name",
    accessor: "name",
    render: (_, row) => (
      <Link
        to="/socials/$socialId"
        params={{ socialId: String(row.id) }}
        className="no-underline text-primary hover:underline font-bold"
      >
        {row.name}
      </Link>
    ),
  },
  { header: "Self Text", accessor: "self_text" },
  { header: "Room Text", accessor: "room_text" },
];

function SocialsManagement() {
  const [search, setSearch] = useState("");
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const { data: socials, isLoading, error } = useSocials();
  const deleteMutation = useDeleteSocial();

  const filtered = (socials ?? []).filter((s) =>
    s.name.toLowerCase().includes(search.toLowerCase()) ||
    s.self_text.toLowerCase().includes(search.toLowerCase())
  );

  const handleDelete = () => {
    if (deleteId == null) return;
    deleteMutation.mutate(deleteId, {
      onSuccess: () => {
        setDeleteId(null);
        showToast("Social deleted", "success");
      },
    });
  };

  if (location.pathname !== "/socials") return <Outlet />;

  if (isLoading) return <div className="loading">Loading socials...</div>;
  if (error) return <div className="error">Failed to load socials: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Social Commands"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => navigate({ to: "/socials/new" })}>
            + Add Social
          </Button>
        }
      />

      <div className="mb-4">
        <input
          type="text"
          placeholder="Search socials by name or text..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      <DataTable
        columns={[
          ...COLUMNS.slice(0, -1), // Remove the actions column from COLUMNS
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
        ]}
        data={filtered}
        getKey={(row) => row.id}
        onRowClick={(row) => navigate({ to: "/socials/$socialId", params: { socialId: String(row.id) } })}
        emptyMessage={
          search ? "No socials match this search." : "No socials found. Create your first social command!"
        }
      />

      <DeleteConfirmation
        open={deleteId != null}
        title="Delete Social Command"
        message="Are you sure you want to delete this social command? This action cannot be undone."
        onConfirm={handleDelete}
        onCancel={() => setDeleteId(null)}
        isLoading={deleteMutation.isPending}
      />
    </PageContainer>
  );
}