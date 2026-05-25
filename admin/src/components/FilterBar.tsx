import { useState, type ReactNode } from "react";
import { Button } from "./Button";

type FilterBarProps = Readonly<{
  children: ReactNode
  showClear?: boolean
  onClear?: () => void
}>

/**
 * Responsive filter bar for management list pages.
 * Desktop: horizontal row. Mobile: vertical stack with optional collapse.
 */
export function FilterBar({ children, showClear, onClear }: FilterBarProps) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="mb-4 sm:mb-6 border border-border rounded-lg bg-surface-muted">
      {/* Mobile toggle */}
      <div className="sm:hidden flex items-center justify-between px-3 py-2">
        <span className="text-sm font-medium text-text-muted">Filters</span>
        <div className="flex items-center gap-2">
          {showClear && onClear && (
            <Button variant="ghost" size="sm" onClick={onClear}>Clear</Button>
          )}
          <Button variant="ghost" size="sm" onClick={() => setExpanded((v) => !v)}>
            {expanded ? "Hide" : "Show"}
          </Button>
        </div>
      </div>

      {/* Content */}
      <div className={[
        "p-3 sm:p-4 flex flex-col sm:flex-row gap-3 sm:items-end flex-wrap",
        expanded ? "block" : "hidden sm:flex",
      ].join(" ")}>
        {children}
        {showClear && onClear && (
          <div className="hidden sm:block">
            <Button variant="ghost" size="sm" onClick={onClear}>Clear</Button>
          </div>
        )}
      </div>
    </div>
  );
}
