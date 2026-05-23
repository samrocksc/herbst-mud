import { Link } from "@tanstack/react-router";
import type { ReactNode } from "react";

type PageHeaderProps = Readonly<{
  title: string
  showBack?: boolean
  backTo?: string
  backLabel?: string
  actions?: ReactNode
}>

export function PageHeader({ title, showBack, backTo, backLabel = "← Dashboard", actions }: PageHeaderProps) {
  return (
    <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-4 sm:mb-6">
      <div className="flex items-center gap-2 sm:gap-3 min-w-0">
        {showBack && backTo && (
          <Link
            to={backTo as "/dashboard" | "/map" | "/npcs" | "/items" | "/abilities" | "/quests" | "/logs" | "/skills" | "/factions"}
            className="no-underline px-2.5 py-1.5 sm:px-3 sm:py-1.5 rounded border border-border hover:border-primary transition-colors text-sm font-medium shrink-0"
          >
            {backLabel}
          </Link>
        )}
        <h1 className="m-0 text-lg sm:text-xl font-bold text-text truncate">{title}</h1>
      </div>
      {actions && (
        <div className="flex items-center gap-2 flex-wrap">{actions}</div>
      )}
    </div>
  );
}
