import type { InputHTMLAttributes } from 'react';
import { FieldLabel, INPUT_CLASS } from './FieldLabel';

type NumberFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
  disabled?: boolean
  tooltip?: string
}> & Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'value' | 'className' | 'id' | 'type'>

export function NumberField({ label, id, placeholder, value, onChange, min, max, disabled, tooltip, ...rest }: NumberFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-');
  return (
    <div>
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <input
        id={fieldId}
        type="text"
        inputMode="numeric"
        value={Number.isFinite(value) ? String(value) : ''}
        onChange={(e) => {
          const raw = e.target.value;
          if (raw === '' || raw === '-') { onChange(0); return; }
          const parsed = parseInt(raw, 10);
          if (!isNaN(parsed)) onChange(parsed);
        }}
        placeholder={placeholder}
        min={min}
        max={max}
        disabled={disabled}
        className={`${INPUT_CLASS} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        {...rest}
      />
    </div>
  );
}