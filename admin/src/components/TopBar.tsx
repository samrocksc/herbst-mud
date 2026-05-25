import { useNavigate, useLocation } from "@tanstack/react-router";
import { navGroups, findActiveGroup } from "./navConfig";

export function TopBar() {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const activeGroup = findActiveGroup(pathname);

  return (
    <nav className="fixed top-0 left-0 right-0 z-40 h-14 bg-surface border-b border-border">
      <div className="flex items-center h-full px-3 gap-1 overflow-x-auto [&::-webkit-scrollbar]:hidden">
        {navGroups.map((group) => {
          const isActive = group === activeGroup;
          return (
            <button
              key={group.label}
              onClick={() => {
                if (!isActive) navigate({ to: group.items[0].path });
              }}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
                isActive
                  ? "bg-primary/10 text-primary"
                  : "text-text-muted hover:text-text hover:bg-surface-muted"
              }`}
            >
              <group.Icon />
              <span className="hidden sm:inline">{group.label}</span>
            </button>
          );
        })}
      </div>
    </nav>
  );
}
