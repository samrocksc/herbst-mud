import { Link } from '@tanstack/react-router'
import type { Room } from './types'
import type { NPC } from './types'

type MapSidebarProps = {
  rooms: Room[]
  npcs: NPC[]
  zLevels: Map<number, number>
  currentZLevel: number
  selectedRoom: Room | null
  setCurrentZLevel: (z: number) => void
  setSelectedRoom: (room: Room | null) => void
  setShowCreateModal: (show: boolean) => void
}

export function MapSidebar({
  rooms,
  npcs,
  zLevels,
  currentZLevel,
  selectedRoom,
  setCurrentZLevel,
  setSelectedRoom,
  setShowCreateModal,
}: MapSidebarProps) {
  const zLevelRange = Array.from(new Set(Array.from(zLevels.values()))).sort((a, b) => a - b)

  return (
    <div className="w-[220px] bg-surface-muted border-r border-border flex flex-col">
      <div className="p-4 border-b border-border">
        <Link
          to="/dashboard"
          className="block text-primary no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker"
        >
          ← Dashboard
        </Link>
        <button
          onClick={() => setShowCreateModal(true)}
          className="w-full p-2 bg-primary border-2 border-black rounded text-white cursor-pointer hover:bg-primary-hover"
        >
          + Add Room
        </button>
      </div>

      <div className="p-3 border-b border-border">
        <label className="text-text-muted text-xs block mb-2">Floor (Z-Level)</label>
        <div className="flex gap-1 flex-wrap">
          {zLevelRange.map(z => (
            <button
              key={z}
              onClick={() => setCurrentZLevel(z)}
              className={`px-2 py-1 rounded text-xs cursor-pointer ${
                currentZLevel === z
                  ? 'bg-primary border-primary-hover border'
                  : 'bg-surface-dark border-border border'
              } text-white`}
            >
              {z === 0 ? 'G' : z > 0 ? `+${z}` : `${z}`}
            </button>
          ))}
        </div>
      </div>

      <div className="p-3 text-text-muted text-xs border-b border-border">
        <div>Total: {rooms.length} rooms</div>
        <div>
          Floor {currentZLevel}: {Array.from(zLevels.values()).filter(z => z === currentZLevel).length}
        </div>
        <div>NPCs: {npcs.length}</div>
      </div>

      <div className="flex-1 overflow-y-auto p-3">
        <h4 className="m-0 mb-2 text-text-muted text-xs">Rooms on Floor {currentZLevel}</h4>
        <div className="flex flex-col gap-1">
          {rooms
            .filter(r => (zLevels.get(r.id) || 0) === currentZLevel)
            .map(room => (
              <div
                key={room.id}
                onClick={() => setSelectedRoom(room)}
                className={`p-2 cursor-pointer rounded text-xs room-node ${
                  selectedRoom?.id === room.id ? 'room-node--selected' : ''
                }`}
              >
                <span className="truncate">{room.name}</span>
                {room.isStartingRoom && <span> ⭐</span>}
              </div>
            ))}
        </div>
      </div>
    </div>
  )
}
