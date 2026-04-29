interface StatCardProps {
  label: string
  value: string | number
  accent?: 'primary' | 'warning' | 'accent' | 'secondary' | 'success' | 'danger'
  loading?: boolean
}

const accentMap = {
  primary: 'text-primary',
  warning: 'text-warning',
  accent: 'text-accent',
  secondary: 'text-secondary',
  success: 'text-success',
  danger: 'text-danger',
} as const

export function StatCard({ label, value, accent = 'primary', loading = false }: StatCardProps) {
  return (
    <div className="bg-surface-muted rounded-lg p-6 text-center">
      <div className={`text-2xl font-bold ${accentMap[accent]}`}>
        {loading ? '--' : value}
      </div>
      <div className="text-text-muted text-sm mt-1">{label}</div>
    </div>
  )
}
