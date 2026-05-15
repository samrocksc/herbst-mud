import { Link } from '@tanstack/react-router';
import { Button } from './Button';
import type { ReactNode } from 'react';

// ─── Nav link ───────────────────────────────────────────────────────────────

type NavLinkProps = Readonly<{
  to: string
  icon: ReactNode
  label: string
  collapsed?: boolean
}>

function NavLink({ to, icon, label, collapsed }: NavLinkProps) {
  return (
    <Link
      to={to}
      activeProps={{
        className: 'bg-primary/10 text-primary border-l-4 border-primary font-semibold',
      }}
      inactiveProps={{
        className: 'text-text-muted hover:bg-surface-muted hover:text-text',
      }}
      className={[
        'flex items-center gap-3 px-3 py-2 rounded text-sm',
        'no-underline transition-colors',
        collapsed ? 'justify-center' : '',
      ].join(' ')}
      title={collapsed ? label : undefined}
    >
      <span className="flex-shrink-0">{icon}</span>
      <span
        className={[
          'whitespace-nowrap transition-opacity duration-300 min-w-0',
          collapsed ? 'opacity-0 pointer-events-none' : 'opacity-100',
        ].join(' ')}
      >
        {label}
      </span>
    </Link>
  );
}

// ─── Props ─────────────────────────────────────────────────────────────────

type SecondarySidebarProps = Readonly<{
  /** Home/back link at top of sidebar */
  homeNav: { to: string; icon: ReactNode; label: string }
  /** Primary action button — shown below home nav (e.g. "Add Room") */
  action?: {
    label: string
    onClick: () => void
    variant?: 'primary' | 'secondary' | 'danger' | 'accent'
    disabled?: boolean
  }
  /** Section title */
  title: string
  /** Stat/count shown below title */
  count?: string
  /** Secondary nav links below the action button */
  secondaryNav?: ReadonlyArray<Readonly<{ label: string; to: string; icon: ReactNode }>>
  /** Custom content slot (lists, floor selectors, etc.) */
  children?: ReactNode
}>

// ─── Component ─────────────────────────────────────────────────────────────

export function SecondarySidebar({
  homeNav,
  action,
  title,
  count,
  secondaryNav,
  children,
}: SecondarySidebarProps) {
  return (
    <div className="w-[220px] bg-surface-muted border-r border-border flex flex-col flex-shrink-0">
      {/* Home nav link */}
      <div className="p-3 border-b border-border flex flex-col gap-1">
        <NavLink to={homeNav.to} icon={homeNav.icon} label={homeNav.label} />

        {secondaryNav?.map((item) => (
          <NavLink key={item.to} to={item.to} icon={item.icon} label={item.label} />
        ))}
      </div>

      {/* Primary action button */}
      {action && (
        <div className="p-3 border-b border-border">
          <Button
            variant={action.variant ?? 'primary'}
            size="md"
            fullWidth
            onClick={action.onClick}
            disabled={action.disabled}
          >
            {action.label}
          </Button>
        </div>
      )}

      {/* Header with title + count */}
      <div className="p-3 border-b border-border">
        <h2 className="m-0 text-text text-base font-semibold">{title}</h2>
        {count && (
          <p className="text-text-muted text-xs mt-0.5 m-0">{count}</p>
        )}
      </div>

      {/* Custom content */}
      {children && (
        <div className="flex-1 overflow-y-auto">{children}</div>
      )}
    </div>
  );
}