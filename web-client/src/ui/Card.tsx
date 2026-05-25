import { type ReactNode } from "react";

type CardProps = {
  children: ReactNode
  className?: string
  onClick?: () => void
  hover?: boolean
};

export function Card({ children, className = "", onClick, hover = false }: CardProps) {
  const base = `
    rounded border border-border bg-surface
    font-mono transition-colors
  `;
  const hoverClasses = hover
    ? "cursor-pointer hover:bg-surface-alt hover:opacity-80 active:opacity-60"
    : "";

  if (onClick) {
    return (
      <button
        type="button"
        onClick={onClick}
        className={`${base} ${hoverClasses} text-left w-full ${className}`}
      >
        {children}
      </button>
    );
  }

  return <div className={`${base} ${className}`}>{children}</div>;
}

type CardHeaderProps = {
  title: string
  subtitle?: string
  className?: string
};

export function CardHeader({ title, subtitle, className = "" }: CardHeaderProps) {
  return (
    <div className={`px-4 py-3 border-b border-border ${className}`}>
      <h3 className="font-bold text-sm">{title}</h3>
      {subtitle && <p className="text-[11px] text-muted mt-0.5">{subtitle}</p>}
    </div>
  );
}

type CardBodyProps = {
  children: ReactNode
  className?: string
};

export function CardBody({ children, className = "" }: CardBodyProps) {
  return <div className={`px-4 py-3 ${className}`}>{children}</div>;
}

type CardFooterProps = {
  children: ReactNode
  className?: string
};

export function CardFooter({ children, className = "" }: CardFooterProps) {
  return (
    <div className={`px-4 py-2 border-t border-border flex items-center gap-2 ${className}`}>
      {children}
    </div>
  );
}