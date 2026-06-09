import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useQuests, type Quest } from "../../hooks/useQuests";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/quests")({
  component: QuestsManagement,
});

const QUEST_TYPE_BADGES: Record<string, { label: string; color: string }> = {
  general:   { label: "General",   color: "bg-accent/20 text-accent" },
  hunter:    { label: "Hunter",    color: "bg-danger/20 text-danger" },
  collector: { label: "Collector", color: "bg-success/20 text-success" },
  explorer:  { label: "Explorer",  color: "bg-primary/20 text-primary" },
};

const REPEAT_MODE_BADGES: Record<string, string> = {
  none:     "One-time",
  cooldown: "Cooldown",
  always:   "Repeatable",
};

const OBJ_LABELS: Record<string, string> = {
  kill: "Kill", explore: "Visit", collect: "Collect",
};

function objectiveSummary(obj: Quest["objectives"][number]): string {
  const verb = OBJ_LABELS[obj.type] || obj.type;
  const count = obj.count > 1 ? ` (×${obj.count})` : "";
  return `${verb} ${obj.labels?.[0] || obj.target_id || "?"}${count}`;
}

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
  {
    header: "Type",
    accessor: "main_type",
    render: (val) => {
      const t = (val as string) || "general";
      const badge = QUEST_TYPE_BADGES[t] || QUEST_TYPE_BADGES.general;
      return <span className={`px-1.5 py-0.5 rounded text-xs font-medium ${badge.color}`}>{badge.label}</span>;
    },
  },
  {
    header: "Status",
    accessor: "is_active",
    render: (val) => val
      ? <span className="px-1.5 py-0.5 rounded text-xs font-medium bg-success/20 text-success">Active</span>
      : <span className="px-1.5 py-0.5 rounded text-xs font-medium bg-surface-muted text-text-muted">Inactive</span>,
  },
  {
    header: "Repeat",
    accessor: "repeat_mode",
    render: (val) => {
      const mode = (val as string) || "none";
      const label = REPEAT_MODE_BADGES[mode] || mode;
      return <span className="text-xs text-text-muted">{label}</span>;
    },
  },
  {
    header: "Objectives",
    accessor: "objectives",
    render: (_, row) => {
      const objs = (row.objectives ?? []) as Quest["objectives"];
      if (objs.length === 0) return <span className="text-xs text-text-muted">—</span>;
      const preview = objs.slice(0, 2).map(objectiveSummary).join(", ");
      const remainder = objs.length > 2 ? ` +${objs.length - 2} more` : "";
      return <span className="text-xs text-text" title={objs.map(objectiveSummary).join("\n")}>{preview}{remainder}</span>;
    },
  },
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
