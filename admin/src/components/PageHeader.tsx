import { Link } from '@tanstack/react-router'

interface PageHeaderProps {
  title: string
  showBack?: boolean
  backTo?: string
  backLabel?: string
  actions?: React.ReactNode
}

export function PageHeader({ title, showBack, backTo, backLabel = '← Dashboard', actions }: PageHeaderProps) {
  return (
    <div className="flex justify-between items-center mb-6">
      <div className="flex items-center gap-3">
        {showBack && backTo && (
          <Link
            to={backTo as any}
            className="text-primary no-underline px-3 py-1.5 rounded border border-border hover:border-primary transition-colors text-sm font-medium"
          >
            {backLabel}
          </Link>
        )}
        <h1 className="m-0 text-xl font-bold text-text">{title}</h1>
      </div>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  )
}
