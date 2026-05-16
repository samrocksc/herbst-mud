 
import type { TextareaHTMLAttributes } from "react";
import { FieldLabel, INPUT_CLASS } from "./FieldLabel";

type TextareaFieldProps = Readonly<{
  label: string
  id?: string
  placeholder?: string
  value: string
  onChange: (value: string) => void
  rows?: number
  disabled?: boolean
  tooltip?: string
}> & Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, "onChange" | "value" | "className" | "id">

export function TextareaField({ label, id, placeholder, value, onChange, rows = 3, disabled, tooltip, ...rest }: TextareaFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, "-");
  return (
    <div>
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <textarea
        id={fieldId}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        rows={rows}
        disabled={disabled}
        className={`${INPUT_CLASS} resize-y ${disabled ? "opacity-50 cursor-not-allowed" : ""}`}
        {...rest}
      />
    </div>
  );
}