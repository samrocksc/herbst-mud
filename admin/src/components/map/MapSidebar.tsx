/* eslint-disable functional/no-mixed-types */
import { Link } from "@tanstack/react-router";
import { Button } from "../Button";
import { DashboardIcon } from "../icons/DashboardIcon";
import { NPCsIcon } from "../icons/NPCsIcon";
import { ItemsIcon } from "../icons/ItemsIcon";
import type { Room, NPC } from "./types";
import { ZonesPanel } from "./ZonesPanel";

type MapSidebarProps = Readonly<{
  rooms: Room[]
  npcs: NPC[]
  zLevels: Map<number, number>
  currentZLevel: number
  selectedRoom: Room | null
  setCurrentZLevel: (z: number) => void
  setSelectedRoom: (room: Room | null) => void
}>

export function MapSidebar({
  rooms,
  npcs,
  zLevels,
  currentZLevel,
  selectedRoom,
  setCurrentZLevel,
  setSelectedRoom,
}: MapSidebarProps) {
  const zLevelRange = Array.from(new Set(Array.from(zLevels.values()))).sort((a, b) => a - b);
  const roomsOnFloor = Array.from(zLevels.values()).filter((z) => z === currentZLevel).length;

  return (
    <div className="w-[220px] h-full bg-surface-muted border-r border-border flex flex-col flex-shrink-0">
      {/* Mobile close + sidebar toggle on top — shown only on mobile */}
      <div className="flex items-center justify-between p-3 border-b border-border lg:hidden">
        <Button variant="ghost" size="sm" onClick={() => { /* handled by parent backdrop */ }} aria-label="Close">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </Button>
      </div>
      {/* Dashboard + secondary nav */}
      <div className="p-3 border-b border-border flex flex-col gap-1">
        <Link
          to="/dashboard"
          activeProps={{
            className: "bg-primary/10 text-primary border-l-4 border-primary font-semibold",
          }}
          inactiveProps={{
            className: "text-text-muted hover:bg-surface-muted hover:text-text",
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <DashboardIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">Dashboard</span>
        </Link>
        <Link
          to="/npcs"
          activeProps={{
            className: "bg-primary/10 text-primary border-l-4 border-primary font-semibold",
          }}
          inactiveProps={{
            className: "text-text-muted hover:bg-surface-muted hover:text-text",
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <NPCsIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">NPCs</span>
        </Link>
        <Link
          to="/items"
          activeProps={{
            className: "bg-primary/10 text-primary border-l-4 border-primary font-semibold",
          }}
          inactiveProps={{
            className: "text-text-muted hover:bg-surface-muted hover:text-text",
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <ItemsIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">Items</span>
        </Link>
      </div>

      {/* Zones section (admin can manually add zones, manage room chips) */}
      <ZonesPanel />

      {/* Add Room button */}
      <div className="p-3 border-b border-border">
        <Link to="/map/rooms/new" className="no-underline">
          <Button variant="primary" size="md" fullWidth>
            + Add Room
          </Button>
        </Link>
      </div>

      {/* Floor selector */}
      <div className="p-3 border-b border-border">
        <label className="text-text-muted text-xs block mb-2">
          Floor (Z-Level)
        </label>
        <div className="flex gap-1 flex-wrap">
          {zLevelRange.map((z) => (
            <Button
              key={z}
              variant={currentZLevel === z ? "primary" : "secondary"}
              size="sm"
              onClick={() => setCurrentZLevel(z)}
            >
              {z === 0 ? "G" : z > 0 ? `+${z}` : `${z}`}
            </Button>
          ))}
        </div>
      </div>

      {/* Stats */}
      <div className="p-3 text-text-muted text-xs border-b border-border">
        <div>Total: {rooms.length} rooms</div>
        <div>Floor {currentZLevel}: {roomsOnFloor}</div>
        <div>NPCs: {npcs.length}</div>
      </div>

      {/* Room list */}
      <div className="flex-1 overflow-y-auto p-3">
        <h4 className="m-0 mb-2 text-text-muted text-xs font-semibold uppercase tracking-wide">
          Rooms on Floor {currentZLevel}
        </h4>
        <div className="flex flex-col gap-1">
          {rooms
            .filter((r) => (zLevels.get(r.id) || 0) === currentZLevel)
            .map((room) => (
              <div
                key={room.id}
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  setSelectedRoom(room);
                }}
                className={[
                  "p-2 rounded text-xs cursor-pointer transition-colors",
                  selectedRoom?.id === room.id
                    ? "bg-primary/10 text-text border border-primary/30"
                    : "text-text-muted hover:bg-surface hover:text-text",
                ].join(" ")}
              >
                <span className="truncate block">{room.name}</span>
                {room.isRootRoom && (
                  <span className="text-accent text-[10px]"> 🏠 Root</span>
                )}
                {!room.isRootRoom && room.isStartingRoom && (
                  <span className="text-warning text-[10px]"> ⭐ Start</span>
                )}
              </div>
            ))}
        </div>
      </div>
    </div>
  );
}
