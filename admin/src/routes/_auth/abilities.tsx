/* eslint-disable functional/immutable-data */
import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useAbilities } from "../../hooks/useAbilities";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { FilterBar } from "../../components/FilterBar";
import type { Ability } from "../../hooks/useAbilities";

export const Route = createFileRoute("/_auth/abilities")({
  component: AbilitiesManagement,
});

const COLUMNS: Column<Ability>[] = [
  {
    header: "Name",
    accessor: "name",
    render: (_, row) => (
      <Link
        to="/abilities/$abilityId"
        params={{ abilityId: String(row.id) }}
        className="no-underline text-primary hover:underline font-bold"
      >
        {row.name}
      </Link>
    ),
  },
  {
    header: "Description",
    accessor: "description",
    // Refinement #11: show a 1-line truncated description instead of the full text
    render: (val: unknown) => (
      <span className="text-text-muted text-xs" title={String(val ?? "")}>
        {String(val ?? "").slice(0, 80)}
        {String(val ?? "").length > 80 ? "…" : ""}
      </span>
    ),
  },
  {
    header: "Type",
    accessor: "ability_type",
    // Refinement #11: capitalize raw type values
    render: (val: unknown) => (
      <span className="text-xs px-2 py-0.5 rounded bg-primary/15 text-text border border-primary/30 capitalize">
        {String(val)}
      </span>
    ),
  },
  {
    header: "Class",
    accessor: "ability_class",
    render: (val: unknown) => (
      <span className="text-xs px-2 py-0.5 rounded bg-accent/15 text-text border border-accent/30 capitalize">
        {String(val)}
      </span>
    ),
  },
  {
    header: "Costs",
    accessor: "mana_cost",
    render: (_, row: Ability) => {
      const parts: React.ReactNode[] = [];
      if (row.mana_cost > 0) parts.push(<span key="mp" className="cost-badge" title="Mana Cost">MP: {row.mana_cost}</span>);
      if (row.stamina_cost > 0) parts.push(<span key="sp" className="cost-badge" title="Stamina Cost">SP: {row.stamina_cost}</span>);
      if (row.hp_cost > 0) parts.push(<span key="hp" className="cost-badge" title="HP Cost">HP: {row.hp_cost}</span>);
      if (row.cooldown_seconds > 0) parts.push(<span key="cd" className="cost-badge" title="Cooldown">CD: {row.cooldown_seconds}s</span>);
      return parts.length > 0 ? parts : <span className="text-muted">—</span>;
    },
  },
];

export function AbilitiesManagement() {
  const [filterType, setFilterType] = useState<string>("");
  const [filterClass, setFilterClass] = useState<string>("");
  const navigate = useNavigate();
  const location = useLocation();
  const { data: abilities, isLoading, error } = useAbilities({
    type: filterType || undefined,
    abilityClass: filterClass || undefined,
  });

  if (location.pathname !== "/abilities") return <Outlet />;

  if (isLoading) return <div className="loading">Loading abilities...</div>;
  if (error) return <div className="error">Failed to load abilities: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Abilities"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => navigate({ to: "/abilities/new" })}>
            + Add Ability
          </Button>
        }
      />

      <FilterBar
        showClear={!!(filterType || filterClass)}
        onClear={() => { setFilterType(""); setFilterClass(""); }}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Type:</label>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          >
            <option value="">All Types</option>
            <option value="combat">Combat</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
            <option value="healing">Healing</option>
            <option value="support">Support</option>
            <option value="defensive">Defensive</option>
          </select>
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Class:</label>
          <select
            value={filterClass}
            onChange={(e) => setFilterClass(e.target.value)}
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          >
            <option value="">All Classes</option>
            <option value="active">Active</option>
            <option value="passive">Passive</option>
            <option value="toggle">Toggle</option>
          </select>
        </div>
      </FilterBar>

      <DataTable
        columns={COLUMNS}
        data={abilities ?? []}
        getKey={(row) => row.id}
        emptyMessage={
          filterType
            ? "No abilities match this filter."
            : "No abilities found. Create your first ability!"
        }
      />
    </PageContainer>
  );
}