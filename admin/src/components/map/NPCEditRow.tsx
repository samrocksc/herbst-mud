import { Button } from '../Button'
import { NumberField } from '../fields/NumberField'
import { FormError } from '../fields/FormError'
import type { NPCInstanceView, EditFormData } from './NPCInstanceManager'

type NPCEditRowProps = Readonly<{
  inst: NPCInstanceView
  editForm: EditFormData
  setEditForm: React.Dispatch<React.SetStateAction<EditFormData>>
  onSave: () => void
  onCancel: () => void
  isPending: boolean
  error: Error | null
}>

export function NPCEditRow({ inst, editForm, setEditForm, onSave, onCancel, isPending, error }: NPCEditRowProps) {
  return (
    <div className="space-y-1">
      <div className="font-medium">{inst.name}</div>
      {error && <FormError message={(error as Error)?.message || 'Update failed'} />}
      <div className="flex gap-1">
        <NumberField label="Level" value={editForm.level} onChange={(v) => setEditForm((f) => ({ ...f, level: v }))} disabled={isPending} />
        <NumberField label="HP" value={editForm.hitpoints} onChange={(v) => setEditForm((f) => ({ ...f, hitpoints: v }))} disabled={isPending} />
      </div>
      <NumberField label="Room ID" value={editForm.room_id} onChange={(v) => setEditForm((f) => ({ ...f, room_id: v }))} disabled={isPending} />
      <NumberField label="Starting Room" value={editForm.starting_room_id} onChange={(v) => setEditForm((f) => ({ ...f, starting_room_id: v }))} disabled={isPending} />
      <div className="flex gap-1">
        <Button variant="primary" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={onSave} disabled={isPending}>
          {isPending ? 'Saving...' : 'Save'}
        </Button>
        <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={onCancel}>Cancel</Button>
      </div>
    </div>
  )
}

type NPCInstanceRowProps = Readonly<{
  inst: NPCInstanceView
  editingId: number | null
  confirmDeleteId: number | null
  editForm: EditFormData
  setEditForm: React.Dispatch<React.SetStateAction<EditFormData>>
  startEdit: (inst: NPCInstanceView) => void
  handleUpdate: () => void
  handleDelete: (id: number) => void
  setEditingId: (id: number | null) => void
  setConfirmDeleteId: (id: number | null) => void
  isUpdatePending: boolean
  updateError: Error | null
}>

export function NPCInstanceRow({ inst, editingId, confirmDeleteId, editForm, setEditForm,
  startEdit, handleUpdate, handleDelete, setEditingId, setConfirmDeleteId, isUpdatePending, updateError }: NPCInstanceRowProps) {
  const isEditing = editingId === inst.id
  const isConfirmDelete = confirmDeleteId === inst.id && !isEditing
  return (
    <div className="p-1 bg-surface-muted rounded text-xs text-text">
      {isEditing ? (
        <NPCEditRow inst={inst} editForm={editForm} setEditForm={setEditForm}
          onSave={handleUpdate} onCancel={() => setEditingId(null)}
          isPending={isUpdatePending} error={updateError} />
      ) : (
        <div className="flex justify-between items-center">
          <div><span className="font-medium">{inst.name}</span>{' '}<span className="text-text-muted">{inst.race} lv.{inst.level} HP:{inst.hitpoints}/{inst.max_hitpoints}</span></div>
          <div className="flex gap-0.5">
            <Button variant="ghost" size="sm" className="!px-0.5 !py-0" onClick={() => startEdit(inst)} aria-label={`Edit ${inst.name}`}>✏️</Button>
            <Button variant={isConfirmDelete ? 'secondary' : 'ghost'} size="sm" className="!px-0.5 !py-0"
              onClick={() => isConfirmDelete ? handleDelete(inst.id) : setConfirmDeleteId(inst.id)}
              aria-label={isConfirmDelete ? `Confirm delete ${inst.name}` : `Delete ${inst.name}`}>
              {isConfirmDelete ? '❓' : '🗑'}
            </Button>
          </div>
        </div>
      )}
      {isConfirmDelete && (
        <div className="mt-1 text-[10px] text-text">Confirm delete?{' '}
          <Button variant="danger" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => handleDelete(inst.id)}>Yes</Button>{' '}
          <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => setConfirmDeleteId(null)}>No</Button>
        </div>
      )}
    </div>
  )
}