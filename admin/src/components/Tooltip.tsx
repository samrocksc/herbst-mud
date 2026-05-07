import { useId, type ReactNode } from 'react'

// ─── Tooltip — CSS-only hover tooltip (no JS positioning library) ──────────

type TooltipProps = Readonly<{
  children: ReactNode
  content: string
  placement?: 'top' | 'bottom' | 'left' | 'right'
}>

export function Tooltip({ children, content, placement = 'top' }: TooltipProps) {
  const id = useId()

  return (
    <span className="tooltip-wrapper" aria-describedby={id}>
      {children}
      <span
        id={id}
        role="tooltip"
        className={`tooltip-text tooltip-${placement}`}
      >
        {content}
      </span>
    </span>
  )
}

// ─── TooltipIcon — compact ⓘ trigger for form labels ─────────────────────

export function TooltipIcon({ content }: Readonly<{ content: string }>) {
  return (
    <Tooltip content={content}>
      <span
        className="inline-flex items-center justify-center w-4 h-4 rounded-full bg-primary/20 text-primary text-[10px] font-bold cursor-help ml-1 select-none"
        aria-label="More info"
        tabIndex={0}
      >
        ⓘ
      </span>
    </Tooltip>
  )
}
