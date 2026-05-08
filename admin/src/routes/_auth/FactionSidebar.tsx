import { Button } from '../../components/Button'
import { EMPTY_FORM, type Faction, type FactionForm } from './factionTypes'

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
  return (
    <>
      <div className="p-3 border-b border-border">
        <h2 className="m-0 text-text text-lg">Factions</h2>
        <p className="text-text-muted text-xs mt-1 mb-0">{factions.length} factions</p>
      </div>
      <div className="p-3 border-b border-border">
        <input type="text" placeholder="Search factions..." value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        <div className="flex flex-col gap-1">
          {filteredFactions.map((f) => (
            <div key={f.id}
              onClick={() => { setSelectedFaction(f); setShowCreateForm(false) }}
              className={`p-2 cursor-pointer rounded text-xs ${
                selectedFaction?.id === f.id ? 'text-primary bg-primary/20 font-medium' : 'text-text'
              }`}>
              <div className="font-bold">{f.display_name || f.name}</div>
              <div className="text-text-muted">{f.is_universal ? 'universal' : `standing: ${f.standing ?? 0}`}</div>
            </div>
          ))}
          {filteredFactions.length === 0 && <div className="text-text-muted text-center py-4">No factions found</div>}
        </div>
      </div>
      <div className="p-3 border-t border-border">
        <Button variant="primary" size="md" fullWidth
          onClick={() => { setShowCreateForm(true); setSelectedFaction(null); setForm(EMPTY_FORM) }}>
          + Add Faction
        </Button>
      </div>
    </>
  )
}