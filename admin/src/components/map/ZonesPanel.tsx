// ZonesPanel — accordion section for the map editor sidebar.
//
// Lets the admin:
//   1. Create a zone manually (id, name, color, min_level)
//   2. En masse add rooms to a zone (chip fashion; rooms that don't exist
//      are rejected — the user types/searches and only valid room ids can
//      be added)
//   3. See ghost rooms in red chips (rooms that have been removed from the
//      world but are still referenced in the zone)
//
// This component is world-scoped — it reads the current world from
// WorldStoreContext and does not expose world_id as a form field.

import { useState, useMemo } from "react";
import { useZones, useZoneRooms, type Zone, type ZoneRoom } from "../../hooks/useZones";
import { useRooms } from "../../hooks/useRooms";

type ZonesPanelProps = Readonly<{
  // Optional rooms prop — if not provided, useRooms is called internally.
  // Tests inject rooms to skip the network call.
  roomsProp?: ReadonlyArray<{ id: number; name: string }>;
}>;

function genZoneId(): string {
  return "z_" + Math.random().toString(36).slice(2, 10);
}

export function ZonesPanel({ roomsProp }: ZonesPanelProps) {
  const { zones, isLoading, createZone, isCreating } = useZones();
  const { rooms: hookRooms } = useRooms();
  const rooms = roomsProp ?? hookRooms;

  const [expanded, setExpanded] = useState(false);
  const [showCreate, setShowCreate] = useState(false);

  return (
    <div className="p-3 border-b border-border">
      <button
        type="button"
        onClick={() => setExpanded((v) => !v)}
        className="w-full flex items-center justify-between text-text font-semibold text-sm"
      >
        <span>Zones</span>
        <span className="text-text-muted text-xs">
          {isLoading ? "…" : `${zones.length}`} {expanded ? "▾" : "▸"}
        </span>
      </button>

      {expanded && (
        <div className="mt-2 space-y-2">
          <button
            type="button"
            onClick={() => setShowCreate((v) => !v)}
            className="text-xs text-primary hover:underline"
            disabled={isCreating}
          >
            {showCreate ? "− cancel" : "+ New zone"}
          </button>
          {showCreate && (
            <CreateZoneForm
              onSubmit={async (input) => {
                await createZone({
                  id: input.id,
                  name: input.name,
                  color: input.color,
                  min_level: input.min_level,
                  description: input.description,
                });
                setShowCreate(false);
              }}
              isSubmitting={isCreating}
            />
          )}
          {zones.map((z) => (
            <ZoneRow key={z.id} zone={z} rooms={rooms} />
          ))}
          {zones.length === 0 && !isLoading && (
            <div className="text-text-muted text-xs italic">No zones yet.</div>
          )}
        </div>
      )}
    </div>
  );
}

