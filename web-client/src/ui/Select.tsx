import { type SelectHTMLAttributes } from "react";

type SelectOption = {
  value: string
  label: string
};

type SelectProps = {
  label?: string;
  error?: string;
  options: SelectOption[];
  placeholder?: string;
  className?: string;
} & Omit<SelectHTMLAttributes<HTMLSelectElement>, "className">;

export function Select({ label, error, options, placeholder, className = "", ...rest }: SelectProps) {
  return (
    <div className={`space-y-1 ${className}`}>
      {label && <label className="block text-xs font-mono text-muted">{label}</label>}
        <div className="relative">
        <select
          {...rest}
          className={`
            w-full rounded border px-3 py-2 pr-8
            bg-surface-alt text-foreground
            font-mono text-sm outline-none
            appearance-none
            focus:border-accent
            disabled:opacity-40
            ${error ? "border-danger" : "border-border"}
          `}
        >
          {placeholder && (
            <option value="" disabled hidden>{placeholder}</option>
          )}
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>
        {/* chevron */}
        <svg
          className="absolute right-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none"
          fill="none"
          stroke="currentColor"
          strokeWidth={2}
          viewBox="0 0 24 24"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </div>
      {error && <p className="text-[11px] font-mono text-danger">{error}</p>}
    </div>
  );
}