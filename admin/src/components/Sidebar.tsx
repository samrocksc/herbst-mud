 
import { useState, useEffect } from "react";
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

const navItems = [
  { label: "Dashboard", path: "/dashboard", Icon: DashboardIcon },
  { label: "XP", path: "/xp", Icon: XpIcon },
  { label: "Config", path: "/config", Icon: ConfigIcon },
  { label: "Factions", path: "/factions", Icon: FactionsIcon },
  { label: "Items", path: "/items", Icon: ItemsIcon },
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

/** Toggle button for collapsing/expanding the sidebar. Named component for DevTools clarity. */
function SidebarCollapseToggle({
  collapsed,
  onToggle,
}: Readonly<{
  collapsed: boolean
  onToggle: () => void
}>) {
  return (
    <button
      onClick={onToggle}
      aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      className={[
        "flex-shrink-0 flex items-center justify-center",
        "w-8 h-8 rounded",
        "hover:bg-surface-muted",
        "transition-colors duration-200",
        "focus:outline-none focus:ring-2 focus:ring-primary",
      ].join(" ")}
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

// Collapsed state key
const COLLAPSED_KEY = "sidebar-collapsed";

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(() => {
    try {
      return localStorage.getItem(COLLAPSED_KEY) === "true";
    } catch {
      return false;
    }
  });

  // Update localStorage when collapsed changes
  useEffect(() => {
    try {
      localStorage.setItem(COLLAPSED_KEY, String(collapsed));
    } catch {
      // Ignore errors
    }
  }, [collapsed]);

  return (
    <nav
      className={[
        "h-screen bg-surface border-r border-border",
        "flex flex-col",
        "transition-all duration-300 ease-in-out",
        "relative",
        // Mobile: sidebar slides in/out as an overlay
        "fixed inset-y-0 left-0 z-40",
        "max-w-[220px]",
        collapsed ? "w-[64px] max-w-[64px]" : "w-[220px]",
        // Desktop: always visible, not fixed
        "lg:relative lg:inset-auto lg:z-auto",
      ].join(" ")}
    >
      {/* Mobile overlay backdrop — only visible when collapsed on mobile */}
      {/* Header + toggle */}
      <div className="flex items-center border-b border-border flex-shrink-0 h-14 px-1">
        <div
          className={[
            "flex-1 min-w-0 px-1 overflow-hidden",
            "transition-opacity duration-300",
            collapsed ? "opacity-0 select-none" : "opacity-100",
          ].join(" ")}
        >
          <WorldTitle />
        </div>
        <SidebarCollapseToggle
          collapsed={collapsed}
          onToggle={() => setCollapsed((c) => !c)}
        />
      </div>

      {/* Nav items — the only scrollable region */}
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
              collapsed ? "justify-center px-0" : "",
            ].join(" ")}
            title={collapsed ? item.label : undefined}
          >
            <span className="flex-shrink-0">
              <item.Icon />
            </span>
            <span
              className={[
                "whitespace-nowrap transition-opacity duration-300 min-w-0",
                collapsed ? "opacity-0 pointer-events-none w-0 overflow-hidden" : "opacity-100",
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