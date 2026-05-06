import { useState } from 'react'
import { ALL_DIRECTIONS } from './DirectionUtils'
import { Button } from '../Button'
import { SearchableSelect } from '../SearchableSelect'
import { useRooms } from '../../hooks/useRooms'
import type { Room } from './types'

type RoomEditorProps = {
  room: Room
  onCancel: () => void
}

export function RoomEditor({
  room,
  onCancel,
}: RoomEditorProps) {
  const { updateRoom, isUpdating } = useRooms()
  const { rooms } = useRooms() // We'll need this for the room picker

  // Local state for editing
  const [form, setForm] = useState({
    name: room.name,
    description: room.description,
    exits: { ...room.exits },
    isStartingRoom: room.isStartingRoom,
  })

  const handleSave = () => {
    updateRoom({ 
      id: room.id, 
      update: {
        name: form.name,
        description: form.description,
        exits: form.exits,
        isStartingRoom: form.isStartingRoom,
        version: room.version
      } 
    })
    // Note: In a real app, we'd handle the mutation onSuccess to call onCancel
  }

  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-base font-semibold">Edit Room</h3>
        <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">
          ×
        </Button>
      </div>

      <div className="p-3 flex-1 overflow-y-auto">
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-1">Name</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-1">Description</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            rows={4}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
          />
        </div>
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-2">Exits</label>
          {ALL_DIRECTIONS.map((dir) => (
            <div key={dir} className="flex items-center gap-2 mb-2">
              <span className="w-[60px] text-text-muted text-xs">{dir}:</span>
              <SearchableSelect
                options={rooms.map(r => ({
                  id: String(r.id),
                  name: `${r.name} (ID: ${r.id})`,
                }))}
                value={form.exits[dir] ? String(form.exits[dir]) : ''}
                onChange={(val) => setForm({ 
                  ...form, 
                  exits: { ...form.exits, [dir]: val ? parseInt(val) : 0 } 
                })}
                placeholder="Pick destination room..."
              />
            </div>
          ))}
        </div>
      </div>

      <div className="p-3 border-t border-border flex gap-2">
        <Button
          variant="primary"
          size="md"
          fullWidth
          onClick={handleSave}
          disabled={isUpdating}
        >
          {isUpdating ? 'Saving...' : 'Save Changes'}
        </Button>
        <Button variant="secondary" size="md" fullWidth onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </>
  )
}
