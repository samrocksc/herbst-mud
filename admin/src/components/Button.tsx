import type { ButtonHTMLAttributes } from 'react';

// ─── Variant map ─────────────────────────────────────────────────────────────

const variantClasses = {
  primary:
    'bg-primary text-white hover:bg-primary-hover border-primary hover:border-primary-hover',
  secondary:
    'bg-transparent text-text border border-border hover:bg-surface-muted hover:border-text-muted',
  danger:
    'bg-danger text-white hover:bg-danger-hover border-danger hover:border-danger-hover',
  ghost:
    'bg-transparent text-text-muted hover:bg-surface-muted hover:text-text border-transparent',
  outline:
    'bg-transparent text-text border-border hover:border-primary hover:text-primary',
  accent:
    'bg-accent text-white hover:bg-accent-hover border-accent hover:border-accent-hover',
  success:
    'bg-success text-white hover:opacity-90 border-success',
} as const;

// ─── Size map ────────────────────────────────────────────────────────────────

const sizeClasses = {
  sm: 'px-2.5 py-1 text-xs',
  md: 'px-4 py-2 text-sm',
  lg: 'px-6 py-3 text-base',
} as const;

// ─── Component ───────────────────────────────────────────────────────────────

type ButtonProps = Readonly<{
  variant?: keyof typeof variantClasses
  size?: keyof typeof sizeClasses
  fullWidth?: boolean
}> &
  ButtonHTMLAttributes<HTMLButtonElement>

export function Button({
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  className = '',
  disabled,
  ...rest
}: ButtonProps) {
  const classes = [
    // Base — no global reset interference, all Tailwind
    'inline-flex items-center justify-center gap-2',
    'font-medium rounded-md border-2',
    'cursor-pointer select-none',
    'transition-colors duration-150',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-1',
    'disabled:opacity-50 disabled:cursor-not-allowed',
    // Variant
    variantClasses[variant],
    // Size
    sizeClasses[size],
    // Width
    fullWidth ? 'w-full' : '',
    // User classes
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return <button className={classes} disabled={disabled} {...rest} />;
}
