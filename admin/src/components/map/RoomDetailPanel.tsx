import { useState } from 'react'
import type { Room, NPC, Equipment } from './types'
import { ALL_DIRECTIONS } from './DirectionUtils'

type RoomDetailPanelProps = {
  selectedRoom: Room
  rooms: Room[]
  zLevels: Map<number, number>
  npcs: NPC[]
  roomEquipment: Record<number, Equipment[]>
  onSelectRoom: (room: Room | null) => void
  onEditRoom: (room: Room) => void
  onDeleteRoom: (roomId: number) => void
}

export function RoomDetailPanel({
  selectedRoom,
  rooms,
  zLevels,
  npcs,
  roomEquipment,
  onSelectRoom,
  onEditRoom,
  onDeleteRoom,
}: RoomDetailPanelProps) {
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)

  const roomNpcs = npcs.filter(npc => npc.currentRoomId === selectedRoom.id)
  const roomItems = roomEquipment[selectedRoom.id] || []

  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-base">
          {selectedRoom.name}
          {selectedRoom.isStartingRoom && <span className="text-warning"> ⭐</span>}
        </h3>
        <button
          onClick={() => onSelectRoom(null)}
          className="bg-transparent border-none text-text-muted cursor-pointer text-xl"
        >
          ×
        </button>
      </div>
      <div className="p-3 flex-1 overflow-y-auto">
        <div className="text-text-muted text-[10px] mb-2">
          Room ID: {selectedRoom.id}
          {selectedRoom.atmosphere && ` • ${selectedRoom.atmosphere}`}
        </div>
        <div className="text-text mb-3 text-sm">{selectedRoom.description}</div>

        {roomNpcs.length > 0 && (
          <div className="mb-3">
            <strong className="text-warning text-xs">NPCs:</strong>
            <div className="mt-1">
              {roomNpcs.map(npc => (
                <div
                  key={npc.id}
                  className="p-1 bg-surface-dark rounded mb-1 text-xs text-text"
                >
                  {npc.name}{' '}
                  <span className="text-text-muted">
                    ({npc.race} {npc.class} lv.{npc.level})
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}

        {roomItems.length > 0 && (
          <div className="mb-3">
            <strong className="text-success text-xs">Items:</strong>
            <div className="mt-1">
              {roomItems.map(item => (
                <div key={item.id} className="p-1 bg-surface-dark rounded mb-1 text-xs text-text">
                  {item.name}
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="mb-3">
          <strong className="text-accent text-xs">Exits:</strong>
          <div className="mt-1">
            {ALL_DIRECTIONS.map(dir => {
              const targetId = selectedRoom.exits?.[dir]
              const targetRoom = rooms.find(r => r.id === targetId)
              const isZExit = dir === 'up' || dir === 'down'

              if (targetId && targetRoom) {
                return (
                  <div
                    key={dir}
                    onClick={() => {
                      if (isZExit) {
                        // z-level change handled by parent
                      }
                      onSelectRoom(targetRoom)
                    }}
                    className={`p-1 my-1 rounded cursor-pointer text-xs ${
                      isZExit
                        ? dir === 'up'
                          ? 'bg-warning/20 border border-warning'
                          : 'bg-success/20 border border-success'
                        : 'bg-surface-dark border-none'
                    }`}
                  >
                    <strong>{dir}</strong> → {targetRoom.name}
                    {isZExit && (
                      <span className="text-text-muted ml-1 text-[10px]">
                        (z={zLevels.get(targetId) || 0})
                      </span>
                    )}
                  </div>
                )
              } else if (targetId) {
                return (
                  <div key={dir} className="p-1 my-1 rounded text-xs bg-surface-dark border-none">
                    <strong>{dir}</strong> → <span className="text-text-muted">Room #{targetId}</span>
                  </div>
                )
              } else {
                return (
                  <div key={dir} className="flex items-center gap-2 my-1">
                    <div className="flex-1 p-1 rounded text-xs bg-surface-muted border border-border text-text-muted">
                      <strong>{dir}</strong> → none
                    </div>
                  </div>
                )
              }
            })}
          </div>
        </div>
      </div>
      <div className="p-3 border-t border-border flex gap-2">
        <button
          onClick={() => onEditRoom(selectedRoom)}
          className="flex-1 p-2 bg-accent border-2 border-black rounded text-white cursor-pointer"
        >
          Edit Room
        </button>
        <button
          onClick={() => {
            if (confirmDelete === selectedRoom.id) {
              onDeleteRoom(selectedRoom.id)
              setConfirmDelete(null)
            } else {
              setConfirmDelete(selectedRoom.id)
            }
          }}
          className={`flex-1 p-2 border-none rounded text-white cursor-pointer ${
            confirmDelete === selectedRoom.id ? 'bg-warning hover:bg-warning/80' : 'bg-danger hover:bg-danger-hover'
          }`}
        >
          {confirmDelete === selectedRoom.id ? 'Confirm Delete?' : 'Delete Room'}
        </button>
      </div>
    </>
  )
}
