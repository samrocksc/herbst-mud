// RoomZonesField — multi-select for the zones a room belongs to.
//
// Used inside the room editor. Lists all zones for the current world with
// checkboxes; user toggles membership. Selected zones show as removable
// chips. Unselected zones are listed in a scrollable list.
//
// Updates are committed via the parent room editor's save flow
// (we do not auto-save; we just hold the chosen ids in form state).

import { useState, useMemo } from "react";
import { useZones } from "../../hooks/useZones";

type RoomZonesFieldProps = Readonly<{
  selectedZoneIds: string[];
  onChange: (next: string[]) => void;
}>;

export function RoomZonesField({ selectedZoneIds, onChange }: RoomZonesFieldProps) {
  const { zones, isLoading } = useZones();
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return zones;
    return zones.filter(
      (z) =>
        z.id.toLowerCase().includes(q) ||
        z.name.toLowerCase().includes(q),
    );
  }, [zones, search]);

  const selectedSet = useMemo(() => new Set(selectedZoneIds), [selectedZoneIds]);

  const toggle = (id: string) => {
    if (selectedSet.has(id)) {
      onChange(selectedZoneIds.filter((z) => z !== id));
    } else {
      onChange([...selectedZoneIds, id]);
    }
  };

  const removeChip = (id: string) => onChange(selectedZoneIds.filter((z) => z !== id));

  return (
    <div className="space-y-2">
      <label className="text-text-muted text-xs block">Zones</label>
      <div className="flex flex-wrap gap-1 min-h-[24px]">
        {selectedZoneIds.length === 0 ? (
          <span className="text-text-muted text-xs italic">Not in any zone.</span>
        ) : (
          selectedZoneIds.map((zid) => {
            const z = zones.find((z) => z.id === zid);
            const color = z?.color;
            return (
              <span
                key={zid}
                className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded border border-border bg-surface text-xs"
              >
                {color && (
                  <span
                    className="inline-block w-2.5 h-2.5 rounded-sm"
                    style={{ backgroundColor: color }}
                    aria-hidden="true"
                  />
                )}
                {z?.name ?? zid}
                <button
                  type="button"
                  onClick={() => removeChip(zid)}
                  className="text-text-muted hover:text-danger"
                  aria-label={"Remove zone " + (z?.name ?? zid)}
                >
                  ×
                </button>
              </span>
            );
          })
        )}
      </div>
      <input
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Filter zones…"
        className="w-full bg-surface border border-border rounded px-1.5 py-0.5 text-xs"
      />
      <div className="max-h-32 overflow-y-auto bg-surface border border-border rounded">
        {isLoading ? (
          <div className="px-2 py-1 text-text-muted text-xs italic">Loading…</div>
        ) : filtered.length === 0 ? (
          <div className="px-2 py-1 text-text-muted text-xs italic">
            {zones.length === 0
              ? "No zones defined. Create one in the Zones section of the sidebar."
              : "No zones match the filter."}
          </div>
        ) : (
          filtered.map((z) => (
            <label
              key={z.id}
              className="flex items-center gap-2 px-2 py-0.5 hover:bg-surface-muted text-xs cursor-pointer"
            >
              <input
                type="checkbox"
                checked={selectedSet.has(z.id)}
                onChange={() => toggle(z.id)}
                className="w-3.5 h-3.5 rounded border-border bg-surface text-primary focus:ring-primary"
              />
              {z.color && (
                <span
                  className="inline-block w-2.5 h-2.5 rounded-sm"
                  style={{ backgroundColor: z.color }}
                  aria-hidden="true"
                />
              )}
              <span className="flex-1 truncate">{z.name}</span>
              <span className="text-text-muted text-[10px] font-mono">{z.id}</span>
            </label>
          ))
        )}
      </div>
    </div>
  );
}
