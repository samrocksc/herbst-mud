import { useEffect } from "react";
import { Link } from "@tanstack/react-router";
import { WorldTitle } from "./WorldTitle";

import { DashboardIcon } from "./icons/DashboardIcon";
import { XpIcon } from "./icons/XpIcon";
import { ConfigIcon } from "./icons/ConfigIcon";
import { FactionsIcon } from "./icons/FactionsIcon";
import { ItemsIcon } from "./icons/ItemsIcon";
import { AbilitiesIcon } from "./icons/AbilitiesIcon";
import { SkillsIcon } from "./icons/SkillsIcon";
import { PlayersIcon } from "./icons/PlayersIcon";
import { MapIcon } from "./icons/MapIcon";
import { NPCsIcon } from "./icons/NPCsIcon";
import { ChevronLeftIcon } from "./icons/ChevronIcons";
import { ChevronRightIcon } from "./icons/ChevronIcons";

import { DocsIcon } from "./icons/DocsIcon";
import { TagsIcon } from "./icons/TagsIcon";
import { RacesIcon } from "./icons/RacesIcon";
import { LogsIcon } from "./icons/LogsIcon";
import { EffectsIcon } from "./icons/EffectsIcon";
import { QuestsIcon } from "./icons/QuestsIcon";
import { SocialsIcon } from "./icons/SocialsIcon";
import { ChannelsIcon } from "./icons/ChannelsIcon";
import { WorldIcon } from "./icons/WorldIcon";
import { CraftingIcon } from "./icons/CraftingIcon";

const navItems = [
  { label: "Dashboard", path: "/dashboard", Icon: DashboardIcon },
  { label: "XP", path: "/xp", Icon: XpIcon },
  { label: "Config", path: "/config", Icon: ConfigIcon },
  { label: "Factions", path: "/factions", Icon: FactionsIcon },
  { label: "Items", path: "/items", Icon: ItemsIcon },
  { label: "Recipes", path: "/recipes", Icon: CraftingIcon },
  { label: "Abilities", path: "/abilities", Icon: AbilitiesIcon },
  { label: "Effects", path: "/effects", Icon: EffectsIcon },
  { label: "Socials", path: "/socials", Icon: SocialsIcon },
  { label: "Channels", path: "/channels", Icon: ChannelsIcon },
  { label: "Skills", path: "/skills", Icon: SkillsIcon },
  { label: "Worlds", path: "/worlds", Icon: WorldIcon },
  { label: "Tags", path: "/tags", Icon: TagsIcon },
  { label: "Players", path: "/players", Icon: PlayersIcon },
  { label: "Characters", path: "/characters", Icon: PlayersIcon },
  { label: "Races", path: "/races", Icon: RacesIcon },
  { label: "Map", path: "/map", Icon: MapIcon },
  { label: "NPCs", path: "/npcs", Icon: NPCsIcon },
  { label: "Quests", path: "/quests", Icon: QuestsIcon },
  { label: "Logs", path: "/logs", Icon: LogsIcon },
  { label: "Docs", path: "/docs", Icon: DocsIcon },
];

/** Toggle button for collapsing/expanding the sidebar on desktop. */
function SidebarCollapseToggle({
  collapsed,
  onToggle,
}: Readonly<{
  collapsed: boolean;
  onToggle: () => void;
}>) {
  return (
    <button
      onClick={onToggle}
      aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      className="flex-shrink-0 flex items-center justify-center w-8 h-8 rounded hover:bg-surface-muted transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary"
      style={{ color: "var(--color-primary)" }}
    >
      {collapsed ? (
        <ChevronRightIcon stroke="var(--color-primary)" />
      ) : (
        <ChevronLeftIcon stroke="var(--color-primary)" />
      )}
    </button>
  );
}

/** Close button (X) for mobile dropdown. */
function SidebarCloseButton({ onClose }: Readonly<{ onClose: () => void }>) {
  return (
    <button
      onClick={onClose}
      aria-label="Close menu"
      title="Close menu"
      className="flex-shrink-0 flex items-center justify-center w-8 h-8 rounded hover:bg-surface-muted transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary md:hidden"
      style={{ color: "var(--color-primary)" }}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <line x1="18" y1="6" x2="6" y2="18" />
        <line x1="6" y1="6" x2="18" y2="18" />
      </svg>
    </button>
  );
}

export function Sidebar({
  mobileOpen,
  onMobileClose,
  collapsed,
  onToggleCollapse,
}: Readonly<{
  mobileOpen: boolean;
  onMobileClose: () => void;
  collapsed: boolean;
  onToggleCollapse: () => void;
}>) {
  // Sync collapsed prop to localStorage for persistence
  useEffect(() => {
    try {
      localStorage.setItem("sidebar-collapsed", String(collapsed));
    } catch {
      // ignore
    }
  }, [collapsed]);

  return (
    <nav
      className={[
        "fixed inset-y-0 z-50 flex flex-col bg-surface",
        // Position: right edge on mobile, left edge on desktop
        "right-0 md:left-0 md:right-auto",
        // Border: left border on mobile (drawer's left edge), right on desktop
        "border-l border-border md:border-r md:border-l-0",
        "transform transition-transform duration-300 ease-in-out",
        // Mobile: slide from right; desktop: always visible
        !mobileOpen && "translate-x-full md:translate-x-0",
        mobileOpen && "translate-x-0",
        // Width: mobile full-screen, desktop sidebar
        "w-full md:w-[64px]",
        collapsed ? "lg:w-[64px]" : "lg:w-[220px]",
      ].filter(Boolean).join(" ")}
    >
      {/* Header: WorldTitle + close button (mobile) + collapse toggle (desktop) */}
      <div className="flex items-center border-b border-border flex-shrink-0 h-14 px-1">
        <div
          className={[
            "flex-1 min-w-0 px-1 overflow-hidden",
            "transition-opacity duration-300",
            collapsed ? "md:opacity-0 md:select-none" : "md:opacity-100",
          ].join(" ")}
        >
          <WorldTitle />
        </div>
        <SidebarCloseButton onClose={onMobileClose} />
        <div className="hidden md:block">
          <SidebarCollapseToggle
            collapsed={collapsed}
            onToggle={onToggleCollapse}
          />
        </div>
      </div>

      {/* Nav items */}
      <div className="flex flex-col p-2 gap-1 flex-1 overflow-y-auto">
        {navItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            search={(prev: Record<string, string>) => prev}
            activeProps={{
              className:
                "bg-primary !text-white font-semibold border-l-4 border-primary",
            }}
            inactiveProps={{
              className:
                "text-text-muted hover:bg-surface-muted hover:text-text",
            }}
            className={[
              "flex items-center gap-3 px-3 py-2.5 rounded text-sm",
              "no-underline transition-colors",
              collapsed ? "md:justify-center md:px-0" : "",
            ].join(" ")}
            title={collapsed ? item.label : undefined}
            onClick={onMobileClose}
          >
            <span className="flex-shrink-0">
              <item.Icon />
            </span>
            <span
              className={[
                "whitespace-nowrap transition-opacity duration-300 min-w-0",
                collapsed ? "md:opacity-0 md:pointer-events-none md:w-0 md:overflow-hidden" : "md:opacity-100",
              ].join(" ")}
            >
              {item.label}
            </span>
          </Link>
        ))}
      </div>
    </nav>
  );
}
