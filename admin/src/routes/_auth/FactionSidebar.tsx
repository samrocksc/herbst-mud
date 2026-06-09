import { Button } from "../../components/Button";
import { EMPTY_FORM, type Faction, type FactionForm } from "./-factionTypes";

export function FactionSidebar({
  factions,
  searchQuery,
  setSearchQuery,
  filteredFactions,
  selectedFaction,
  setSelectedFaction,
  setShowCreateForm,
  setForm,
}: Readonly<{
  factions: Faction[]
  searchQuery: string
  setSearchQuery: (q: string) => void
  filteredFactions: Faction[]
  selectedFaction: Faction | null
  setSelectedFaction: (f: Faction | null) => void
  setShowCreateForm: (v: boolean) => void
  setForm: (f: FactionForm) => void
}>) {
  // Sort: most recent first (highest id), then by display_name
  const sorted = [...filteredFactions].sort((a, b) => {
    if (b.id !== a.id) return b.id - a.id;
    return (a.display_name || a.name).localeCompare(b.display_name || b.name);
  });

  return (
    <>
      <div className="p-3 border-b border-border flex items-center justify-between">
        <div>
          <h2 className="m-0 text-text text-lg">Factions</h2>
          <p className="text-text-muted text-xs mt-1 mb-0">{factions.length} total</p>
        </div>
        <Button variant="primary" size="sm" onClick={() => { setShowCreateForm(true); setSelectedFaction(null); setForm(EMPTY_FORM); }}>
          + New
        </Button>
      </div>
      <div className="p-3 border-b border-border">
        <input
          type="text"
          placeholder="Search factions..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        <div className="flex flex-col gap-1">
          {sorted.map((f) => {
            const emoji = (f as Faction & { emoji?: string }).emoji;
            return (
              <div
                key={f.id}
                onClick={() => { setSelectedFaction(f); setShowCreateForm(false); }}
                className={`p-2 cursor-pointer rounded text-xs flex items-center gap-2 ${
                  selectedFaction?.id === f.id ? "bg-primary/20 text-primary font-medium" : "text-text hover:bg-surface-hover"
                }`}
              >
                {emoji && <span className="text-base shrink-0">{emoji}</span>}
                <div className="flex-1 min-w-0">
                  <div className="font-bold truncate">{f.display_name || f.name}</div>
                  <div className="text-text-muted text-xs">
                    {f.member_count != null ? `${f.member_count} member${f.member_count === 1 ? "" : "s"}` : "0 members"}
                  </div>
                </div>
              </div>
            );
          })}
          {sorted.length === 0 && (
            <div className="text-center py-6 space-y-2">
              <p className="text-text-muted text-sm">{searchQuery ? "No factions match" : "No factions yet"}</p>
              {!searchQuery && (
                <Button variant="accent" size="sm" onClick={() => { setShowCreateForm(true); setSelectedFaction(null); setForm(EMPTY_FORM); }}>
                  + Create your first faction
                </Button>
              )}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
