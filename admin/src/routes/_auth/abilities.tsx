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
  { header: "Description", accessor: "description" },
  {
    header: "Type",
    accessor: "ability_type",
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: "Class",
    accessor: "ability_class",
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
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
        <div className="flex items-center gap-2 mb-2">
          <span className="text-xs text-text-muted font-medium mr-1">Quick:</span>
          <button onClick={() => { setFilterClass("passive"); setFilterType(""); }}
            className={`px-2 py-1 text-xs rounded border ${filterClass === "passive" ? "bg-primary/20 border-primary text-text" : "bg-surface border-border text-muted hover:border-primary"}`}>
            Weapon Skills (Passive)
          </button>
          <button onClick={() => { setFilterClass("active"); setFilterType(""); }}
            className={`px-2 py-1 text-xs rounded border ${filterClass === "active" ? "bg-primary/20 border-primary text-text" : "bg-surface border-border text-muted hover:border-primary"}`}>
            Active Abilities
          </button>
          <button onClick={() => { setFilterClass(""); setFilterType(""); }}
            className={`px-2 py-1 text-xs rounded border ${filterClass === "" && filterType === "" ? "bg-primary/20 border-primary text-text" : "bg-surface border-border text-muted hover:border-primary"}`}>
            All
          </button>
        </div>
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