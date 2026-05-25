import { type ReactNode, type ButtonHTMLAttributes } from "react";

type ButtonVariant = "primary" | "secondary" | "ghost" | "danger" | "success" | "warning" | "info";
type ButtonSize = "sm" | "md" | "lg";

type ButtonProps = {
  readonly children: ReactNode;
  readonly variant?: ButtonVariant;
  readonly size?: ButtonSize;
  readonly fullWidth?: boolean;
  readonly className?: string;
} & Omit<ButtonHTMLAttributes<HTMLButtonElement>, "className">;

const variantMap: Record<ButtonVariant, string> = {
  primary:   "bg-accent text-background hover:opacity-90 active:opacity-80",
  secondary: "bg-surface text-foreground border border-border hover:bg-surface-alt active:bg-border",
  ghost:     "bg-transparent text-muted hover:text-foreground hover:bg-surface active:bg-surface-alt",
  danger:    "bg-danger text-background hover:opacity-90 active:opacity-80",
  success:   "bg-success text-background hover:opacity-90 active:opacity-80",
  warning:   "bg-warning text-background hover:opacity-90 active:opacity-80",
  info:      "bg-info text-background hover:opacity-90 active:opacity-80",
};

const sizeMap: Record<ButtonSize, string> = {
  sm: "px-2 py-1 text-[11px]",
  md: "px-3 py-1.5 text-xs",
  lg: "px-4 py-2 text-sm",
};

export function Button({
  children,
  variant = "primary",
  size = "md",
  fullWidth = false,
  className = "",
  disabled,
  ...rest
}: Readonly<ButtonProps>) {
  return (
    <button
      type="button"
      disabled={disabled}
      className={`
        inline-flex items-center justify-center gap-1 rounded
        font-mono font-bold transition-opacity
        disabled:opacity-40 disabled:cursor-not-allowed
        ${variantMap[variant]}
        ${sizeMap[size]}
        ${fullWidth ? "w-full" : ""}
        ${className}
      `}
      {...rest}
    >
      {children}
    </button>
  );
}

/** Action chip used inside game panels (exits, characters, items, hotkeys) */
type ActionChipProps = {
  readonly children: ReactNode;
  readonly onClick?: () => void;
  readonly className?: string;
};

export function ActionChip({
  children,
  onClick,
  className = "",
}: Readonly<ActionChipProps>) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`
        inline-flex items-center gap-1 px-2 py-1 rounded
        border border-border text-[11px] font-mono
        hover:bg-surface-alt active:bg-accent active:text-background
        transition-colors cursor-pointer
        ${className}
      `}
    >
      {children}
    </button>
  );
}