import type { Room, NPC, Equipment } from './types'
import { DirectionShortLabels } from './DirectionUtils'

type RoomNodeProps = {
  room: Room
  pos: { x: number; y: number }
  isSelected: boolean
  roomNpcs: NPC[]
  roomItems: Equipment[]
  rooms: Room[]
  onSelect: (room: Room) => void
  isDragging: boolean
  onDragEnd: (roomId: number, x: number, y: number) => void
}

export function RoomNode({ room, pos, isSelected, roomNpcs, roomItems, rooms, onSelect, isDragging, onDragEnd }: RoomNodeProps) {
  const isColored = room.isStartingRoom || isSelected

  // --- Drag support ---
  const dragRef = { startMouseX: 0, startMouseY: 0, startPosX: 0, startPosY: 0 }

  function handleMouseDown(e: React.MouseEvent) {
    if ((e.target as HTMLElement).closest('button')) return
    e.preventDefault()
    dragRef.startMouseX = e.clientX
    dragRef.startMouseY = e.clientY
    dragRef.startPosX = pos.x
    dragRef.startPosY = pos.y
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }

  function handleMouseMove(e: MouseEvent) {
    const dx = e.clientX - dragRef.startMouseX
    const dy = e.clientY - dragRef.startMouseY
    const nodeEl = (e.currentTarget as HTMLElement)
    nodeEl.style.left = `${dragRef.startPosX + dx}px`
    nodeEl.style.top = `${dragRef.startPosY + dy}px`
  }

  function handleMouseUp(e: MouseEvent) {
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
    const dx = e.clientX - dragRef.startMouseX
    const dy = e.clientY - dragRef.startMouseY
    const newX = dragRef.startPosX + dx
    const newY = dragRef.startPosY + dy
    onDragEnd(room.id, newX, newY)
  }

  /** Resolve a room ID to its name, falling back to "Unknown Room" */
  function resolveRoomName(roomId: number): string {
    const found = rooms.find((r) => r.id === roomId)
    return found ? found.name : 'Unknown Room'
  }

  return (
    <div
      onClick={() => onSelect(room)}
      onMouseDown={handleMouseDown}
      className={`room-node absolute w-[120px] min-h-[65px] p-2 rounded-lg cursor-grab transition-all select-none ${
        isDragging ? 'opacity-50 cursor-grabbing' : ''
      } ${
        room.isStartingRoom
          ? 'bg-primary text-white'
          : isSelected
            ? 'bg-primary-hover text-white shadow-lg border-2 border-accent'
            : 'bg-surface text-text border-2 border-border hover:border-primary'
      }`}
      style={{ left: pos.x, top: pos.y, zIndex: isDragging ? 50 : 1 }}
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
      {room.exits && (
        <div className="flex flex-col items-center gap-0.5 mt-0.5">
          {Object.entries(room.exits).map(([dir, targetId]) => {
            const label = DirectionShortLabels[dir as keyof typeof DirectionShortLabels] ?? dir
            const targetName = resolveRoomName(targetId)
            return (
              <span
                key={dir}
                className={`text-[8px] leading-tight ${
                  dir === 'up'
                    ? isColored ? 'text-white/90' : 'text-warning'
                    : dir === 'down'
                      ? isColored ? 'text-white/80' : 'text-success'
                      : isColored ? 'text-white/70' : 'text-text-muted'
                }`}
              >
                {label}: {targetName}
              </span>
            )
          })}
        </div>
      )}
    </div>
  )
}