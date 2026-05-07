import { Button } from '../../components/Button'

/** Modal for spawning a new item instance in a room. */
export function SpawnModal({ open, onClose, spawnRoomId, setSpawnRoomId, onSpawn, isLoading, error }: Readonly<{
  open: boolean
  onClose: () => void
  spawnRoomId: string
  setSpawnRoomId: (v: string) => void
  onSpawn: () => void
  isLoading: boolean
  error: string | null
}>) {
  if (!open) return null
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header"><h3>Add Instance</h3><Button variant="ghost" size="sm" onClick={onClose}>x</Button></div>
        <div className="modal-body space-y-4">
          <div><label className="block text-text text-sm font-medium mb-1">Room ID</label>
            <input type="number" className="w-full bg-surface border border-border rounded-md px-3 py-2 text-text text-sm"
              placeholder="Enter room ID" value={spawnRoomId} onChange={(e) => setSpawnRoomId(e.target.value)} /></div>
          {error && <div className="text-danger text-xs">{error}</div>}
        </div>
        <div className="modal-footer">
          <Button variant="primary" size="sm" disabled={!spawnRoomId || isLoading} onClick={onSpawn}>{isLoading ? 'Spawning...' : 'Spawn'}</Button>
          <Button variant="secondary" size="sm" onClick={onClose}>Cancel</Button>
        </div>
      </div>
    </div>
  )
}