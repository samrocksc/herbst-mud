import type { ReactNode } from 'react'
import { TooltipIcon } from '../Tooltip'

type FieldLabelProps = Readonly<{
  htmlFor?: string
  children: ReactNode
  tooltip?: string
}>

export const INPUT_CLASS = 'w-full p-2 bg-surface border border-border rounded text-text text-sm'

export function FieldLabel({ htmlFor, children, tooltip }: FieldLabelProps) {
  return (
    <label htmlFor={htmlFor} className="text-text-muted text-xs block mb-1">
      {children}
      {tooltip && <TooltipIcon content={tooltip} />}
    </label>
  )
}