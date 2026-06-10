/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types */
import { useState, memo } from "react";
import type { Room } from "./types";
import { ALL_DIRECTIONS } from "./DirectionUtils";
import { Button } from "../Button";
import { NPCInstanceManager } from "./NPCInstanceManager";
import { ItemInstanceManager } from "./ItemInstanceManager";
import { NewRoomModal } from "./NewRoomModal";
import { RoomDeleteModal } from "./RoomDeleteModal";

type RoomDetailPanelProps = {
  selectedRoom: Room;
  rooms: Room[];
  zLevels: Map<number, number>;
  onSelectRoom: (room: Room | null) => void;
  onEditRoom: (room: Room) => void;
  onDeleteRoom: (roomId: number) => void;
  onRequestDeleteRoom?: (roomId: number) => void;
  onAddRoom?: (room: Room, dir: string) => void;
  addRoomModal?: {
    open: boolean;
    fromRoom: Room | null;
    dir: string | null;
  };
  onConfirmAddRoom?: (input: { name: string; description: string }) => void;
  onCancelAddRoom?: () => void;
  isAddingRoom?: boolean;
  deleteRoomModalOpen?: boolean;
  onConfirmDeleteRoom?: () => void;
  onCancelDeleteRoom?: () => void;
  isDeletingRoom?: boolean;
};

export const RoomDetailPanel = memo(function RoomDetailPanel({
  selectedRoom,
  rooms,
  zLevels,
  onSelectRoom,
  onEditRoom,
  onDeleteRoom,
  onRequestDeleteRoom,
  onAddRoom,
  addRoomModal,
  onConfirmAddRoom,
  onCancelAddRoom,
  isAddingRoom,
  deleteRoomModalOpen,
  onConfirmDeleteRoom,
  onCancelDeleteRoom,
  isDeletingRoom,
}: RoomDetailPanelProps) {

  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center flex-shrink-0">
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

      <div className="p-3 flex-1 min-h-0 overflow-y-auto">
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
                    className={[
                      "p-1 my-1 rounded text-xs flex items-center gap-2 transition-colors hover:bg-surface",
                      isZExit
                        ? dir === "up"
                          ? "bg-warning/20 border border-warning"
                          : "bg-success/20 border border-success"
                        : "bg-surface-muted",
                    ].join(" ")}
                  >
                    <Button
                      variant="ghost"
                      size="sm"
                      className="!px-1 !py-0.5 flex-1 !justify-start text-left"
                      onClick={() => onSelectRoom(targetRoom)}
                      aria-label={`Navigate to ${targetRoom.name} via ${dir} exit`}
                    >
                      <strong>{dir}</strong> → {targetRoom.name}
                      {isZExit && (
                        <span className="text-text-muted ml-1 text-[10px]">
                          (z={zLevels.get(targetId) || 0})
                        </span>
                      )}
                    </Button>
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

      <div className="p-3 border-t border-border flex gap-2 flex-shrink-0">
        <Button variant="accent" size="md" fullWidth onClick={() => onEditRoom(selectedRoom)}>
          Edit Room
        </Button>
        {onRequestDeleteRoom && (
          <Button
            variant="danger"
            size="md"
            fullWidth
            onClick={() => onRequestDeleteRoom(selectedRoom.id)}
          >
            Delete Room
          </Button>
        )}
      </div>

      {onRequestDeleteRoom && deleteRoomModalOpen && onConfirmDeleteRoom && onCancelDeleteRoom && (
        <RoomDeleteModal
          open={deleteRoomModalOpen}
          room={selectedRoom}
          onConfirm={onConfirmDeleteRoom}
          onCancel={onCancelDeleteRoom}
          isLoading={isDeletingRoom}
        />
      )}

      {addRoomModal && onConfirmAddRoom && onCancelAddRoom && (
        <NewRoomModal
          open={addRoomModal.open}
          parentRoom={addRoomModal.fromRoom}
          direction={addRoomModal.dir}
          onConfirm={onConfirmAddRoom}
          onCancel={onCancelAddRoom}
          isLoading={isAddingRoom}
        />
      )}
    </>
  );
});
