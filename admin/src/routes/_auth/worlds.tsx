/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useWorlds, useSetActiveWorld } from "../../hooks/useWorlds";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import type { World } from "../../hooks/useWorlds";

export const Route = createFileRoute("/_auth/worlds")({
  component: WorldsManagement,
});

const COLUMNS: Column<World>[] = [
  {
    header: "ID",
    accessor: "id",
    render: (_, row) => <span className="font-mono text-xs">{row.id}</span>,
  },
  {
    header: "Name",
    accessor: "name",
    render: (_, row) => (
      <Link
        to="/worlds/$worldId"
        params={{ worldId: String(row.id) }}
        className="no-underline text-primary hover:underline font-bold"
      >
        {row.name}
      </Link>
    ),
  },
  { header: "Title", accessor: "title" },
  { header: "Description", accessor: "description" },
  {
    header: "Active",
    accessor: "active",
    render: (_: unknown, row: World) => <ActiveToggle world={row} />,
  },
];

function ActiveToggle({ world }: { world: World }) {
  const setActive = useSetActiveWorld();
  const isChecked = world.active ?? false;

  const handleToggle = () => {
    setActive.mutate({ id: world.id, active: !isChecked });
  };

  return (
    <label className="flex items-center cursor-pointer gap-2">
      <input
        type="checkbox"
        checked={isChecked}
        onChange={handleToggle}
        disabled={setActive.isPending}
        className="accent-primary"
      />
      <span className="text-sm">{isChecked ? "Active" : "Inactive"}</span>
    </label>
  );
}

export function WorldsManagement() {
  const navigate = useNavigate();
  const location = useLocation();
  const { data: worlds, isLoading, error } = useWorlds();

  if (location.pathname !== "/worlds") return <Outlet />;

  if (isLoading) return <div className="loading">Loading worlds...</div>;
  if (error) return <div className="error">Failed to load worlds: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Worlds"
        showBack
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => navigate({ to: "/worlds/new" })}>
            + Add World
          </Button>
        }
      />

      <DataTable
        columns={COLUMNS}
        data={worlds ?? []}
        getKey={(row) => row.id}
        emptyMessage="No worlds found. Create your first world!"
      />
    </PageContainer>
  );
}