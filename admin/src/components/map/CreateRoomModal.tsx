import { useState } from 'react'
import { Modal } from '../Modal'
import { Button } from '../Button'

type CreateRoomModalProps = {
  isOpen: boolean
  onClose: () => void
  onCreate: (input: { name: string; description: string }) => void
  creating: boolean
}

export function CreateRoomModal({ isOpen, onClose, onCreate, creating }: CreateRoomModalProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const handleSubmit = () => {
    if (!name.trim()) return
    onCreate({ name: name.trim(), description: description.trim() })
    setName('')
    setDescription('')
  }

  if (!isOpen) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Create New Room">
      <div className="mb-4">
        <label className="text-text-muted text-xs block mb-1">Room Name</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter room name"
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>
      <div className="mb-4">
        <label className="text-text-muted text-xs block mb-1">Description</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Enter room description"
          rows={4}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
        />
      </div>
      <div className="flex gap-2">
        <Button variant="primary" size="md" fullWidth onClick={handleSubmit} disabled={creating || !name.trim()}>
          {creating ? 'Creating...' : 'Create Room'}
        </Button>
        <Button variant="secondary" size="md" fullWidth onClick={onClose}>
          Cancel
        </Button>
      </div>
    </Modal>
  )
}