export type ModalProps = Readonly<{
  isOpen: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
}>