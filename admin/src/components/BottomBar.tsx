import { Link, useLocation } from "@tanstack/react-router";
import { navGroups, findActiveGroup } from "./navConfig";

export function BottomBar() {
  const { pathname } = useLocation();
  const activeGroup = findActiveGroup(pathname);
  const items = activeGroup.items;

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-surface border-t border-border">
      <div className="flex items-center h-16 px-2 gap-1 overflow-x-auto [&::-webkit-scrollbar]:hidden">
        {items.map((item) => {
          const isActive =
            pathname === item.path || pathname.startsWith(item.path + "/");
          return (
            <Link
              key={item.path}
              to={item.path}
              search={(prev: Record<string, string>) => prev}
              className={`flex flex-col items-center justify-center gap-0.5 min-w-[64px] px-2 py-1.5 rounded-lg no-underline transition-colors whitespace-nowrap ${
                isActive
                  ? "text-primary"
                  : "text-text-muted hover:text-text hover:bg-surface-muted"
              }`}
            >
              <item.Icon />
              <span className="text-[10px] font-medium">{item.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
