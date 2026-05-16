 
import type { SelectHTMLAttributes } from "react";
import { FieldLabel, INPUT_CLASS } from "./FieldLabel";

type SelectOption = Readonly<{ value: string; label: string }>

type SelectFieldProps = Readonly<{
  label: string
  id?: string
  value: string
  onChange: (value: string) => void
  options: SelectOption[]
  placeholder?: string
  disabled?: boolean
  tooltip?: string
}> & Omit<SelectHTMLAttributes<HTMLSelectElement>, "onChange" | "value" | "className" | "id">

export function SelectField({ label, id, value, onChange, options, placeholder, disabled, tooltip, ...rest }: SelectFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, "-");
  return (
    <div>
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <select
        id={fieldId}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        className={`${INPUT_CLASS} ${disabled ? "opacity-50 cursor-not-allowed" : ""}`}
        {...rest}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>{opt.label}</option>
        ))}
      </select>
    </div>
  );
}