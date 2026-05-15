import { createFileRoute, Link, Outlet, useLocation, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { useAbilities } from '../../hooks/useAbilities';
import { PageHeader } from '../../components/PageHeader';
import { DataTable, type Column } from '../../components/DataTable';
import { Button } from '../../components/Button';
import type { Ability } from '../../hooks/useAbilities';

export const Route = createFileRoute('/_auth/abilities')({
  component: AbilitiesManagement,
});

const COLUMNS: Column<Ability>[] = [
  {
    header: 'Name',
    accessor: 'name',
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
  { header: 'Description', accessor: 'description' },
  {
    header: 'Type',
    accessor: 'ability_type',
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: 'Class',
    accessor: 'ability_class',
    render: (val: unknown) => (
      <span className={`talent-effect talent-effect-${String(val)}`}>{String(val)}</span>
    ),
  },
  {
    header: 'Costs',
    accessor: 'mana_cost',
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

function AbilitiesManagement() {
  const [filterType, setFilterType] = useState<string>('');
  const [filterClass, setFilterClass] = useState<string>('');
  const navigate = useNavigate();
  const location = useLocation();
  const { data: abilities, isLoading, error } = useAbilities({
    type: filterType || undefined,
    abilityClass: filterClass || undefined,
  });

  if (location.pathname !== '/abilities') return <Outlet />;

  if (isLoading) return <div className="loading">Loading abilities...</div>;
  if (error) return <div className="error">Failed to load abilities: {error.message}</div>;

  return (
    <div className="management-page">
      <PageHeader
        title="Abilities"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => navigate({ to: '/abilities/new' })}>
            + Add Ability
          </Button>
        }
      />

      <div className="filters-bar">
        <div className="filter-group">
          <label>Type:</label>
          <select value={filterType} onChange={(e) => setFilterType(e.target.value)}>
            <option value="">All Types</option>
            <option value="combat">Combat</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
            <option value="healing">Healing</option>
            <option value="support">Support</option>
            <option value="defensive">Defensive</option>
          </select>
        </div>
        <div className="filter-group">
          <label>Class:</label>
          <select value={filterClass} onChange={(e) => setFilterClass(e.target.value)}>
            <option value="">All Classes</option>
            <option value="active">Active</option>
            <option value="passive">Passive</option>
            <option value="toggle">Toggle</option>
          </select>
        </div>
        {(filterType || filterClass) && (
          <Button variant="ghost" size="sm" onClick={() => { setFilterType(''); setFilterClass(''); }}>
            Clear Filters
          </Button>
        )}
      </div>

      <DataTable
        columns={COLUMNS}
        data={abilities ?? []}
        getKey={(row) => row.id}
        emptyMessage={
          filterType
            ? 'No abilities match this filter.'
            : 'No abilities found. Create your first ability!'
        }
      />
    </div>
  );
}