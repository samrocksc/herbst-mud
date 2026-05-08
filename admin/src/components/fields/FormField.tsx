import type { InputHTMLAttributes } from 'react'
import { FieldLabel, INPUT_CLASS } from './FieldLabel'

type FormFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: string
  onChange: (value: string) => void
  disabled?: boolean
  tooltip?: string
}> & Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'value' | 'className' | 'id'>

export function FormField({ label, id, placeholder, value, onChange, disabled, tooltip, ...rest }: FormFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <input
        id={fieldId}
        type={rest.type ?? 'text'}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        disabled={disabled}
        className={`${INPUT_CLASS} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      />
    </div>
  )
}