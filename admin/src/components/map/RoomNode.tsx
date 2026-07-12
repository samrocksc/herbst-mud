import { useCallback, memo } from "react";
import { Handle, Position } from "@xyflow/react";
import { NODE_W, NODE_H } from "./constants";
import { DirectionShortLabels } from "./DirectionUtils";
import type { Room, NPC, Equipment } from "./types";

// eslint-disable-next-line functional/no-mixed-types -- component props conventionally mix data fields and callback handlers
type RoomNodeProps = Readonly<{
  room: Room
  pos: Readonly<{ x: number; y: number }>
  isSelected: boolean
  roomNpcs: ReadonlyArray<NPC>
  roomItems: ReadonlyArray<Equipment>
  rooms: ReadonlyArray<Room>
  zoom: number
  onSelect: (room: Room) => void
}>

export const RoomNode = memo(function RoomNode({ room, isSelected, roomNpcs, roomItems, rooms, onSelect }: RoomNodeProps) {
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === "Enter" || e.key === " ") {
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
    return found ? found.name : "Unknown Room";
  };

  return (
    <div
      data-room-id={room.id}
      role="button"
      tabIndex={0}
      aria-label={`${room.name} (Room ${room.id})`}
      onClick={() => { onSelect(room); }}
      onKeyDown={handleKeyDown}
      className={`room-node overflow-hidden p-2 rounded-lg cursor-pointer transition-all select-none ${
        room.isRootRoom
          ? "bg-accent text-white"
          : room.isStartingRoom
            ? "bg-primary text-white"
            : isSelected
            ? "bg-yellow-400 text-black shadow-2xl border-4 border-red-600 ring-4 ring-yellow-200 scale-110 z-50"
            : "bg-surface text-text border-2 border-border hover:border-primary"
      }`}
      style={{ zIndex: isSelected ? 50 : 1, width: NODE_W, height: NODE_H }}
    >
      <div className={`font-bold text-xs text-center truncate w-full ${isColored ? "text-white" : "text-text"}`}>
        {room.name}
        {room.isRootRoom && " 🏠"}
        {room.isStartingRoom && !room.isStartingRoom && " ⭐"}
      </div>
      <div className={`text-xs text-center ${isColored ? "text-white/80" : "text-text-muted"}`}>
        #{room.id}
      </div>
      <div className="flex justify-center gap-1 mt-1">
        {roomNpcs.length > 0 && (
          <span className={`text-[10px] ${isColored ? "text-white/90" : "text-warning"}`} title={`${roomNpcs.length} NPCs`}>
            👥{roomNpcs.length}
          </span>
        )}
        {roomItems.length > 0 && (
          <span className={`text-[10px] ${isColored ? "text-white/80" : "text-success"}`} title={`${roomItems.length} items`}>
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
                className={`text-[8px] leading-tight w-full truncate text-center ${
                  dir === "up" ? isColored ? "text-white/90" : "text-warning"
                  : dir === "down" ? isColored ? "text-white/80" : "text-success"
                  : isColored ? "text-white/70" : "text-text-muted"
                }`}
              >
                {label}: {resolveRoomName(targetId)}
              </span>
            );
          })}
        </div>
      )}
      <Handle type="target" position={Position.Top} id="north" style={{ visibility: "hidden" }} />
      <Handle type="source" position={Position.Top} id="north" style={{ visibility: "hidden" }} />
      <Handle type="target" position={Position.Bottom} id="south" style={{ visibility: "hidden" }} />
      <Handle type="source" position={Position.Bottom} id="south" style={{ visibility: "hidden" }} />
      <Handle type="target" position={Position.Right} id="east" style={{ visibility: "hidden" }} />
      <Handle type="source" position={Position.Right} id="east" style={{ visibility: "hidden" }} />
      <Handle type="target" position={Position.Left} id="west" style={{ visibility: "hidden" }} />
      <Handle type="source" position={Position.Left} id="west" style={{ visibility: "hidden" }} />
    </div>
  );
});
