import type { InputHTMLAttributes, SelectHTMLAttributes, TextareaHTMLAttributes, ReactNode } from 'react'

// ─── Shared input style (the 56-copy-paste winner) ───────────────────────

const inputClass = 'w-full p-2 bg-surface border border-border rounded text-text text-sm'

// ─── Label component ─────────────────────────────────────────────────────

type FieldLabelProps = Readonly<{
  htmlFor?: string
  children: ReactNode
}>

function FieldLabel({ htmlFor, children }: FieldLabelProps) {
  return (
    <label htmlFor={htmlFor} className="text-text-muted text-xs block mb-1">
      {children}
    </label>
  )
}

// ─── FormField — text input + label ──────────────────────────────────────

type FormFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: string
  onChange: (value: string) => void
  disabled?: boolean
} & Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'value' | 'className' | 'id'>>

export function FormField({ label, id, placeholder, value, onChange, disabled, ...rest }: FormFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <input
        id={fieldId}
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        disabled={disabled}
        className={`${inputClass} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      />
    </div>
  )
}

// ─── NumberField — numeric input (no spinners) + label ───────────────────

type NumberFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
  disabled?: boolean
} & Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'value' | 'className' | 'id' | 'type'>>

export function NumberField({ label, id, placeholder, value, onChange, min, max, disabled, ...rest }: NumberFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <input
        id={fieldId}
        type="text"
        inputMode="numeric"
        value={Number.isFinite(value) ? String(value) : ''}
        onChange={(e) => {
          const raw = e.target.value
          if (raw === '' || raw === '-') {
            onChange(0)
            return
          }
          const parsed = parseInt(raw, 10)
          if (!isNaN(parsed)) onChange(parsed)
        }}
        placeholder={placeholder}
        min={min}
        max={max}
        disabled={disabled}
        className={`${inputClass} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      />
    </div>
  )
}

// ─── TextareaField — textarea + label ────────────────────────────────────

type TextareaFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: string
  onChange: (value: string) => void
  rows?: number
  disabled?: boolean
} & Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, 'onChange' | 'value' | 'className' | 'id'>>

export function TextareaField({ label, id, placeholder, value, onChange, rows = 3, disabled, ...rest }: TextareaFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <textarea
        id={fieldId}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        rows={rows}
        disabled={disabled}
        className={`${inputClass} resize-y ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      />
    </div>
  )
}

// ─── SelectField — select dropdown + label ───────────────────────────────

type SelectOption = Readonly<{
  value: string
  label: string
}>

type SelectFieldProps = Readonly<{
  label: string
  id?: string
  value: string
  onChange: (value: string) => void
  options: SelectOption[]
  placeholder?: string
  disabled?: boolean
} & Omit<SelectHTMLAttributes<HTMLSelectElement>, 'onChange' | 'value' | 'className' | 'id'>>

export function SelectField({ label, id, value, onChange, options, placeholder, disabled, ...rest }: SelectFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <select
        id={fieldId}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        className={`${inputClass} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
    </div>
  )
}

// ─── CheckboxField — checkbox + label (row layout) ───────────────────────

type CheckboxFieldProps = Readonly<{
  label: string
  id?: string
  checked: boolean
  onChange: (checked: boolean) => void
  disabled?: boolean
}>

export function CheckboxField({ label, id, checked, onChange, disabled }: CheckboxFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
        <input
          id={fieldId}
          type="checkbox"
          checked={checked}
          onChange={(e) => onChange(e.target.checked)}
          disabled={disabled}
        />
        {label}
      </label>
    </div>
  )
}

// ─── Error banner for forms ──────────────────────────────────────────────

type FormErrorProps = Readonly<{
  message: string
}>

export function FormError({ message }: FormErrorProps) {
  return (
    <div className="p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
      {message}
    </div>
  )
}
