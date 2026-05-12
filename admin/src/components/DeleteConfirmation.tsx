import { Modal } from './Modal'
import { Button } from './Button'

type DeleteConfirmationProps = Readonly<{
  open: boolean
  title?: string
  message: string
  onConfirm: () => void
  onCancel: () => void
  isLoading?: boolean
}>

export function DeleteConfirmation({ open, title = 'Confirm Delete', message, onConfirm, onCancel, isLoading }: DeleteConfirmationProps) {
  return (
    <Modal isOpen={open} onClose={onCancel} title={title}>
      <p className="text-text mb-4">{message}</p>
      <div className="flex gap-2 justify-end">
        <Button variant="secondary" onClick={onCancel} disabled={isLoading} type="button">Cancel</Button>
        <Button variant="danger" onClick={onConfirm} disabled={isLoading} type="button">
          {isLoading ? 'Deleting…' : 'Delete'}
        </Button>
      </div>
    </Modal>
  )
}