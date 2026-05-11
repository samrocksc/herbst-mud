import { useState } from 'react'
import { Button } from '../../components/Button'

type TargetType = 'room' | 'character'

/** Modal for spawning a new item instance in a room or assigned to a character. */
export function SpawnModal({ open, onClose, onSpawn, isLoading, error }: Readonly<{
  open: boolean
  onClose: () => void
  onSpawn: (targetType: TargetType, targetId: number) => void
  isLoading: boolean
  error: string | null
}>) {
  const [targetType, setTargetType] = useState<TargetType>('room')
  const [targetId, setTargetId] = useState('')

  if (!open) return null

  const canSpawn = targetId && Number(targetId) > 0

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Add Instance</h3>
          <Button variant="ghost" size="sm" onClick={onClose}>x</Button>
        </div>
        <div className="modal-body space-y-4">
          <div>
            <label className="block text-text text-sm font-medium mb-1">Assign to</label>
            <div className="flex gap-2">
              <button
                className={`px-3 py-1.5 rounded text-sm ${targetType === 'room' ? 'bg-primary text-white' : 'bg-surface-muted text-text'}`}
                onClick={() => setTargetType('room')}
              >Room</button>
              <button
                className={`px-3 py-1.5 rounded text-sm ${targetType === 'character' ? 'bg-primary text-white' : 'bg-surface-muted text-text'}`}
                onClick={() => setTargetType('character')}
              >Character</button>
            </div>
          </div>
          <div>
            <label className="block text-text text-sm font-medium mb-1">
              {targetType === 'room' ? 'Room ID' : 'Character ID'}
            </label>
            <input
              type="number"
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-text text-sm"
              placeholder={targetType === 'room' ? 'Enter room ID' : 'Enter character ID'}
              value={targetId}
              onChange={(e) => setTargetId(e.target.value)}
            />
          </div>
          {error && <div className="text-danger text-xs">{error}</div>}
        </div>
        <div className="modal-footer">
          <Button variant="primary" size="sm" disabled={!canSpawn || isLoading} onClick={() => onSpawn(targetType, Number(targetId))}>
            {isLoading ? 'Spawning...' : 'Spawn'}
          </Button>
          <Button variant="secondary" size="sm" onClick={onClose}>Cancel</Button>
        </div>
      </div>
    </div>
  )
}