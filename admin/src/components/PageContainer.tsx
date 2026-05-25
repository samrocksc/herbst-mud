import type { ReactNode } from "react";

type PageContainerProps = Readonly<{
  children: ReactNode
  className?: string
  fullWidth?: boolean
}>

/**
 * Responsive page wrapper used on every management page.
 * Mobile: tight padding, no max-width constraint.
 * Desktop: comfortable padding with optional max-width.
 */
export function PageContainer({ children, className = "", fullWidth = false }: PageContainerProps) {
  return (
    <div
      className={[
        "min-h-full bg-surface",
        "px-3 py-4 sm:px-6 sm:py-6 lg:px-8 lg:py-8",
        fullWidth ? "" : "max-w-[1400px]",
        className,
      ].join(" ")}
    >
      {children}
    </div>
  );
}
