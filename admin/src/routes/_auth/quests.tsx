import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useQuests, type Quest } from "../../hooks/useQuests";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/quests")({
  component: QuestsManagement,
});

const COLUMNS: Column<Quest>[] = [
  {
    header: "Name",
    accessor: "name",
    render: (_, row) => (
      <Link to="/quests/$questId" params={{ questId: String(row.id) }} className="no-underline text-primary hover:underline font-bold">
        {row.name}
      </Link>
    ),
  },
  { header: "Active", accessor: "is_active", render: (val) => val ? "✓" : "✗" },
  { header: "Repeat", accessor: "repeat_mode" },
  { header: "Objectives", accessor: "objectives", render: (val) => (val as unknown as unknown[])?.length ?? 0 },
  { header: "XP", accessor: "rewards", render: (val) => (val as { xp?: number })?.xp ?? 0 },
];

export function QuestsManagement() {
  const navigate = useNavigate();
  const location = useLocation();
  const { data: quests, isLoading, error } = useQuests();

  if (location.pathname !== "/quests") return <Outlet />;

  if (isLoading) return <div className="loading">Loading quests...</div>;
  if (error) return <div className="error">Failed to load quests: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Quests"
        backTo="/dashboard"
        actions={
          <Link to="/quests/new">
            <Button variant="primary" size="sm">+ Add Quest</Button>
          </Link>
        }
      />
      <DataTable
        columns={COLUMNS}
        data={quests ?? []}
        getKey={(row) => row.id}
        onRowClick={(row) => navigate({ to: "/quests/$questId", params: { questId: String(row.id) } })}
        emptyMessage="No quests found. Create your first quest!"
      />
    </PageContainer>
  );
}
