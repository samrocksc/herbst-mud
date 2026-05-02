import type { Room, NPC, Equipment } from './types'

type RoomNodeProps = {
  room: Room
  pos: { x: number; y: number }
  isSelected: boolean
  roomNpcs: NPC[]
  roomItems: Equipment[]
  onSelect: (room: Room) => void
}

export function RoomNode({ room, pos, isSelected, roomNpcs, roomItems, onSelect }: RoomNodeProps) {
  const isColored = room.isStartingRoom || isSelected

  return (
    <div
      onClick={() => onSelect(room)}
      className={`room-node absolute w-[120px] min-h-[65px] p-2 rounded-lg cursor-pointer transition-all ${
        room.isStartingRoom
          ? 'bg-primary text-white'
          : isSelected
            ? 'bg-primary-hover text-white shadow-lg border-2 border-accent'
            : 'bg-surface text-text border-2 border-border hover:border-primary'
      }`}
      style={{ left: pos.x, top: pos.y, zIndex: 1 }}
    >
      <div
        className={`font-bold text-xs text-center truncate ${
          isColored ? 'text-white' : 'text-text'
        }`}
      >
        {room.name}
        {room.isStartingRoom && ' ⭐'}
      </div>
      <div className={`text-xs text-center ${isColored ? 'text-white/80' : 'text-text-muted'}`}>
        #{room.id}
      </div>
      <div className="flex justify-center gap-1 mt-1">
        {roomNpcs.length > 0 && (
          <span
            className={`text-[10px] ${isColored ? 'text-white/90' : 'text-warning'}`}
            title={`${roomNpcs.length} NPCs`}
          >
            👥{roomNpcs.length}
          </span>
        )}
        {roomItems.length > 0 && (
          <span
            className={`text-[10px] ${isColored ? 'text-white/80' : 'text-success'}`}
            title={`${roomItems.length} items`}
          >
            📦{roomItems.length}
          </span>
        )}
      </div>
      <div className="flex justify-center gap-0.5 mt-0.5">
        {room.exits?.up && (
          <span className={`text-[8px] ${isColored ? 'text-white/90' : 'text-warning'}`}>
            ▲{room.exits.up}
          </span>
        )}
        {room.exits?.down && (
          <span className={`text-[8px] ${isColored ? 'text-white/80' : 'text-success'}`}>
            ▼{room.exits.down}
          </span>
        )}
      </div>
    </div>
  )
}
