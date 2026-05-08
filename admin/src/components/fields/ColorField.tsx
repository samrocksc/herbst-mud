import { FieldLabel } from './FieldLabel'

type ColorFieldProps = Readonly<{
  label: string
  id?: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
}>

const DEFAULT_COLOR = 'var(--color-tag-default)'

export function ColorField({ label, id, value, onChange, placeholder }: ColorFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <div className="flex items-center gap-2">
        <input
          id={fieldId}
          type="color"
          value={value || DEFAULT_COLOR}
          onChange={(e) => onChange(e.target.value)}
          className="w-10 h-8 p-0.5 cursor-pointer"
        />
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder ?? 'CSS color / hex'}
          pattern="^#[0-9a-fA-F]{6}$"
          className="w-28 p-2 bg-surface border border-border rounded text-text text-sm"
        />
        <span
          className="inline-block w-3 h-3 rounded-full shrink-0"
          style={{ backgroundColor: value || DEFAULT_COLOR }}
        />
      </div>
    </div>
  )
}