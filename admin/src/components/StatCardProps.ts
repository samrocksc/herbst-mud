export type StatCardProps = Readonly<{
  label: string
  value: string | number
  accent?: 'primary' | 'warning' | 'accent' | 'secondary' | 'success' | 'danger'
  loading?: boolean
}>