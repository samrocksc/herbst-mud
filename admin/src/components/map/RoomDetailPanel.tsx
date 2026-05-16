/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types */
import { useState, memo } from "react";
import type { Room } from "./types";
import { ALL_DIRECTIONS } from "./DirectionUtils";
import { Button } from "../Button";
import { NPCInstanceManager } from "./NPCInstanceManager";
import { ItemInstanceManager } from "./ItemInstanceManager";

type RoomDetailPanelProps = {
  selectedRoom: Room
  rooms: Room[]
  zLevels: Map<number, number>
  onSelectRoom: (room: Room | null) => void
  onEditRoom: (room: Room) => void
  onDeleteRoom: (roomId: number) => void
  onAddRoom?: (room: Room, dir: string) => void
}

export const RoomDetailPanel = memo(function RoomDetailPanel({
  selectedRoom,
  rooms,
  zLevels,
  onSelectRoom,
  onEditRoom,
  onDeleteRoom,
  onAddRoom,
}: RoomDetailPanelProps) {
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-base font-semibold">
          {selectedRoom.name}
          {selectedRoom.isStartingRoom && (
            <span className="text-warning ml-1">⭐</span>
          )}
        </h3>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onSelectRoom(null)}
          aria-label="Close"
        >
          ×
        </Button>
      </div>

      <div className="p-3 flex-1 overflow-y-auto">
        <div className="text-text-muted text-[10px] mb-2">
          Room ID: <a href={`/map?room=${selectedRoom.id}`} className="text-primary hover:underline" onClick={(e) => { e.preventDefault(); navigator.clipboard.writeText(`${window.location.origin}/map?room=${selectedRoom.id}`); }}>#{selectedRoom.id}</a>
          {selectedRoom.atmosphere && ` • ${selectedRoom.atmosphere}`}
        </div>
        <div className="text-text mb-3 text-sm">{selectedRoom.description}</div>

        <NPCInstanceManager roomId={selectedRoom.id} />
        <ItemInstanceManager roomId={selectedRoom.id} />

        <div className="mb-3">
          <strong className="text-accent text-xs">Exits:</strong>
          <div className="mt-1">
            {ALL_DIRECTIONS.map((dir) => {
              const targetId = selectedRoom.exits?.[dir];
              const targetRoom = rooms.find((r) => r.id === targetId);
              const isZExit = dir === "up" || dir === "down";

              if (targetId && targetRoom) {
                return (
                  <div
                    key={dir}
                    onClick={() => onSelectRoom(targetRoom)}
                    className={[
                      "p-1 my-1 rounded cursor-pointer text-xs transition-colors",
                      isZExit
                        ? dir === "up"
                          ? "bg-warning/20 border border-warning"
                          : "bg-success/20 border border-success"
                        : "bg-surface-muted",
                    ].join(" ")}
                  >
                    <strong>{dir}</strong> → {targetRoom.name}
                    {isZExit && (
                      <span className="text-text-muted ml-1 text-[10px]">
                        (z={zLevels.get(targetId) || 0})
                      </span>
                    )}
                  </div>
                );
              } else if (targetId) {
                return (
                  <div
                    key={dir}
                    className="p-1 my-1 rounded text-xs bg-surface-muted"
                  >
                    <strong>{dir}</strong> →{" "}
                    <span className="text-text-muted">Room #{targetId}</span>
                  </div>
                );
              } else {
                return (
                  <div key={dir} className="flex items-center gap-1 my-1">
                    <div className="flex-1 p-1 rounded text-xs bg-surface-muted border border-border text-text-muted">
                      <strong>{dir}</strong> → none
                    </div>
                    {onAddRoom && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="!px-1 !py-0.5"
                        onClick={() => onAddRoom(selectedRoom, dir)}
                        aria-label={`Add room to the ${dir}`}
                      >
                        +
                      </Button>
                    )}
                  </div>
                );
              }
            })}
          </div>
        </div>
      </div>

      <div className="p-3 border-t border-border flex gap-2">
        <Button variant="accent" size="md" fullWidth onClick={() => onEditRoom(selectedRoom)}>
          Edit Room
        </Button>
        <Button
          variant={confirmDelete === selectedRoom.id ? "secondary" : "danger"}
          size="md"
          fullWidth
          onClick={() => {
            if (confirmDelete === selectedRoom.id) {
              onDeleteRoom(selectedRoom.id);
              setConfirmDelete(null);
            } else {
              setConfirmDelete(selectedRoom.id);
            }
          }}
        >
          {confirmDelete === selectedRoom.id ? "Confirm Delete?" : "Delete Room"}
        </Button>
      </div>
    </>
  );
});
