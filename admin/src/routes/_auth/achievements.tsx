/* eslint-disable functional/immutable-data */
import { createFileRoute } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { useAchievements, type Achievement } from "../../hooks/useAchievements";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { FilterBar } from "../../components/FilterBar";

export const Route = createFileRoute("/_auth/achievements")({
  component: AchievementsPage,
});

type CompletionFilter = "all" | "completed" | "incomplete";

function AchievementsPage() {
  const { data: achievements, isLoading, error } = useAchievements();
  const [completionFilter, setCompletionFilter] = useState<CompletionFilter>("all");
  const [search, setSearch] = useState("");

  const allAchievements = achievements ?? [];

  // Derived: filter by completion and search text
  const filtered = useMemo(() => {
    let result = allAchievements;

    if (completionFilter === "completed") {
      result = result.filter((a) => getCompletedCount(a) > 0);
    } else if (completionFilter === "incomplete") {
      result = result.filter((a) => getCompletedCount(a) === 0);
    }

    if (search.trim()) {
      const q = search.toLowerCase();
      result = result.filter((a) =>
        a.name.toLowerCase().includes(q) ||
        a.description.toLowerCase().includes(q) ||
        a.criteria.toLowerCase().includes(q)
      );
    }

    return result;
  }, [allAchievements, completionFilter, search]);

  const columns: Column<Achievement>[] = [
    {
      header: "Icon",
      accessor: "icon",
      render: (val: unknown) => (
        <span className="text-xl text-center block" title={String(val ?? "")}>
          {String(val ?? "🏆")}
        </span>
      ),
      className: "w-12 text-center",
      align: "center",
    },
    {
      header: "Name",
      accessor: "name",
      render: (val: unknown) => (
        <span className="font-bold text-primary">{String(val ?? "")}</span>
      ),
    },
    {
      header: "Description",
      accessor: "description",
      render: (val: unknown) => (
        <span className="text-text-muted text-xs" title={String(val ?? "")}>
          {String(val ?? "").slice(0, 100)}
          {String(val ?? "").length > 100 ? "…" : ""}
        </span>
      ),
    },
    {
      header: "Criteria",
      accessor: "criteria",
      render: (val: unknown) => (
        <code className="text-xs text-text-muted/80 bg-surface-muted px-1.5 py-0.5 rounded">
          {String(val ?? "")}
        </code>
      ),
    },
    {
      header: "XP Reward",
      accessor: "xp_reward",
      render: (val: unknown) => (
        <span className="text-sm text-text font-medium">
          {val != null ? `${val} XP` : "—"}
        </span>
      ),
      align: "center",
    },
    {
      header: "Completions",
      accessor: "completed_count",
      render: (_val: unknown, row: Achievement) => {
        const count = getCompletedCount(row);
        const badge = count > 0
          ? "bg-emerald-900/40 text-emerald-300 border-emerald-600/30"
          : "bg-slate-700/50 text-slate-400 border-slate-600/30";
        return (
          <span className={`text-xs px-2 py-0.5 rounded border font-medium ${badge}`}>
            {count > 0 ? `${count} completed` : "None yet"}
          </span>
        );
      },
      align: "center",
    },
  ];

  if (isLoading) return <div className="loading p-8 text-center text-text-muted">Loading achievements...</div>;
  if (error) return <div className="error p-8 text-center text-danger">Failed to load achievements: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Achievements"
        backTo="/dashboard"
      />

      <p className="text-sm text-muted mb-4">
        View all game achievements, their completion criteria, and how many characters have completed each one.
      </p>

      <FilterBar
        showClear={completionFilter !== "all" || !!search.trim()}
        onClear={() => { setCompletionFilter("all"); setSearch(""); }}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Status:</label>
          <select
            value={completionFilter}
            onChange={(e) => setCompletionFilter(e.target.value as CompletionFilter)}
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          >
            <option value="all">All Achievements</option>
            <option value="completed">Completed Only</option>
            <option value="incomplete">Incomplete Only</option>
          </select>
        </div>

        <div className="flex flex-col gap-1 flex-1 min-w-[200px]">
          <label className="text-xs text-text-muted">Search:</label>
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search name, description, or criteria..."
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          />
        </div>
      </FilterBar>

      {/* Summary stats */}
      <div className="flex gap-4 mb-4 text-sm">
        <div className="bg-surface-muted rounded-lg px-4 py-2 border border-border">
          <span className="text-text-muted">Total: </span>
          <span className="text-text font-bold">{allAchievements.length}</span>
        </div>
        <div className="bg-surface-muted rounded-lg px-4 py-2 border border-border">
          <span className="text-text-muted">Completed: </span>
          <span className="text-emerald-400 font-bold">
            {allAchievements.filter((a) => getCompletedCount(a) > 0).length}
          </span>
        </div>
        <div className="bg-surface-muted rounded-lg px-4 py-2 border border-border">
          <span className="text-text-muted">Incomplete: </span>
          <span className="text-amber-400 font-bold">
            {allAchievements.filter((a) => getCompletedCount(a) === 0).length}
          </span>
        </div>
      </div>

      <DataTable
        columns={columns}
        data={filtered}
        getKey={(row) => row.id}
        emptyMessage={
          completionFilter !== "all" || search.trim()
            ? "No achievements match your current filters."
            : "No achievements found."
        }
      />
    </PageContainer>
  );
}

/**
 * Extract a completed count from an achievement.
 * The API may return this as a top-level field or nested in metadata.
 */
function getCompletedCount(achievement: Achievement): number {
  const raw = (achievement as Record<string, unknown>);
  if (typeof raw["completed_count"] === "number") return raw["completed_count"];
  if (typeof raw["completions"] === "number") return raw["completions"];
  return 0;
}