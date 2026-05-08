import { Button } from '../Button'
import type { ItemInstanceView } from './types'

type ItemInstanceRowProps = Readonly<{
  inst: ItemInstanceView
  confirmDeleteId: number | null
  onEdit: () => void
  onDelete: () => void
}>

export function ItemInstanceRow({ inst, confirmDeleteId, onEdit, onDelete }: ItemInstanceRowProps) {
  return (
    <div className="flex justify-between items-center">
      <div>
        <span className="font-medium">{inst.name}</span>{' '}
        <span className="text-text-muted">{inst.itemType} lv.{inst.level}</span>
        {!inst.isVisible && <span className="text-warning ml-1 text-[10px]">(hidden)</span>}
        {inst.isImmovable && <span className="text-danger ml-1 text-[10px]">(immovable)</span>}
      </div>
      <div className="flex gap-0.5">
        <Button variant="ghost" size="sm" className="!px-0.5 !py-0" onClick={onEdit} aria-label={`Edit ${inst.name}`}>✏️</Button>
        <Button
          variant={confirmDeleteId === inst.id ? 'secondary' : 'ghost'}
          size="sm" className="!px-0.5 !py-0"
          onClick={onDelete}
          aria-label={confirmDeleteId === inst.id ? `Confirm delete ${inst.name}` : `Delete ${inst.name}`}
        >
          {confirmDeleteId === inst.id ? '❓' : '🗑'}
        </Button>
      </div>
    </div>
  )
}