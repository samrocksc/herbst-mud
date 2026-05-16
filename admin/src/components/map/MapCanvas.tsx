/* eslint-disable functional/no-mixed-types */
import { useEffect } from "react";
import { ExitLines } from "./ExitLineRenderer";
import { RoomNode } from "./RoomNode";
import { CANVAS_W, CANVAS_H } from "./constants";
import type { Room, NPC, Equipment } from "./types";

type MapCanvasProps = Readonly<{
  rooms: Room[]
  nodePositions: Map<number, { x: number; y: number }>
  selectedRoom: Room | null
  zoom: number
  panOffset: { x: number; y: number }
  isDragging: boolean
  onWheel: (e: WheelEvent) => void
  onSelectRoom: (room: Room | null) => void
  onDragStart: (roomId: number) => void
  onDragEnd: (roomId: number, x: number, y: number) => void
  getNPCsInRoom: (roomId: number) => NPC[]
  getEquipmentInRoom: (roomId: number) => Equipment[]
  viewportRef: React.RefObject<HTMLDivElement | null>
}>

export function MapCanvas({
  rooms, nodePositions, selectedRoom, zoom, panOffset, isDragging,
  onWheel, onSelectRoom, onDragStart, onDragEnd,
  getNPCsInRoom, getEquipmentInRoom, viewportRef,
}: MapCanvasProps) {
  useEffect(() => {
    const el = viewportRef.current;
    if (!el) return;
    el.addEventListener("wheel", onWheel, { passive: false });
    return () => el.removeEventListener("wheel", onWheel);
  }, [onWheel, viewportRef]);

  return (
    <div ref={viewportRef} className="mt-[50px] h-[calc(100%-50px)] overflow-hidden p-6">
      <div
        className="relative"
        style={{
          width: CANVAS_W, height: CANVAS_H,
          transform: `translate(${panOffset.x}px, ${panOffset.y}px) scale(${zoom})`,
          transformOrigin: "top left",
        }}
      >
        <ExitLines rooms={rooms} nodePositions={nodePositions} />
        {rooms.map(room => {
          const pos = nodePositions.get(room.id);
          if (!pos) return null;
          return (
            <RoomNode
              key={room.id}
              room={room}
              pos={pos}
              isSelected={selectedRoom?.id === room.id}
              roomNpcs={getNPCsInRoom(room.id)}
              roomItems={getEquipmentInRoom(room.id)}
              rooms={rooms}
              zoom={zoom}
              onSelect={onSelectRoom}
              isDragging={isDragging}
              onDragStart={onDragStart}
              onDragEnd={onDragEnd}
            />
          );
        })}
      </div>
    </div>
  );
}