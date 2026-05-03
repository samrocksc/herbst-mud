import { Modal } from '../Modal'
import { Button } from '../Button'

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
          onChange={(e) =>
            setNewRoomForm({ ...newRoomForm, name: e.target.value })
          }
          placeholder="Enter room name"
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>
      <div className="mb-4">
        <label className="text-text-muted text-xs block mb-1">Description</label>
        <textarea
          value={newRoomForm.description}
          onChange={(e) =>
            setNewRoomForm({ ...newRoomForm, description: e.target.value })
          }
          placeholder="Enter room description"
          rows={4}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
        />
      </div>
      <div className="flex gap-2">
        <Button
          variant="primary"
          size="md"
          fullWidth
          onClick={onCreate}
          disabled={creating}
        >
          {creating ? 'Creating...' : 'Create Room'}
        </Button>
        <Button variant="secondary" size="md" fullWidth onClick={onClose}>
          Cancel
        </Button>
      </div>
    </Modal>
  )
}
