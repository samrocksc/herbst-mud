import { type InputHTMLAttributes, type ReactNode } from "react";

type InputProps = {
  readonly label?: string;
  readonly error?: string;
  readonly className?: string;
  readonly icon?: ReactNode;
} & Omit<InputHTMLAttributes<HTMLInputElement>, "className">;

export function Input({ label, error, className = "", icon, ...rest }: Readonly<InputProps>) {
  return (
    <div className={`space-y-1 ${className}`}>
      {label && (
        <label className="block text-xs font-mono text-muted">{label}</label>
      )}
      <div className="flex items-center gap-2">
        {icon && <span className="text-muted">{icon}</span>}
        <input
          {...rest}
          className={`
            w-full rounded border px-3 py-2
            bg-surface-alt text-foreground
            font-mono text-sm outline-none
            placeholder-muted
            focus:border-accent
            disabled:opacity-40
            ${error ? "border-danger" : "border-border"}
          `}
        />
      </div>
      {error && <p className="text-[11px] font-mono text-danger">{error}</p>}
    </div>
  );
}