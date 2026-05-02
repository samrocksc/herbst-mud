import { ALL_DIRECTIONS } from './DirectionUtils'

type RoomEditorProps = {
  editForm: {
    name: string
    description: string
    exits: Record<string, string>
  }
  setEditForm: (form: { name: string; description: string; exits: Record<string, string> }) => void
  onSave: () => void
  onCancel: () => void
  saving: boolean
}

export function RoomEditor({ editForm, setEditForm, onSave, onCancel, saving }: RoomEditorProps) {
  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-base">Edit Room</h3>
        <button
          onClick={onCancel}
          className="bg-transparent border-none text-text-muted cursor-pointer text-xl"
        >
          ×
        </button>
      </div>
      <div className="p-3 flex-1 overflow-y-auto">
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-1">Name</label>
          <input
            type="text"
            value={editForm.name}
            onChange={e => setEditForm({ ...editForm, name: e.target.value })}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-1">Description</label>
          <textarea
            value={editForm.description}
            onChange={e => setEditForm({ ...editForm, description: e.target.value })}
            rows={4}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
          />
        </div>
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-2">Exits</label>
          {ALL_DIRECTIONS.map(dir => (
            <div key={dir} className="flex items-center gap-2 mb-1">
              <span className="w-[50px] text-text-muted text-xs">{dir}:</span>
              <input
                type="text"
                value={editForm.exits[dir] || ''}
                onChange={e =>
                  setEditForm({ ...editForm, exits: { ...editForm.exits, [dir]: e.target.value } })
                }
                placeholder="room id"
                className="flex-1 p-1 bg-surface border border-border rounded text-text text-xs"
              />
            </div>
          ))}
        </div>
      </div>
      <div className="p-3 border-t border-border">
        <button
          onClick={onSave}
          disabled={saving}
          className="w-full p-2 bg-primary border-2 border-black rounded text-white cursor-pointer mb-2 disabled:opacity-70"
        >
          {saving ? 'Saving...' : 'Save Changes'}
        </button>
        <button
          onClick={onCancel}
          className="w-full p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
        >
          Cancel
        </button>
      </div>
    </>
  )
}
