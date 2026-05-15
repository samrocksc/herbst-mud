type CheckboxFieldProps = Readonly<{
  label: string
  id?: string
  checked: boolean
  onChange: (checked: boolean) => void
  disabled?: boolean
}>

export function CheckboxField({ label, id, checked, onChange, disabled }: CheckboxFieldProps) {
  const fieldId = id ?? label.toLowerCase().replace(/\s+/g, '-');
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
  );
}