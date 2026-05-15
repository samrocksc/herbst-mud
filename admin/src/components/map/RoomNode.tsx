import { useRef, useCallback, memo } from 'react';
import { DRAG_THRESHOLD, GRID } from './constants';
import { DirectionShortLabels } from './DirectionUtils';
import type { Room, NPC, Equipment } from './types';

type RoomNodeProps = {
  room: Room
  pos: { x: number; y: number }
  isSelected: boolean
  roomNpcs: NPC[]
  roomItems: Equipment[]
  rooms: Room[]
  zoom: number
  onSelect: (room: Room) => void
  isDragging: boolean
  onDragStart: (roomId: number) => void
  onDragEnd: (roomId: number, x: number, y: number) => void
}

export const RoomNode = memo(function RoomNode({ room, pos, isSelected, roomNpcs, roomItems, rooms, zoom, onSelect, isDragging, onDragStart, onDragEnd }: RoomNodeProps) {
  const dragRef = useRef({ startMouseX: 0, startMouseY: 0, startPosX: 0, startPosY: 0, didDrag: false });

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    if ((e.target as HTMLElement).closest('button')) return;
    e.preventDefault();
    dragRef.current = {
      startMouseX: e.clientX, startMouseY: e.clientY,
      startPosX: pos.x, startPosY: pos.y,
      didDrag: false,
    };
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  }, [pos.x, pos.y, zoom]);

  const handleMouseMove = useCallback((e: MouseEvent) => {
    const dx = (e.clientX - dragRef.current.startMouseX) / zoom;
    const dy = (e.clientY - dragRef.current.startMouseY) / zoom;
    if (!dragRef.current.didDrag && Math.abs(dx) < DRAG_THRESHOLD && Math.abs(dy) < DRAG_THRESHOLD) return;
    if (!dragRef.current.didDrag) {
      dragRef.current.didDrag = true;
      onDragStart(room.id);
    }
    const rawX = dragRef.current.startPosX + dx;
    const rawY = dragRef.current.startPosY + dy;
    const snapX = Math.round(rawX / GRID) * GRID;
    const snapY = Math.round(rawY / GRID) * GRID;
    const nodeEl = document.querySelector(`[data-room-id="${room.id}"]`) as HTMLElement;
    if (nodeEl) {
      nodeEl.style.left = `${snapX}px`;
      nodeEl.style.top = `${snapY}px`;
    }
  }, [room.id, onDragStart, zoom]);

  const handleMouseUp = useCallback((e: MouseEvent) => {
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
    if (!dragRef.current.didDrag) return;
    const dx = (e.clientX - dragRef.current.startMouseX) / zoom;
    const dy = (e.clientY - dragRef.current.startMouseY) / zoom;
    onDragEnd(room.id, dragRef.current.startPosX + dx, dragRef.current.startPosY + dy);
  }, [room.id, onDragEnd, handleMouseMove, zoom]);

  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onSelect(room);
    }
  }, [room, onSelect]);

  const isColored = room.isRootRoom || room.isStartingRoom || isSelected;

  const validExits: Array<[string, number]> = room.exits
    ? (Object.entries(room.exits) as Array<[string, number]>).filter(([, tid]) => rooms.some(r => r.id === tid))
    : [];

  const resolveRoomName = (id: number) => {
    const found = rooms.find(r => r.id === id);
    return found ? found.name : 'Unknown Room';
  };

  return (
    <div
      data-room-id={room.id}
      role="button"
      tabIndex={0}
      aria-label={`${room.name} (Room ${room.id})`}
      onClick={() => { if (!dragRef.current.didDrag) onSelect(room); }}
      onMouseDown={handleMouseDown}
      onKeyDown={handleKeyDown}
      className={`room-node absolute w-[120px] min-h-[65px] p-2 rounded-lg cursor-grab transition-all select-none ${
        isDragging ? 'opacity-50 cursor-grabbing' : ''
      } ${
        room.isRootRoom
          ? 'bg-accent text-white'
          : room.isStartingRoom
            ? 'bg-primary text-white'
            : isSelected
            ? 'bg-primary-hover text-white shadow-lg border-2 border-accent'
            : 'bg-surface text-text border-2 border-border hover:border-primary'
      }`}
      style={{ left: pos.x, top: pos.y, zIndex: isDragging ? 50 : 1 }}
    >
      <div className={`font-bold text-xs text-center truncate ${isColored ? 'text-white' : 'text-text'}`}>
        {room.name}
        {room.isRootRoom && ' 🏠'}
        {room.isStartingRoom && !room.isRootRoom && ' ⭐'}
      </div>
      <div className={`text-xs text-center ${isColored ? 'text-white/80' : 'text-text-muted'}`}>
        #{room.id}
      </div>
      <div className="flex justify-center gap-1 mt-1">
        {roomNpcs.length > 0 && (
          <span className={`text-[10px] ${isColored ? 'text-white/90' : 'text-warning'}`} title={`${roomNpcs.length} NPCs`}>
            👥{roomNpcs.length}
          </span>
        )}
        {roomItems.length > 0 && (
          <span className={`text-[10px] ${isColored ? 'text-white/80' : 'text-success'}`} title={`${roomItems.length} items`}>
            📦{roomItems.length}
          </span>
        )}
      </div>
      {validExits.length > 0 && (
        <div className="flex flex-col items-center gap-0.5 mt-0.5">
          {validExits.map(([dir, targetId]) => {
            const label = DirectionShortLabels[dir as keyof typeof DirectionShortLabels] ?? dir;
            return (
              <span
                key={dir}
                className={`text-[8px] leading-tight ${
                  dir === 'up' ? isColored ? 'text-white/90' : 'text-warning'
                  : dir === 'down' ? isColored ? 'text-white/80' : 'text-success'
                  : isColored ? 'text-white/70' : 'text-text-muted'
                }`}
              >
                {label}: {resolveRoomName(targetId)}
              </span>
            );
          })}
        </div>
      )}
    </div>
  );
});