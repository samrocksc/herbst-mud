import { Modal } from '../Modal'

type CreateRoomModalProps = {
  isOpen: boolean
  onClose: () => void
  newRoomForm: { name: string; description: string }
  setNewRoomForm: (form: { name: string; description: string }) => void
  onCreate: () => void
  creating: boolean
}

export function CreateRoomModal({
  isOpen,
  onClose,
  newRoomForm,
  setNewRoomForm,
  onCreate,
  creating,
}: CreateRoomModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Create New Room">
      <div className="mb-4">
        <label className="text-text-muted text-xs block mb-1">Room Name</label>
        <input
          type="text"
          value={newRoomForm.name}
          onChange={e => setNewRoomForm({ ...newRoomForm, name: e.target.value })}
          placeholder="Enter room name"
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>
      <div className="mb-4">
        <label className="text-text-muted text-xs block mb-1">Description</label>
        <textarea
          value={newRoomForm.description}
          onChange={e => setNewRoomForm({ ...newRoomForm, description: e.target.value })}
          placeholder="Enter room description"
          rows={4}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
        />
      </div>
      <div className="flex gap-2">
        <button
          onClick={onCreate}
          disabled={creating}
          className="flex-1 p-2 bg-primary border-2 border-black rounded text-white cursor-pointer disabled:opacity-70"
        >
          {creating ? 'Creating...' : 'Create Room'}
        </button>
        <button
          onClick={onClose}
          className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
        >
          Cancel
        </button>
      </div>
    </Modal>
  )
}
