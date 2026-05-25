import { Link } from "@tanstack/react-router";
import { useLocation } from "@tanstack/react-router";
import { DashboardIcon } from "./icons/DashboardIcon";
import { MapIcon } from "./icons/MapIcon";
import { PlayersIcon } from "./icons/PlayersIcon";
import { ItemsIcon } from "./icons/ItemsIcon";
import { LogsIcon } from "./icons/LogsIcon";

const mobileTabs = [
  { label: "Dashboard", path: "/dashboard", Icon: DashboardIcon },
  { label: "Map", path: "/map", Icon: MapIcon },
  { label: "Players", path: "/players", Icon: PlayersIcon },
  { label: "Items", path: "/items", Icon: ItemsIcon },
  { label: "Logs", path: "/logs", Icon: LogsIcon },
];

export function MobileNavBar() {
  const { pathname } = useLocation();

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-surface border-t border-border md:hidden">
      <div className="flex items-center justify-around h-16">
        {mobileTabs.map((tab) => {
          const active = pathname === tab.path || pathname.startsWith(tab.path + "/");
          return (
            <Link
              key={tab.path}
              to={tab.path}
              className={[
                "flex flex-col items-center justify-center gap-0.5 w-16 h-full no-underline",
                active ? "text-primary" : "text-text-muted",
              ].join(" ")}
            >
              <span className={active ? "text-primary" : "text-text-muted"}>
                <tab.Icon />
              </span>
              <span className="text-[10px] font-medium">{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