function CreateZoneForm({
  onSubmit,
  isSubmitting,
}: {
  onSubmit: (input: { id: string; name: string; color: string; min_level: number; description: string }) => Promise<void>;
  isSubmitting: boolean;
}) {
  const [id, setId] = useState(genZoneId());
  const [name, setName] = useState("");
  const [color, setColor] = useState("#8b4513");
  const [minLevel, setMinLevel] = useState(1);
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    try {
      await onSubmit({
        id: id.trim(),
        name: name.trim(),
        color: color || "#8b4513",
        min_level: Math.max(0, minLevel),
        description: description.trim(),
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-surface-muted rounded p-2 space-y-1.5 text-xs">
      <div>
        <label className="block text-text-muted">ID</label>
        <input
          value={id}
          onChange={(e) => setId(e.target.value)}
          className="w-full bg-surface border border-border rounded px-1 py-0.5 text-xs font-mono"
        />
      </div>
      <div>
        <label className="block text-text-muted">Name *</label>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full bg-surface border border-border rounded px-1 py-0.5"
        />
      </div>
      <div className="flex gap-2">
        <div className="flex-1">
          <label className="block text-text-muted">Color</label>
          <input
            type="color"
            value={color}
            onChange={(e) => setColor(e.target.value)}
            className="w-full h-6 bg-surface border border-border rounded"
          />
        </div>
        <div className="w-16">
          <label className="block text-text-muted">Min Lv</label>
          <input
            type="number"
            value={minLevel}
            onChange={(e) => setMinLevel(parseInt(e.target.value, 10) || 0)}
            className="w-full bg-surface border border-border rounded px-1 py-0.5"
            min={0}
          />
        </div>
      </div>
      <div>
        <label className="block text-text-muted">Description</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={2}
          className="w-full bg-surface border border-border rounded px-1 py-0.5"
        />
      </div>
      {error && <div className="text-danger">{error}</div>}
      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full bg-primary text-white rounded px-2 py-1 text-xs hover:bg-primary/90 disabled:opacity-50"
      >
        {isSubmitting ? "Creating…" : "Create zone"}
      </button>
    </form>
  );
}

function ZoneRow({
  zone,
  rooms,
}: {
  zone: Zone;
  rooms: ReadonlyArray<{ id: number; name: string }>;
}) {
  const [expanded, setExpanded] = useState(false);
  const { deleteZone, isDeleting } = useZones();

  return (
    <div className="bg-surface-muted rounded">
      <div className="flex items-center gap-2 p-2">
        <span
          className="inline-block w-3 h-3 rounded-sm border border-border"
          style={{ backgroundColor: zone.color || "#888" }}
          aria-label={"Zone color " + (zone.color || "default")}
        />
        <button
          type="button"
          onClick={() => setExpanded((v) => !v)}
          className="flex-1 text-left text-xs font-medium text-text hover:text-primary"
        >
          {zone.name} <span className="text-text-muted">({zone.id})</span>
        </button>
        <button
          type="button"
          onClick={() => {
            if (confirm("Delete zone " + zone.name + "?")) {
              void deleteZone(zone.id);
            }
          }}
          disabled={isDeleting}
          className="text-danger hover:underline text-xs"
        >
          ×
        </button>
      </div>
      {expanded && <ZoneRoomsEditor zone={zone} rooms={rooms} />}
    </div>
  );
}

function ZoneRoomsEditor({
  zone,
  rooms,
}: {
  zone: Zone;
  rooms: ReadonlyArray<{ id: number; name: string }>;
}) {
  const { rooms: zoneRooms, isLoading, addRooms, removeRoom, isAdding, isRemoving } =
    useZoneRooms(zone.id);
  const [search, setSearch] = useState("");
  const [error, setError] = useState("");

  const roomById = useMemo(() => {
    const m = new Map<number, string>();
    for (const r of rooms) m.set(r.id, r.name);
    return m;
  }, [rooms]);

  const candidates = useMemo(() => {
    const present = new Set<number>();
    for (const zr of zoneRooms) present.add(zr.id);
    const q = search.trim().toLowerCase();
    return rooms
      .filter((r) => !present.has(r.id))
      .filter((r) => {
        if (!q) return true;
        return (
          String(r.id).includes(q) ||
          r.name.toLowerCase().includes(q)
        );
      })
      .slice(0, 25);
  }, [rooms, zoneRooms, search]);

  const handleAdd = async (id: number) => {
    setError("");
    try {
      await addRooms([id]);
      setSearch("");
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  const handleRemove = async (id: number) => {
    setError("");
    try {
      await removeRoom(id);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  const roomLabel = (zr: ZoneRoom) =>
    zr.exists
      ? `${zr.id} · ${zr.name ?? roomById.get(zr.id) ?? ""}`
      : `${zr.id} · (removed)`;

  return (
    <div className="p-2 space-y-1.5 text-xs border-t border-border">
      <div className="flex flex-wrap gap-1">
        {zoneRooms.map((zr) => {
          const cls = zr.exists
            ? "bg-surface border-border text-text"
            : "bg-danger/10 border-danger text-danger";
          return (
            <span
              key={zr.id}
              title={zr.exists ? roomLabel(zr) : (zr.message ?? "room missing")}
              className={`inline-flex items-center gap-1 px-1.5 py-0.5 rounded border ${cls}`}
            >
              {roomLabel(zr)}
              <button
                type="button"
                onClick={() => void handleRemove(zr.id)}
                disabled={isRemoving}
                className="text-text-muted hover:text-danger"
                aria-label="Remove room from zone"
              >
                ×
              </button>
            </span>
          );
        })}
        {isLoading && <span className="text-text-muted">loading…</span>}
        {!isLoading && zoneRooms.length === 0 && (
          <span className="text-text-muted italic">No rooms in zone yet.</span>
        )}
      </div>

      <div className="space-y-1">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search room by id or name…"
          className="w-full bg-surface border border-border rounded px-1 py-0.5 text-xs"
        />
        {search && (
          <div className="max-h-32 overflow-y-auto bg-surface border border-border rounded">
            {candidates.length === 0 ? (
              <div className="px-2 py-1 text-text-muted italic">
                No matching rooms.
              </div>
            ) : (
              candidates.map((r) => (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => void handleAdd(r.id)}
                  disabled={isAdding}
                  className="w-full text-left px-2 py-0.5 hover:bg-surface-muted text-xs"
                >
                  {r.id} · {r.name}
                </button>
              ))
            )}
          </div>
        )}
      </div>

      {error && <div className="text-danger">{error}</div>}
    </div>
  );
}
