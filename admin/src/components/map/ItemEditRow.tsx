import { FormField } from '../fields/FormField';
import { NumberField } from '../fields/NumberField';
import { FormError } from '../fields/FormError';
import { Button } from '../Button';
import type { ItemInstanceView, EditFormData } from './types';

type ItemEditRowProps = Readonly<{
  inst: ItemInstanceView
  editForm: Partial<EditFormData>
  setEditForm: React.Dispatch<React.SetStateAction<Partial<EditFormData>>>
  onSave: () => void
  onCancel: () => void
  isPending: boolean
  error: Error | null
}>

export function ItemEditRow({
  inst, editForm, setEditForm, onSave, onCancel, isPending, error,
}: ItemEditRowProps) {
  return (
    <div className="space-y-1">
      <div className="font-medium">{inst.name}</div>
      {error && <FormError message={error.message || 'Failed to update item'} />}
      <div className="flex gap-1">
        <div className="flex-1">
          <FormField
            label="Name"
            value={editForm.name ?? inst.name}
            onChange={(v) => setEditForm((f) => ({ ...f, name: v }))}
          />
        </div>
        <div className="flex-1">
          <NumberField
            label="Level"
            value={editForm.level ?? inst.level}
            onChange={(v) => setEditForm((f) => ({ ...f, level: v }))}
          />
        </div>
      </div>
      <div className="flex gap-1">
        <div className="flex-1">
          <FormField
            label="Slot"
            value={editForm.slot ?? inst.slot}
            onChange={(v) => setEditForm((f) => ({ ...f, slot: v }))}
          />
        </div>
        <div className="flex-1">
          <NumberField
            label="Weight"
            value={editForm.weight ?? inst.weight}
            onChange={(v) => setEditForm((f) => ({ ...f, weight: v }))}
          />
        </div>
      </div>
      <div className="flex gap-1">
        <Button variant="primary" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={onSave} disabled={isPending}>
          {isPending ? 'Saving...' : 'Save'}
        </Button>
        <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </div>
  );
}