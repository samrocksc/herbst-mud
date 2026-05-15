import { createFileRoute, Link, Outlet, useLocation } from '@tanstack/react-router';
import { useState, useMemo } from 'react';
import { useCharacters, type Character } from '../../hooks/useCharacters';
import { PageHeader } from '../../components/PageHeader';
import { DataTable, type Column } from '../../components/DataTable';
import { fuzzyMatch } from '../../components/fuzzyMatch';

export const Route = createFileRoute('/_auth/characters')({
  component: CharactersIndex,
});

function CharactersIndex() {
  const { data: characters, isLoading, isError, error } = useCharacters();
  const [searchQuery, setSearchQuery] = useState('');
  const [showNPCs, setShowNPCs] = useState(false);
  const [showOnlineOnly, setShowOnlineOnly] = useState(false);

  const filteredCharacters = useMemo(() => {
    const list = (characters ?? []).filter((c) => {
      if (!showNPCs && c.isNPC) return false;
      if (searchQuery && !fuzzyMatch(c.name, searchQuery) && !fuzzyMatch(String(c.id), searchQuery)) return false;
      if (showOnlineOnly) {
        if (!c.lastSeenAt) return false;
        const lastSeen = new Date(c.lastSeenAt);
        const fifteenMinAgo = new Date(Date.now() - 15 * 60 * 1000);
        if (lastSeen < fifteenMinAgo) return false;
      }
      return true;
    });
    return list;
  }, [characters, searchQuery, showNPCs, showOnlineOnly]);

  const columns: Column<Character>[] = [
    {
      header: 'ID',
      accessor: 'id',
      render: (_, row) => <span className="font-mono text-xs">{row.id}</span>,
    },
    {
      header: 'Name',
      accessor: 'name',
      className: 'font-bold',
      render: (_, row) => (
        <Link
          to="/characters/$characterId"
          params={{ characterId: String(row.id) }}
          className="no-underline text-primary hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    { header: 'Race', accessor: 'race' },
    { header: 'Class', accessor: 'class' },
    { header: 'Level', accessor: 'level', align: 'center' },
    {
      header: 'Room',
      accessor: 'currentRoomId',
      align: 'center',
      render: (_, row) => <span className="font-mono text-xs">#{row.currentRoomId}</span>,
    },
    {
      header: 'Status',
      accessor: 'lastSeenAt',
      render: (_, row) => {
        if (row.isNPC) return <span className="badge badge-neutral">NPC</span>;
        if (!row.lastSeenAt) return <span className="badge badge-warning">Offline</span>;
        const lastSeen = new Date(row.lastSeenAt);
        const fifteenMinAgo = new Date(Date.now() - 15 * 60 * 1000);
        if (lastSeen >= fifteenMinAgo) {
          return <span className="badge badge-success">Online</span>;
        }
        return <span className="badge badge-warning">Offline</span>;
      },
    },
  ];

  const location = useLocation();
  const isList = location.pathname === '/characters';

  if (!isList) {
    return <Outlet />;
  }

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title="Characters" showBack backTo="/dashboard" />

      <div className="flex gap-3 mb-4 flex-wrap">
        <input
          type="text"
          placeholder="Search by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-xs p-2 bg-surface border border-border rounded text-text text-sm"
        />
        <label className="flex items-center gap-1 text-sm text-text-muted cursor-pointer">
          <input
            type="checkbox"
            checked={showNPCs}
            onChange={(e) => setShowNPCs(e.target.checked)}
            className="accent-primary"
          />
          Show NPCs
        </label>
        <label className="flex items-center gap-1 text-sm text-text-muted cursor-pointer">
          <input
            type="checkbox"
            checked={showOnlineOnly}
            onChange={(e) => setShowOnlineOnly(e.target.checked)}
            className="accent-primary"
          />
          Online only
        </label>
      </div>

      {isLoading && <div className="p-8 text-text-muted text-center text-xs">Loading characters...</div>}
      {isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load characters: {error?.message ?? 'Unknown error'}
        </div>
      )}
      {!isLoading && !isError && (
        <DataTable<Character>
          columns={columns}
          data={filteredCharacters}
          getKey={(row) => row.id}
          emptyMessage="No characters found."
          variant="dark"
        />
      )}
    </div>
  );
}